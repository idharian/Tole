# Tole

<p align="center">
  <strong>Toolkit pemeliharaan Windows untuk pembersihan, optimasi, analisis disk, dan penghapusan aplikasi.</strong>
  <br />
  <sub>CLI native Go untuk power user Windows.</sub>
  <br /><br />
  <a href="https://github.com/idharian/Tole/releases"><img src="https://img.shields.io/badge/versi-v1.5.0-34D399?style=flat-square" alt="Versi 1.5.0"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26"></a>
  <a href="https://github.com/idharian/Tole/blob/main/LICENSE"><img src="https://img.shields.io/badge/lisensi-MIT-94A3B8?style=flat-square" alt="Lisensi MIT"></a>
  <a href="https://github.com/idharian/Tole"><img src="https://img.shields.io/github/stars/idharian/Tole?style=flat-square" alt="GitHub stars"></a>
  <a href="https://github.com/idharian/Tole/releases"><img src="https://img.shields.io/github/v/release/idharian/Tole?style=flat-square" alt="Release terbaru"></a>
  <a href="https://github.com/idharian/Tole/actions/workflows/release.yml"><img src="https://img.shields.io/github/actions/workflow/status/idharian/Tole/release.yml?style=flat-square&label=release" alt="Status release"></a>
  <a href="https://github.com/idharian/Tole/releases"><img src="https://img.shields.io/github/downloads/idharian/Tole/total?style=flat-square" alt="Total unduhan"></a>
  <a href="https://github.com/idharian/Tole/pulls"><img src="https://img.shields.io/badge/PR%20diterima-ya-34D399?style=flat-square" alt="PR diterima"></a>
  <a href="https://github.com/idharian/Tole/issues"><img src="https://img.shields.io/github/issues/idharian/Tole?style=flat-square" alt="Isu terbuka"></a>
</p>

Tole adalah CLI pemeliharaan sistem Windows all-in-one. Ia memindai file sekali pakai, menghapus artifact proyek yang sudah basi, menganalisis penggunaan disk, memantau hardware, dan menangani sisa-sisa aplikasi melalui antarmuka terminal yang fokus dan rapi.

