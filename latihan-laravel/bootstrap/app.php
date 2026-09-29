<?php

use App\Http\Middleware\PeranAdmin;
use Illuminate\Auth\AuthenticationException;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Laravel\Sanctum\Http\Middleware\CheckAbilities;
use Laravel\Sanctum\Http\Middleware\CheckForAnyAbility;
use Symfony\Component\HttpKernel\Exception\AccessDeniedHttpException;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\HttpKernel\Exception\TooManyRequestsHttpException;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__ . '/../routes/web.php',
        api: __DIR__ . '/../routes/api.php',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        // Sejak Laravel 11 alias middleware Sanctum harus didaftarkan manual.
        $middleware->alias([
            'abilities' => CheckAbilities::class,      // harus punya SEMUA kemampuan
            'ability' => CheckForAnyAbility::class,    // cukup punya SALAH SATU
            'peran.admin' => PeranAdmin::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions) {
        $exceptions->render(function (NotFoundHttpException $e, Request
        $request) {
            if ($request->is('api/*')) {
                return response()->json([
                    'sukses' => false,
                    'pesan' => 'Sumber daya tidak ditemukan',
                ], 404);
            }
        });
        $exceptions->render(function (ValidationException $e, Request
        $request) {
            if ($request->is('api/*')) {
                return response()->json([
                    'sukses' => false,
                    'pesan' => 'Data yang dikirim tidak valid',
                    'galat' => $e->errors(),
                ], 422);
            }
        });
        // Langkah 10: token tidak dikirim / tidak valid -> 401
        $exceptions->render(function (AuthenticationException $e, Request
        $request) {
            if ($request->is('api/*')) {
                return response()->json([
                    'sukses' => false,
                    'pesan' => 'Token tidak valid atau belum dikirim',
                ], 401);
            }
        });
        // Token valid tetapi kemampuan tidak cukup (MissingAbilityException
        // diubah Laravel menjadi AccessDeniedHttpException) -> 403
        $exceptions->render(function (AccessDeniedHttpException $e, Request
        $request) {
            if ($request->is('api/*')) {
                return response()->json([
                    'sukses' => false,
                    'pesan' => 'Anda tidak memiliki izin untuk tindakan ini',
                ], 403);
            }
        });
        // Langkah 11: pembatasan laju -> 429
        $exceptions->render(function (TooManyRequestsHttpException $e, Request
        $request) {
            if ($request->is('api/*')) {
                return response()->json([
                    'sukses' => false,
                    'pesan' => 'Terlalu banyak percobaan, coba lagi nanti',
                ], 429, $e->getHeaders());
            }
        });
    })->create();
