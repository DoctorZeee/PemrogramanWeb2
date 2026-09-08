<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\MahasiswaController;
use App\Http\Controllers\MatakuliahController;

Route::get('/', function () {
    return view('welcome');
});

Route::get('/salam', function () {
    return 'Selamat datang di Pemrograman Web II';
});

// data mahasiswa ================================================
Route::get('/data-mahasiswa', [
    MahasiswaController::class,
    'index'
])->name('mahasiswa.index');

Route::get('/data-mahasiswa/{nim}', [
    MahasiswaController::class,
    'show'
])->name('mahasiswa.show');

Route::get('/cari-mahasiswa', [MahasiswaController::class, 'cari']);
// =================================================================

// data matakuliah =================================================
Route::get('/data-matakuliah', [
    MatakuliahController::class,
    'index'
])->name('matakuliah.index');

Route::get('/data-matakuliah/{kode}', [
    MatakuliahController::class,
    'show'
])->name('matakuliah.show');

Route::get('/cari-matakuliah', [MatakuliahController::class, 'cari']);
// ==================================================================

Route::get('/semester/{angka}', function (int $angka) {
    return 'Semester ke ' . $angka;
})->whereNumber('angka');




