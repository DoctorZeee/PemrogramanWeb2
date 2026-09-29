<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class PeranAdmin
{
    /**
     * Menolak permintaan apabila pengguna yang sedang login bukan admin.
     * Middleware ini harus dipasang SETELAH auth:sanctum.
     */
    public function handle(Request $request, Closure $next): Response
    {
        $pengguna = $request->user();

        if ($pengguna === null || $pengguna->peran !== 'admin') {
            return response()->json([
                'sukses' => false,
                'pesan' => 'Tindakan ini hanya boleh dilakukan oleh admin',
            ], 403);
        }

        return $next($request);
    }
}