Terinspirasi [Mole untuk macOS](https://github.com/tw93/mole), dibangun native untuk Windows dengan Go.

## Fitur

- **Deep clean**: file temp user, crash dump, cache browser, cache developer, cache Windows Update, Prefetch, dan Recycle Bin.
- **Project purge**: menemukan artifact berat seperti `node_modules`, `target`, `.gradle`, `bin`, `obj`, `.venv`, `dist`, dan `.next`.
- **Disk analyzer**: penjelajah ukuran folder interaktif dengan navigasi keyboard.
- **Large file finder**: daftar file terbesar dalam satu tampilan, hapus lewat Recycle Bin.
- **System optimizer**: flush DNS, perawatan icon cache, trimming memory, DISM cleanup, dan SSD TRIM bila didukung.
- **Smart uninstaller**: membaca aplikasi Windows yang terpasang dan mencari lokasi sisa-sisanya.
- **Live monitor**: CPU, RAM, disk, uptime, dan dashboard konteks sistem.
- **Installer cleaner**: menemukan `.msi` lama, setup `.exe`, `.iso`, dan paket terkait di Downloads dan Desktop.
- **Alur aman**: preview dry-run, prompt konfirmasi, dukungan Recycle Bin bila sesuai, dan penanganan file terkunci.
- **Binary Windows native**: tanpa dependensi runtime setelah terpasang.

## Persyaratan

- Windows 10 atau Windows 11.
- Go 1.26+ hanya jika build dari source.
- Terminal Administrator untuk tugas pembersihan dan optimasi level sistem.

## Instalasi

### Unduh release

Unduh binary Windows terbaru dari [Releases](https://github.com/idharian/Tole/releases), lalu letakkan `tole.exe` di direktori yang ada di `PATH` kamu.

### Instal dari source

```powershell
git clone https://github.com/idharian/Tole.git
cd Tole
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

Installer menyalin `tole.exe` ke `%LOCALAPPDATA%\Tole` dan menambahkan direktori tersebut ke `PATH` user. Buka terminal baru setelah instalasi.

Verifikasi:

```powershell
tole version
```

Hapus binary terpasang dan entri `PATH`-nya:

```powershell
powershell -ExecutionPolicy Bypass -File .\uninstall.ps1
```

## Penggunaan

### Menu interaktif

```powershell
tole
```

### Bersihkan cache sistem

Preview dulu:

```powershell
tole clean --dry-run
```

Jalankan pembersihan:

```powershell
tole clean
```

### Bersihkan artifact proyek

Preview pohon proyek:

```powershell
tole purge --path "D:\Projects" --dry-run
```

Scan dan konfirmasi interaktif:

```powershell
tole purge --path "D:\Projects"
```

Lewati konfirmasi:

```powershell
tole purge --path "D:\Projects" --yes
```

### Analisis penggunaan disk

```powershell
tole analyze

tole analyze "C:\Users\musam\Downloads"
tole analyze "D:\"
```

Tombol: `Up` / `Down` navigasi, `Enter` masuk, `Esc` kembali ke folder induk, `d` pindahkan seleksi ke Recycle Bin, `q` keluar.

### Cari file terbesar

Preview file terbesar di suatu path:

```powershell
tole large --dry-run
tole large "D:\Projects" --min 500MB --top 50
```

Pindahkan ke Recycle Bin (masih bisa dipulihkan):

```powershell
tole large "D:\Projects"
```

Flag `--min` menerima `KB`, `MB`, `GB`, `TB`, atau angka byte polos. Default `--min` adalah `100MB`, default `--top` adalah `25`.

### Optimasi Windows

Preview:

```powershell
tole optimize --dry-run
```

Jalankan pemeliharaan:

```powershell
tole optimize
```

Jalankan PowerShell sebagai Administrator untuk mengaktifkan tugas level sistem seperti DISM cleanup dan SSD TRIM.

### Uninstall aplikasi

Daftar aplikasi terpasang:

```powershell
tole uninstall --list
```

Cari berdasarkan nama aplikasi:

```powershell
tole uninstall --search "photoshop"
```

Buka uninstaller interaktif:

```powershell
tole uninstall
```

### Pantau hardware

```powershell
tole status
```

Tekan `q` atau `Esc` untuk keluar.

### Bersihkan file installer lama

Preview:

```powershell
tole installer --dry-run
```

Pindahkan paket installer terpilih ke Recycle Bin:

```powershell
tole installer
```

### Update otomatis dari GitHub Releases

Dari instalasi mana pun, jalankan:

```powershell
tole update
```

Tole akan mengambil release terbaru dari `idharian/Tole`, membandingkan versi, mengunduh `tole-windows-amd64.exe`, memverifikasi `SHA256SUMS`, lalu mengganti binary di `%LOCALAPPDATA%\Tole`. Binary baru diverifikasi bisa dijalankan sebelum backup lama dihapus — jika health check gagal, versi lama dipulihkan otomatis. Tidak perlu masuk folder source, `git`, atau Go.

Release baru dibuat otomatis saat maintainer push tag versi:

```powershell
git tag v1.6.0
git push origin v1.6.0
```

### Update dari source repository

Untuk build dari source secara manual:

```powershell
cd Tole
git pull
go build -o .\tole.exe .\cmd\tole
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

## Referensi perintah

| Perintah | Fungsi |
| --- | --- |
| `tole` | Buka menu interaktif |
| `tole clean` | Bersihkan cache Windows, browser, dan developer |
| `tole purge` | Temukan dan hapus artifact build proyek |
| `tole analyze` | Jelajahi penggunaan disk secara interaktif |
| `tole large` | Temukan dan hapus file terbesar |
| `tole optimize` | Jalankan tugas pemeliharaan Windows |
| `tole uninstall` | Hapus aplikasi dan sisa-sisanya |
| `tole status` | Tampilkan statistik live CPU, RAM, dan disk |
| `tole installer` | Cari paket installer lama |
| `tole version` | Tampilkan versi Tole |
| `tole update` | Ambil dan pasang rilis terbaru dari GitHub |

Sebagian besar perintah destruktif mendukung `--dry-run`. Gunakan dulu sebelum cleanup di mesin yang belum kamu kenal.

## Pengembangan

Clone repositori dan jalankan test:

```powershell
git clone https://github.com/idharian/Tole.git
cd Tole
go test ./...
go build ./...
```

Build binary release Windows:

```powershell
go build -ldflags="-s -w" -o tole.exe ./cmd/tole
```

Struktur proyek:

```text
Tole/
├── cmd/tole/       Perintah CLI dan entrypoint
├── pkg/cleaner/    pemindaian dan pembersihan cache
├── pkg/purge/      pemindai artifact proyek
├── pkg/analyzer/   TUI disk analyzer
├── pkg/largefiles/ pencari file terbesar
├── pkg/optimizer/  tugas pemeliharaan Windows
├── pkg/uninstaller registry dan pembersihan sisa aplikasi
├── pkg/installer/  pemindai paket installer
├── pkg/monitor/    monitor hardware live
├── pkg/ui/         tema Lipgloss, banner, spinner, ringkasan
├── pkg/winapi/     API Windows native
├── install.ps1     script installer
└── uninstall.ps1   script uninstaller
```

## Keamanan

Tole dapat menghapus file secara permanen saat `tole clean`. Preview dulu dengan `--dry-run`. Tutup aplikasi sebelum membersihkan cache-nya. Jalankan perintah elevated hanya jika kamu memahami operasinya.

## Lisensi

Lisensi MIT. Lihat [LICENSE](https://github.com/idharian/Tole/blob/main/LICENSE).

Dikelola oleh [agushariyanto](https://github.com/idharian).
