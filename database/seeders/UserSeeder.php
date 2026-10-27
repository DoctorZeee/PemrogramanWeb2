<?php

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Seeder;

class UserSeeder extends Seeder
{
    /**
     * Pengganti Langkah 6 (tinker): membuat akun admin secara idempoten.
     */
    public function run(): void
    {
        User::updateOrCreate(
            ['email' => 'admin@unsoed.ac.id'],
            [
                'name' => 'Administrator',
                'password' => 'rahasia123', // di-hash otomatis oleh cast
                'peran' => 'admin',
            ]
        );
    }
}
