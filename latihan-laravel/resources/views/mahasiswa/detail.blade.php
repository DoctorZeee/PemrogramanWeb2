@extends('layouts.app')
@section('judul', 'Detail Mahasiswa')
@section('konten')
    <a href="{{ route('mahasiswa.data') }}" class="btn btn-sm btn-outline-secondary mb-3">&laquo; Kembali</a>

    <div class="card mb-4">
        <div class="card-body">
            <h1 class="h4 mb-3">{{ $mahasiswa->nama }}</h1>
            <table class="table table-borderless mb-0">
                <tr>
                    <th style="width: 160px;">NIM</th>
                    <td>{{ $mahasiswa->nim }}</td>
                </tr>
                <tr>
                    <th>Email</th>
                    <td>{{ $mahasiswa->email }}</td>
                </tr>
                <tr>
                    <th>Program Studi</th>
                    <td>{{ $mahasiswa->programStudi->nama }}</td>
                </tr>
                <tr>
                    <th>Angkatan</th>
                    <td>{{ $mahasiswa->angkatan }}</td>
                </tr>
                <tr>
                    <th>IPK</th>
                    <td>{{ $mahasiswa->ipk }}</td>
                </tr>
            </table>
        </div>
    </div>

    <h2 class="h5 mb-3">Daftar Matakuliah</h2>
    <table class="table table-striped bg-white">
        <thead>
            <tr>
                <th>Kode</th>
                <th>Nama Matakuliah</th>
                <th>SKS</th>
                <th>Semester</th>
                <th>Nilai</th>
            </tr>
        </thead>
        <tbody>
            @forelse ($mahasiswa->matakuliah as $matakuliah)
                <tr>
                    <td>{{ $matakuliah->kode }}</td>
                    <td>{{ $matakuliah->nama }}</td>
                    <td>{{ $matakuliah->sks }}</td>
                    <td>{{ $matakuliah->semester }}</td>
                    <td>{{ $matakuliah->pivot->nilai }}</td>
                </tr>
            @empty
                <tr>
                    <td colspan="5" class="text-center">Belum ada matakuliah yang diambil.</td>
                </tr>
            @endforelse
        </tbody>
    </table>
@endsection
