@extends('layouts.app')
@section('judul', 'IPK Tertinggi - Teknik Komputer')
@section('konten')
    <h1 class="h3 mb-4">10 Mahasiswa IPK Tertinggi - Teknik Komputer</h1>
    <table class="table table-striped bg-white">
        <thead>
            <tr>
                <th>#</th>
                <th>NIM</th>
                <th>Nama</th>
                <th>Angkatan</th>
                <th>IPK</th>
            </tr>
        </thead>
        <tbody>
            @foreach ($daftarMahasiswa as $index => $mahasiswa)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $mahasiswa->nim }}</td>
                    <td>{{ $mahasiswa->nama }}</td>
                    <td>{{ $mahasiswa->angkatan }}</td>
                    <td>{{ $mahasiswa->ipk }}</td>
                </tr>
            @endforeach
        </tbody>
    </table>
@endsection
