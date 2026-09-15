<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use App\Models\Mahasiswa;
use App\Models\MataKuliah;

class DatabaseSeeder extends Seeder
{
    use WithoutModelEvents;

    /**
     * Seed the application's database.
     */
    public function run(): void
    {
        $this->call(ProgramStudiSeeder::class);
        $this->call(MatakuliahSeeder::class);

        Mahasiswa::factory()->count(30)->create();

        // Isi tabel pivot mahasiswa_matakuliah: tiap mahasiswa
        // mengambil 3-5 matakuliah acak beserta nilainya.
        $daftarMatakuliah = MataKuliah::pluck('id');

        Mahasiswa::all()->each(function (Mahasiswa $mahasiswa) use ($daftarMatakuliah) {
            $diambil = $daftarMatakuliah->random(random_int(3, 5));
            $data = [];

            foreach ($diambil as $matakuliahId) {
                $data[$matakuliahId] = [
                    'nilai' => fake()->randomFloat(2, 60, 100),
                ];
            }

            $mahasiswa->matakuliah()->attach($data);
        });
    }
}
