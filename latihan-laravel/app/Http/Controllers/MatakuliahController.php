<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class MatakuliahController extends Controller
{
    public function index()
    {
        $daftarMatakuliah = [
            ['kode' => 'MK001', 'nama' => 'Pemrograman Web II', 'sks' => 3],
            ['kode' => 'MK002', 'nama' => 'Basis Data Lanjut', 'sks' => 3],
            ['kode' => 'MK003', 'nama' => 'Jaringan Komputer', 'sks' => 2],
            ['kode' => 'MK004', 'nama' => 'Sistem Operasi', 'sks' => 3],
            ['kode' => 'MK005', 'nama' => 'Kecerdasan Buatan', 'sks' => 2],
        ];
        return view('matakuliah.index', ['daftarMatakuliah' =>
        $daftarMatakuliah]);
    }

    public function show(string $kode)
    {
        return view('matakuliah.show', ['kode' => $kode]);
    }

    public function cari(Request $request)
    {
        $kataKunci = $request->query('q', '');
        return response()->json([
            'kata_kunci' => $kataKunci,
            'metode' => $request->method(),
            'path' => $request->path(),
        ]);
    }
}
