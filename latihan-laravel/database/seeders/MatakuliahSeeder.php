<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\MataKuliah;

class MatakuliahSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $daftar = [
            ['kode' => 'TK301', 'nama' => 'Pemrograman Web II', 'sks' => 3, 'semester' => 5],
            ['kode' => 'TK302', 'nama' => 'Basis Data Lanjut', 'sks' => 3, 'semester' => 5],
            ['kode' => 'TK303', 'nama' => 'Jaringan Komputer', 'sks' => 2, 'semester' => 5],
            ['kode' => 'TK304', 'nama' => 'Sistem Operasi', 'sks' => 3, 'semester' => 4],
            ['kode' => 'TK305', 'nama' => 'Kecerdasan Buatan', 'sks' => 2, 'semester' => 6],
            ['kode' => 'TK306', 'nama' => 'Rekayasa Perangkat Lunak', 'sks' => 3, 'semester' => 4],
        ];

        foreach ($daftar as $item) {
            MataKuliah::create($item);
        }
    }
}
