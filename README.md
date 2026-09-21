# 🐹 Tole

<p align="center">
  <em>All-in-one Windows system maintenance, deep cleaning, and optimization toolkit.</em>
  <br />
  <strong>Inspired by <a href="https://github.com/tw93/mole">Mole for macOS</a>, crafted natively for Windows in Pure Go.</strong>
  <br />
  <sub>Dev by <strong>agushariyanto</strong></sub>
</p>

---

## ✨ Features

- **⚡ Blazing Fast Single Binary**: Ditulis murni dalam bahasa **Go** dengan integrasi langsung ke Win32 API. Tidak memerlukan runtime PowerShell yang lambat atau dependensi eksternal.
- **🎨 Modern Terminal Aesthetics**: Menggunakan styling TUI **Charmbracelet** (`lipgloss`, `bubbletea`, `huh`) dengan visual modern dan responsif.
- **🧹 Deep Clean (`tole clean`)**:
  - **User Essentials**: `%TEMP%`, Windows Error Reporting, Crash Dumps, dan Windows Recycle Bin.
  - **Browsers**: Google Chrome, Microsoft Edge, Brave Browser, dan Mozilla Firefox cache.
  - **Developer Tools**: npm cache, pnpm store, Yarn cache, Python pip cache, Go build cache, Rust Cargo cache, dan NuGet cache.
  - **System Caches** *(dengan hak Administrator)*: `C:\Windows\Temp`, Windows Prefetch, dan Windows Update download cache (`SoftwareDistribution\Download`).
- **📦 Project Build Purge (`tole purge`)**:
  - Memindai folder proyek developer untuk mencari direktori gemuk seperti `node_modules`, `target`, `.gradle`, `bin`/`obj`, `.venv`, `dist`, `.next`.
- **🛡️ Safety First**:
  - Dukungan penuh flag `--dry-run` di setiap perintah untuk melihat preview kalkulasi ruang sebelum dihapus.
  - Menangani *sharing violations* / file locking secara aman tanpa crash.
  - Integrasi native Recycle Bin Win32 (`SHEmptyRecycleBinW`).
- **🎯 Interactive CLI Menu**: Cukup ketik `tole` tanpa argumen untuk menampilkan menu pilihan interaktif.

---

## 🚀 Instalasi Cepat (One-Click Install)

Untuk memasang Tole secara permanen ke sistem Windows Anda agar dapat dipanggil dari folder mana saja:

```powershell
# Jalankan skrip instalasi otomatis
powershell -ExecutionPolicy Bypass -File .\install.ps1
```
*Skrip ini menyalin `tole.exe` ke `%LOCALAPPDATA%\Tole` dan otomatis mendaftarkannya ke variabel Environment `PATH` Windows.*

> Setelah instalasi, buka jendela terminal baru (PowerShell atau Command Prompt) dan Anda bisa langsung mengetik `tole` dari direktori manapun!

Untuk menghapus Tole dari sistem:
```powershell
powershell -ExecutionPolicy Bypass -File .\uninstall.ps1
```

---

## 📖 Panduan Penggunaan Lengkap

### 1. Menu Interaktif (Default)
Jika Anda hanya mengetik `tole` tanpa argumen tambahan, Tole akan menampilkan menu seleksi interaktif yang ramah pengguna:
```powershell
tole
```

### 2. Deep Cleaning (`tole clean`)
Membersihkan cache sementara, browser, dan developer tools:
```powershell
# Pratinjau kalkulasi tanpa menghapus berkas (Sangat Disarankan pertama kali)
tole clean --dry-run

# Melakukan pembersihan menyeluruh
tole clean
```

### 3. Pembersih Build Artifact Proyek (`tole purge`)
Memindai folder proyek developer untuk membersihkan `node_modules`, `target`, `.gradle`, dll.:
```powershell
# Scan folder proyek Anda (simulasi)
tole purge --path "D:\Projects" --dry-run

# Scan dan hapus dengan konfirmasi interaktif
tole purge --path "D:\Projects"

# Otomatis konfirmasi penghapusan tanpa prompt
tole purge --path "D:\Projects" --yes
```

### 4. Visualizer Ruang Disk Interaktif (`tole analyze`)
Melihat folder yang paling memakan ruang disk di terminal (seperti DaisyDisk / WizTree):
```powershell
# Menganalisis folder kerja saat ini
tole analyze

# Menganalisis folder tertentu (misal: Downloads atau Drive D:)
tole analyze "C:\Users\musam\Downloads"
tole analyze "D:\"
```
*Gunakan tombol panah `↑`/`↓` untuk memilih, `Enter` untuk masuk folder, `Esc` untuk kembali, `d` untuk membuang ke Recycle Bin, dan `q` untuk keluar.*

