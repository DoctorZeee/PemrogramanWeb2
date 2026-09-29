<?php

use App\Http\Controllers\Api\AuthController;
use App\Http\Controllers\Api\MahasiswaController;
use Illuminate\Support\Facades\Route;

Route::post('/auth/register', [AuthController::class, 'register']);
// Langkah 11: maksimal 5 percobaan per menit per IP
Route::post('/auth/login', [AuthController::class, 'login'])
    ->middleware('throttle:5,1');

Route::middleware('auth:sanctum')->group(function () {
    Route::get('/auth/profil', [AuthController::class, 'profil']);
    Route::put('/auth/password', [AuthController::class, 'ubahPassword']);
    Route::post('/auth/logout', [AuthController::class, 'logout']);
    Route::post('/auth/logout-semua', [
        AuthController::class,
        'logoutSemua'
    ]);
    Route::get('/mahasiswa', [MahasiswaController::class, 'index']);
    Route::get('/mahasiswa/{mahasiswa}', [
        MahasiswaController::class,
        'show'
    ]);
    Route::middleware('ability:mahasiswa:tulis')->group(function () {
        Route::post('/mahasiswa', [MahasiswaController::class, 'store']);
        Route::put('/mahasiswa/{mahasiswa}', [
            MahasiswaController::class,
            'update'
        ]);
        Route::patch('/mahasiswa/{mahasiswa}', [
            MahasiswaController::class,
            'update'
        ]);
        // Tugas 2: penghapusan hanya untuk peran admin
        Route::delete('/mahasiswa/{mahasiswa}', [
            MahasiswaController::class,
            'destroy'
        ])->middleware('peran.admin');
    });
});
