@extends('layouts.app')
@section('judul', 'Daftar Matakuliah')
@section('konten')
    <h1 class="h3 mb-4">Daftar Matakuliah</h1>
    <table class="table table-bordered bg-white">
        <thead>
            <tr>
                <th>Kode</th>
                <th>Nama</th>
                <th>SKS</th>
                <th>Aksi</th>
            </tr>
        </thead>
        <tbody>
            @forelse ($daftarMatakuliah as $matakuliah)
                <tr>
                    <td>{{ $matakuliah['kode'] }}</td>
                    <td>{{ $matakuliah['nama'] }}</td>
                    <td>
                        <x-badge-sks :sks="$matakuliah['sks']" />
                    </td>
                    <td>
                        <a href="{{ route('matakuliah.show', $matakuliah['kode']) }}" class="btn btn-sm btn-primary">
                            Detail
                        </a>
                    </td>
                </tr>
            @empty
                <tr>
                    <td colspan="4">Data belum tersedia</td>
                </tr>
            @endforelse
        </tbody>
    </table>
@endsection