### 5. Optimasi Sistem Windows (`tole optimize`)
Membersihkan cache DNS, mereset icon cache yang rusak, dan memangkas memori standby RAM:
```powershell
# Simulasi pratinjau optimasi
tole optimize --dry-run

# Eksekusi pemeliharaan sistem
tole optimize
```
*(Buka terminal sebagai Administrator untuk membuka fitur tambahan pembersihan DISM dan SSD TRIM).*

### 6. Uninstaller Pintar & Pembersih Sisa (`tole uninstall`)
Mencopot aplikasi sekaligus membersihkan berkas yang tertinggal di AppData/ProgramData:
```powershell
# Menampilkan seluruh daftar software terinstal
tole uninstall --list

# Mencari aplikasi tertentu berdasarkan nama
tole uninstall --search "photoshop"

# Membuka menu uninstaller interaktif
tole uninstall
```

### 7. Pemantau Hardware Live (`tole status`)
Dashboard pemantau CPU, RAM, dan Disk secara real-time:
```powershell
tole status
```
*(Tekan `q` untuk keluar).*

### 8. Pembersih File Installer Lama (`tole installer`)
Mencari berkas `.msi`, `.exe` setup, dan `.iso` yang terlupakan di folder Downloads:
```powershell
# Pratinjau file installer yang ditemukan
tole installer --dry-run

# Scan dan pindahkan ke Recycle Bin
tole installer
```

### 9. Pembaruan Versi Otomatis (`tole update`)
Memperbarui Tole ke versi kompilasi terbaru secara instan:
```powershell
tole update
```

---

## 📋 Daftar Perintah

| Perintah | Deskripsi | Status |
| :--- | :--- | :---: |
| `tole` | Menu interaktif TUI | ✅ Tersedia |
| `tole clean` | Membersihkan cache sistem, browser, developer tools, dan Recycle Bin | ✅ Tersedia |
| `tole purge` | Mencari dan membersihkan build artifacts (`node_modules`, `target`, dll.) | ✅ Tersedia |
| `tole analyze` | Disk space visualizer interaktif (DaisyDisk / WizTree style TUI) | ✅ Tersedia |
| `tole optimize` | Perbaikan sistem: Flush DNS, trim RAM, icon cache rebuild | ✅ Tersedia |
| `tole uninstall` | Smart uninstaller dengan pembersih berkas sisa (leftover cleaner) | ✅ Tersedia |
| `tole status` | Dashboard pemantau hardware live (CPU, GPU, RAM, Disk, Network) | ✅ Tersedia |
| `tole installer` | Pembersih file installer `.msi`, `.exe`, `.iso` usang di folder Downloads | ✅ Tersedia |
| `tole update` | Memperbarui instalasi Tole ke build terbaru secara otomatis | ✅ Tersedia |

---

## 🏗️ Struktur Proyek

```
Tole/
├── cmd/
│   └── tole/               # CLI Entrypoint & Commands (Cobra)
│       ├── main.go         # Entry point
│       ├── root.go         # Interactive menu default (dry-run + confirm guard)
│       ├── clean.go        # tole clean & --dry-run (permanent delete warning)
│       ├── purge.go        # tole purge & scanner (sorted results)
│       ├── analyze.go      # tole analyze (Bubble Tea disk visualizer)
│       ├── optimize.go     # tole optimize & --dry-run
│       ├── uninstall.go    # tole uninstall (validated index, --list/--search)
│       ├── status.go       # tole status (live monitor)
│       ├── installer.go    # tole installer & --dry-run
│       └── update.go       # tole update (source-channel rebuild)
├── pkg/
│   ├── ui/                 # Styling Lipgloss, ASCII Banner, Spinner, Formatter
│   ├── winapi/             # Win32 APIs (Recycle Bin, Disk Space, Elevation, Paths)
│   ├── cleaner/            # Scanner & Engine Pembersih Cache
│   ├── purge/              # Engine Scanner Build Artifacts
│   ├── analyzer/           # Disk scanner + Bubble Tea TUI model
│   ├── optimizer/          # DNS flush, icon cache, RAM trim, DISM, TRIM
│   ├── uninstaller/        # Registry reader + strict leftover matcher
│   ├── installer/          # Downloads/Desktop installer scanner (recursive)
│   ├── monitor/            # CPU/RAM/Disk telemetry + TUI dashboard
│   └── version/            # CurrentVersion constant
├── go.mod
└── README.md
```

---

## 📄 Lisensi

MIT License. Bebas digunakan dan dikembangkan lebih lanjut.
