<?php

namespace Database\Seeders;

use App\Models\Mahasiswa;
use Illuminate\Database\Seeder;

class MahasiswaSeeder extends Seeder
{
    public function run(): void
    {
        $data = [
            ['program_studi_id' => 1, 'nim' => 'H1A125001', 'nama' => 'Dewi Anggraini', 'email' => 'dewi@example.com', 'angkatan' => 2025, 'ipk' => 3.65, 'aktif' => true],
            ['program_studi_id' => 1, 'nim' => 'H1A125002', 'nama' => 'Budi Santoso', 'email' => 'budi@example.com', 'angkatan' => 2025, 'ipk' => 3.20, 'aktif' => true],
            ['program_studi_id' => 2, 'nim' => 'H1A125003', 'nama' => 'Citra Lestari', 'email' => 'citra@example.com', 'angkatan' => 2024, 'ipk' => 3.85, 'aktif' => true],
            ['program_studi_id' => 2, 'nim' => 'H1A125004', 'nama' => 'Dimas Prakoso', 'email' => 'dimas@example.com', 'angkatan' => 2024, 'ipk' => 3.50, 'aktif' => false],
        ];

        foreach ($data as $item) {
            Mahasiswa::create($item);
        }
    }
}