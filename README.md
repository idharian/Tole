# Tole

<p align="center">
  <strong>Windows maintenance toolkit for cleanup, optimization, disk analysis, and app removal.</strong>
  <br />
  <sub>Native Go CLI for Windows power users.</sub>
  <br /><br />
  <a href="https://github.com/idharian/Tole/releases"><img src="https://img.shields.io/badge/version-v1.5.0-34D399?style=flat-square" alt="Version 1.5.0"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26"></a>
  <a href="https://github.com/idharian/Tole/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-94A3B8?style=flat-square" alt="MIT License"></a>
  <a href="https://github.com/idharian/Tole"><img src="https://img.shields.io/github/stars/idharian/Tole?style=flat-square" alt="GitHub stars"></a>
</p>

Tole is an all-in-one Windows system maintenance CLI. It scans disposable files, removes stale project artifacts, analyzes disk usage, monitors hardware, and handles application leftovers through a focused terminal interface.

Inspired by [Mole for macOS](https://github.com/tw93/mole), built natively for Windows in Go.

## Features

- **Deep cleanup**: user temp files, crash dumps, browser caches, developer caches, Windows Update cache, Prefetch, and Recycle Bin.
- **Project purge**: finds heavy artifacts such as `node_modules`, `target`, `.gradle`, `bin`, `obj`, `.venv`, `dist`, and `.next`.
- **Disk analyzer**: interactive folder size explorer with keyboard navigation.
- **System optimizer**: DNS flush, icon cache maintenance, memory trimming, DISM cleanup, and SSD TRIM where supported.
- **Smart uninstaller**: reads installed Windows applications and searches common leftover locations.
- **Live monitor**: CPU, RAM, disk, uptime, and system context dashboard.
- **Installer cleaner**: finds old `.msi`, setup `.exe`, `.iso`, and related packages in Downloads and Desktop.
- **Safe workflow**: dry-run previews, confirmation prompts, Recycle Bin support where appropriate, and locked-file handling.
- **Native Windows binary**: no runtime dependency after installation.

## Requirements

- Windows 10 or Windows 11.
- Go 1.26+ only when building from source.
- Administrator terminal for system-level cleanup and optimization tasks.

## Install

### Download a release

Download the latest Windows binary from [Releases](https://github.com/idharian/Tole/releases), then place `tole.exe` in a directory included in your `PATH`.

### Install from source

```powershell
git clone https://github.com/idharian/Tole.git
cd Tole
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

The installer copies `tole.exe` to `%LOCALAPPDATA%\Tole` and adds that directory to the user `PATH`. Open a new terminal after installation.

Verify:

```powershell
tole
```

Uninstall the installed binary and remove its `PATH` entry:

```powershell
powershell -ExecutionPolicy Bypass -File .\uninstall.ps1
```

## Usage

### Interactive menu

```powershell
tole
```

### Clean system caches

Preview first:

```powershell
tole clean --dry-run
```

Execute cleanup:

```powershell
tole clean
```

### Purge project artifacts

Preview a project tree:

```powershell
tole purge --path "D:\Projects" --dry-run
```

Scan and confirm interactively:

```powershell
tole purge --path "D:\Projects"
```

Skip confirmation:

```powershell
tole purge --path "D:\Projects" --yes
```

### Analyze disk usage

```powershell
tole analyze

tole analyze "C:\Users\musam\Downloads"
tole analyze "D:\"
```

Keys: `Up` / `Down` navigate, `Enter` open, `Esc` go to parent, `d` move selection to Recycle Bin, `q` exit.

### Optimize Windows

Preview:

```powershell
tole optimize --dry-run
```

Run maintenance:

```powershell
tole optimize
```

Run PowerShell as Administrator to enable system-level tasks such as DISM cleanup and SSD TRIM.

### Uninstall applications

List installed applications:

```powershell
tole uninstall --list
```

Search by application name:

```powershell
tole uninstall --search "photoshop"
```

Open interactive uninstaller:

```powershell
tole uninstall
```

### Monitor hardware

```powershell
tole status
```

Press `q` or `Esc` to exit.

### Clean old installer files

Preview:

```powershell
tole installer --dry-run
```

Move selected installer packages to Recycle Bin:

```powershell
tole installer
```

### Update otomatis dari GitHub Releases

Dari instalasi mana pun, jalankan:

```powershell
tole update
```

Tole akan mengambil release terbaru dari `idharian/Tole`, membandingkan versi, mengunduh `tole-windows-amd64.exe`, memverifikasi `SHA256SUMS`, lalu mengganti binary di `%LOCALAPPDATA%\Tole`. Tidak perlu masuk folder source, `git`, atau Go.

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

## Command reference

| Command | Purpose |
| --- | --- |
| `tole` | Open interactive menu |
| `tole clean` | Clean Windows, browser, and developer caches |
| `tole purge` | Find and remove project build artifacts |
| `tole analyze` | Explore disk usage interactively |
| `tole optimize` | Run Windows maintenance tasks |
| `tole uninstall` | Remove applications and leftovers |
| `tole status` | Show live CPU, RAM, and disk stats |
| `tole installer` | Find old installer packages |
| `tole update` | Pull, rebuild, and reinstall from source |

Most destructive commands support `--dry-run`. Use it before cleanup on an unfamiliar machine.

## Development

Clone the repository and run tests:

```powershell
git clone https://github.com/idharian/Tole.git
cd Tole
go test ./...
go build ./...
```

Build a Windows release binary:

```powershell
go build -ldflags="-s -w" -o tole.exe ./cmd/tole
```

Project layout:

```text
Tole/
├── cmd/tole/       CLI commands and entrypoint
├── pkg/cleaner/    cache scanning and cleanup
├── pkg/purge/      project artifact scanner
├── pkg/analyzer/   disk analyzer TUI
├── pkg/optimizer/  Windows maintenance tasks
├── pkg/uninstaller registry and leftover cleanup
├── pkg/installer/  installer package scanner
├── pkg/monitor/    live hardware monitor
├── pkg/ui/         Lipgloss theme, banner, spinner, summaries
├── pkg/winapi/     native Windows APIs
├── install.ps1     installer script
└── uninstall.ps1   uninstaller script
```

## Safety

Tole can permanently delete files during `tole clean`. Preview with `--dry-run` first. Close applications before cleaning their caches. Run elevated commands only when you understand the requested operation.

## License

MIT License. See [LICENSE](https://github.com/idharian/Tole/blob/main/LICENSE).

Maintained by [agushariyanto](https://github.com/idharian).
