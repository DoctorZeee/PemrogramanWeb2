<?php

namespace Tests\Feature;

use App\Models\Mahasiswa;
use App\Models\ProgramStudi;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Laravel\Sanctum\Sanctum;
use Tests\TestCase;

class MahasiswaApiTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();

        ProgramStudi::create([
            'kode' => 'TK',
            'nama' => 'Teknik Komputer',
            'jenjang' => 'S1',
        ]);
    }

    private function masukSebagaiAdmin(): User
    {
        $admin = User::factory()->create(['peran' => 'admin']);
        Sanctum::actingAs($admin, ['mahasiswa:baca', 'mahasiswa:tulis']);

        return $admin;
    }

    private function dataMahasiswa(array $tambahan = []): array
    {
        return array_merge([
            'program_studi_id' => 1,
            'nim' => 'H1H024043',
            'nama' => 'Dimas Rafif Zaidan',
            'email' => 'dimas@example.com',
            'angkatan' => 2024,
            'ipk' => 3.55,
            'aktif' => true,
        ], $tambahan);
    }

    public function test_daftar_mahasiswa_menolak_permintaan_tanpa_token(): void
    {
        $this->getJson('/api/mahasiswa')->assertUnauthorized();
    }

    public function test_daftar_mahasiswa_mengembalikan_data_dan_meta_paginasi(): void
    {
        Sanctum::actingAs(User::factory()->create(), ['mahasiswa:baca']);
        Mahasiswa::factory()->count(3)->create();

        $this->getJson('/api/mahasiswa?per_halaman=2')
            ->assertOk()
            ->assertJsonStructure([
                'data' => [
                    '*' => ['id', 'nim', 'nama', 'email', 'angkatan', 'ipk', 'aktif'],
                ],
                'meta' => ['current_page', 'last_page', 'per_page', 'total'],
            ])
            ->assertJsonPath('meta.total', 3)
            ->assertJsonPath('meta.per_page', 2);
    }

    public function test_membuat_mahasiswa_dengan_data_valid(): void
    {
        $this->masukSebagaiAdmin();

        $this->postJson('/api/mahasiswa', $this->dataMahasiswa())
            ->assertCreated()
            ->assertJsonPath('sukses', true)
            ->assertJsonPath('data.nim', 'H1H024043');

        $this->assertDatabaseHas('mahasiswas', [
            'nim' => 'H1H024043',
            'email' => 'dimas@example.com',
        ]);
    }

    public function test_membuat_mahasiswa_dengan_data_tidak_valid_ditolak_422(): void
    {
        $this->masukSebagaiAdmin();

        $this->postJson('/api/mahasiswa', $this->dataMahasiswa([
            'email' => 'bukan-email',
            'ipk' => 5.25,
        ]))->assertUnprocessable();

        $this->assertDatabaseCount('mahasiswas', 0);
    }

    public function test_detail_mahasiswa_mengembalikan_data(): void
    {
        Sanctum::actingAs(User::factory()->create(), ['mahasiswa:baca']);
        $mahasiswa = Mahasiswa::factory()->create();

        $this->getJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertOk()
            ->assertJsonPath('data.id', $mahasiswa->id)
            ->assertJsonPath('data.nim', $mahasiswa->nim);
    }

    public function test_memperbarui_mahasiswa_dengan_token_yang_memiliki_kemampuan_tulis(): void
    {
        $this->masukSebagaiAdmin();
        $mahasiswa = Mahasiswa::factory()->create(['nama' => 'Nama Lama']);

        $this->putJson('/api/mahasiswa/' . $mahasiswa->id, [
            'nama' => 'Nama Baru',
        ])
            ->assertOk()
            ->assertJsonPath('data.nama', 'Nama Baru');

        $this->assertDatabaseHas('mahasiswas', [
            'id' => $mahasiswa->id,
            'nama' => 'Nama Baru',
        ]);
    }

    public function test_detail_mahasiswa_yang_tidak_ada_mengembalikan_404(): void
    {
        Sanctum::actingAs(User::factory()->create(), ['mahasiswa:baca']);

        $this->getJson('/api/mahasiswa/999999')->assertNotFound();
    }

    public function test_penghapusan_mahasiswa_oleh_non_admin_ditolak_403(): void
    {
        // Kemampuan tulis dimiliki, tetapi peran bukan admin.
        $pengguna = User::factory()->create(['peran' => 'mahasiswa']);
        Sanctum::actingAs($pengguna, ['mahasiswa:baca', 'mahasiswa:tulis']);
        $mahasiswa = Mahasiswa::factory()->create();

        $this->deleteJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertForbidden();

        $this->assertDatabaseHas('mahasiswas', ['id' => $mahasiswa->id]);
    }

    public function test_penghapusan_mahasiswa_oleh_admin_berhasil(): void
    {
        $admin = User::factory()->create(['peran' => 'admin']);
        Sanctum::actingAs($admin, ['mahasiswa:baca', 'mahasiswa:tulis']);
        $mahasiswa = Mahasiswa::factory()->create();

        $this->deleteJson('/api/mahasiswa/' . $mahasiswa->id)
            ->assertOk()
            ->assertJsonPath('sukses', true);

        $this->assertDatabaseMissing('mahasiswas', ['id' => $mahasiswa->id]);
    }
}
