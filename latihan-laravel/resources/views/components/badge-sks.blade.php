@props(['sks'])

<span class="badge {{ $sks >= 3 ? 'bg-primary' : 'bg-warning text-dark' }}">
    {{ $sks }} SKS
</span>