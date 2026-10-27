<?php

namespace Tests\Feature;

use App\Models\Mahasiswa;
use App\Models\User;
use Database\Seeders\ProgramStudiSeeder;
use Database\Seeders\UserSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\RateLimiter;
use Tests\TestCase;

/**
 * Pengujian otomatis Modul 5 (Laravel Sanctum).
 * Setiap test setara dengan satu skenario pada koleksi Postman.
 */
class AuthSanctumTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();
        $this->seed([UserSeeder::class, ProgramStudiSeeder::class]);
        RateLimiter::clear('');
    }

    /** Login lalu kembalikan token teks-biasa. */
    private function login(string $email, string $password): string
    {
        $respons = $this->postJson('/api/auth/login', [
            'email' => $email,
            'password' => $password,
        ]);
        $respons->assertOk();

        return $respons->json('data.token');
    }

    /**
     * Kirim permintaan dengan Bearer token. Guard di-reset agar tiap
     * permintaan dalam satu test benar-benar dibaca ulang dari database
     * (tanpa ini, pengguna dari permintaan sebelumnya ikut ter-cache).
     */
    private function denganToken(string $token)
    {
        $this->app['auth']->forgetGuards();

        return $this->withHeader('Authorization', 'Bearer ' . $token);
    }

    private function buatMahasiswaBiasa(): string
    {
        $this->postJson('/api/auth/register', [
            'name' => 'Budi',
            'email' => 'budi@unsoed.ac.id',
            'password' => 'rahasia123',
            'password_confirmation' => 'rahasia123',
        ])->assertCreated();

        return $this->login('budi@unsoed.ac.id', 'rahasia123');
    }

    // ---------------------------------------------------------- registrasi

    public function test_registrasi_berhasil_dan_peran_dipaksa_mahasiswa(): void
    {
        $this->postJson('/api/auth/register', [
            'name' => 'Sinta',
            'email' => 'sinta@unsoed.ac.id',
            'password' => 'rahasia123',
            'password_confirmation' => 'rahasia123',
            'peran' => 'admin', // upaya eskalasi hak akses
        ])
            ->assertCreated()
            ->assertJsonPath('data.pengguna.peran', 'mahasiswa')
            ->assertJsonStructure(['data' => ['token']]);

        $this->assertDatabaseHas('users', [
            'email' => 'sinta@unsoed.ac.id',
            'peran' => 'mahasiswa',
        ]);
    }

    public function test_registrasi_ditolak_bila_data_tidak_valid(): void
    {
        $this->postJson('/api/auth/register', [
            'name' => 'X',
            'email' => 'bukan-email',
            'password' => 'pendek',
        ])
            ->assertStatus(422)
            ->assertJsonPath('sukses', false)
            ->assertJsonStructure(['galat' => ['email', 'password']]);
    }

    public function test_token_dari_registrasi_tidak_boleh_menulis(): void
    {
        $token = $this->postJson('/api/auth/register', [
            'name' => 'Sinta',
            'email' => 'sinta@unsoed.ac.id',
            'password' => 'rahasia123',
            'password_confirmation' => 'rahasia123',
        ])->json('data.token');

        $this->denganToken($token)->postJson('/api/mahasiswa', [])
            ->assertStatus(403);
    }

    // ---------------------------------------------------------------- login

    public function test_login_berhasil_mengembalikan_token(): void
    {
        $this->postJson('/api/auth/login', [
            'email' => 'admin@unsoed.ac.id',
            'password' => 'rahasia123',
        ])
            ->assertOk()
            ->assertJsonPath('sukses', true)
            ->assertJsonPath('data.pengguna.peran', 'admin')
            ->assertJsonStructure(['data' => ['token']]);
    }

    public function test_login_gagal_karena_password_salah(): void
    {
        $this->postJson('/api/auth/login', [
            'email' => 'admin@unsoed.ac.id',
            'password' => 'salah-total',
        ])
            ->assertStatus(401)
            ->assertJsonPath('sukses', false)
            ->assertJsonPath('pesan', 'Email atau kata sandi tidak sesuai');
    }

    public function test_login_gagal_karena_email_tidak_terdaftar(): void
    {
        $this->postJson('/api/auth/login', [
            'email' => 'tidak-ada@unsoed.ac.id',
            'password' => 'rahasia123',
        ])->assertStatus(401);
    }

    public function test_login_dibatasi_lima_percobaan_per_menit(): void
    {
        for ($i = 0; $i < 5; $i++) {
            $this->postJson('/api/auth/login', [
                'email' => 'admin@unsoed.ac.id',
                'password' => 'salah',
            ])->assertStatus(401);
        }

        $this->postJson('/api/auth/login', [
            'email' => 'admin@unsoed.ac.id',
            'password' => 'salah',
        ])
            ->assertStatus(429)
            ->assertJsonPath('sukses', false)
            ->assertHeader('Retry-After');
    }

    // Tugas 3
    public function test_login_berhasil_memperbarui_terakhir_login(): void
    {
        $admin = User::where('email', 'admin@unsoed.ac.id')->first();
        $this->assertNull($admin->terakhir_login);

        $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->assertNotNull($admin->fresh()->terakhir_login);
    }

    public function test_login_gagal_tidak_memperbarui_terakhir_login(): void
    {
        $this->postJson('/api/auth/login', [
            'email' => 'admin@unsoed.ac.id',
            'password' => 'salah',
        ])->assertStatus(401);

        $this->assertNull(
            User::where('email', 'admin@unsoed.ac.id')->first()->terakhir_login
        );
    }

    // --------------------------------------------------- proteksi rute (401)

    public function test_akses_tanpa_token_ditolak_401(): void
    {
        $this->getJson('/api/auth/profil')
            ->assertStatus(401)
            ->assertJsonPath('pesan', 'Token tidak valid atau belum dikirim');
    }

    public function test_akses_dengan_token_palsu_ditolak_401(): void
    {
        $this->denganToken('1|tokenpalsu')->getJson('/api/auth/profil')
            ->assertStatus(401);
    }

    public function test_profil_menampilkan_kemampuan_token(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->getJson('/api/auth/profil')
            ->assertOk()
            ->assertJsonPath('data.email', 'admin@unsoed.ac.id')
            ->assertJsonPath(
                'data.kemampuan',
                ['mahasiswa:baca', 'mahasiswa:tulis']
            );
    }

    // ---------------------------------------------------------------- logout

    public function test_token_tidak_berlaku_setelah_logout(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->postJson('/api/auth/logout')
            ->assertOk()
            ->assertJsonPath('pesan', 'Logout berhasil');

        $this->assertDatabaseCount('personal_access_tokens', 0);

        $this->denganToken($token)->getJson('/api/auth/profil')
            ->assertStatus(401);
    }

    public function test_logout_hanya_mencabut_token_perangkat_ini(): void
    {
        $tokenA = $this->login('admin@unsoed.ac.id', 'rahasia123');
        $tokenB = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($tokenA)->postJson('/api/auth/logout')->assertOk();

        $this->denganToken($tokenB)->getJson('/api/auth/profil')->assertOk();
    }

    public function test_logout_semua_mencabut_seluruh_token(): void
    {
        $tokenA = $this->login('admin@unsoed.ac.id', 'rahasia123');
        $tokenB = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($tokenA)->postJson('/api/auth/logout-semua')
            ->assertOk();

        $this->denganToken($tokenB)->getJson('/api/auth/profil')
            ->assertStatus(401);
    }

    // ------------------------------------------- kemampuan token (403)

    public function test_mahasiswa_boleh_membaca_data(): void
    {
        $token = $this->buatMahasiswaBiasa();

        $this->denganToken($token)->getJson('/api/mahasiswa')->assertOk();
    }

    public function test_akses_tanpa_kemampuan_menulis_ditolak_403(): void
    {
        $token = $this->buatMahasiswaBiasa();

        $this->denganToken($token)->postJson('/api/mahasiswa', [
            'program_studi_id' => 1,
            'nim' => 'H1A000001',
            'nama' => 'Uji Coba',
            'email' => 'uji@unsoed.ac.id',
            'angkatan' => 2025,
        ])
            ->assertStatus(403)
            ->assertJsonPath('sukses', false);

        $this->assertDatabaseMissing('mahasiswas', ['nim' => 'H1A000001']);
    }

    public function test_admin_boleh_menulis_data(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->postJson('/api/mahasiswa', [
            'program_studi_id' => 1,
            'nim' => 'H1A000001',
            'nama' => 'Uji Coba',
            'email' => 'uji@unsoed.ac.id',
            'angkatan' => 2025,
        ])->assertCreated();
    }

    // -------------------------------------------- Tugas 2: middleware PeranAdmin

    public function test_hapus_mahasiswa_ditolak_untuk_non_admin_meski_token_menulis(): void
    {
        // Pengguna berperan 'mahasiswa' tetapi memegang token yang
        // kebetulan berkemampuan tulis: PeranAdmin tetap harus menolak.
        $pengguna = User::factory()->create(['peran' => 'mahasiswa']);
        $token = $pengguna->createToken('uji', ['mahasiswa:tulis'])
            ->plainTextToken;
        $mahasiswa = Mahasiswa::factory()->create();

        $this->denganToken($token)
            ->deleteJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertStatus(403)
            ->assertJsonPath(
                'pesan',
                'Tindakan ini hanya boleh dilakukan oleh admin'
            );

        $this->assertDatabaseHas('mahasiswas', ['id' => $mahasiswa->id]);
    }

    public function test_hapus_mahasiswa_berhasil_untuk_admin(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');
        $mahasiswa = Mahasiswa::factory()->create();

        $this->denganToken($token)
            ->deleteJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertOk();

        $this->assertDatabaseMissing('mahasiswas', ['id' => $mahasiswa->id]);
    }

    public function test_hapus_mahasiswa_tanpa_token_ditolak_401(): void
    {
        $mahasiswa = Mahasiswa::factory()->create();

        $this->deleteJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertStatus(401);
    }

    // ------------------------------------------- Tugas 1: PUT /auth/password

    public function test_ubah_password_berhasil(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->putJson('/api/auth/password', [
            'password_lama' => 'rahasia123',
            'password' => 'baru12345',
            'password_confirmation' => 'baru12345',
        ])
            ->assertOk()
            ->assertJsonPath('sukses', true);

        // password lama tidak berlaku, password baru berlaku
        $this->postJson('/api/auth/login', [
            'email' => 'admin@unsoed.ac.id',
            'password' => 'rahasia123',
        ])->assertStatus(401);
        $this->login('admin@unsoed.ac.id', 'baru12345');
    }

    public function test_ubah_password_ditolak_bila_password_lama_salah(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->putJson('/api/auth/password', [
            'password_lama' => 'bukan-password-lama',
            'password' => 'baru12345',
            'password_confirmation' => 'baru12345',
        ])
            ->assertStatus(422)
            ->assertJsonStructure(['galat' => ['password_lama']]);

        $this->assertTrue(
            password_verify(
                'rahasia123',
                User::where('email', 'admin@unsoed.ac.id')->first()->password
            )
        );
    }

    public function test_ubah_password_ditolak_bila_konfirmasi_tidak_cocok(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->putJson('/api/auth/password', [
            'password_lama' => 'rahasia123',
            'password' => 'baru12345',
            'password_confirmation' => 'beda12345',
        ])->assertStatus(422)->assertJsonStructure(['galat' => ['password']]);
    }

    public function test_ubah_password_ditolak_bila_sama_dengan_yang_lama(): void
    {
        $token = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($token)->putJson('/api/auth/password', [
            'password_lama' => 'rahasia123',
            'password' => 'rahasia123',
            'password_confirmation' => 'rahasia123',
        ])->assertStatus(422);
    }

    public function test_ubah_password_mencabut_token_perangkat_lain(): void
    {
        $tokenA = $this->login('admin@unsoed.ac.id', 'rahasia123');
        $tokenB = $this->login('admin@unsoed.ac.id', 'rahasia123');

        $this->denganToken($tokenA)->putJson('/api/auth/password', [
            'password_lama' => 'rahasia123',
            'password' => 'baru12345',
            'password_confirmation' => 'baru12345',
        ])->assertOk();

        $this->denganToken($tokenA)->getJson('/api/auth/profil')->assertOk();
        $this->denganToken($tokenB)->getJson('/api/auth/profil')
            ->assertStatus(401);
    }

    public function test_ubah_password_tanpa_token_ditolak_401(): void
    {
        $this->putJson('/api/auth/password', [])->assertStatus(401);
    }
}
