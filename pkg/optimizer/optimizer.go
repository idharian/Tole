package optimizer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"tole/pkg/winapi"
	"unsafe"

	"golang.org/x/sys/windows"
)

type OptimizeTask struct {
	ID          string
	Name        string
	Description string
	AdminOnly   bool
	Action      func(dryRun bool) (string, error)
}

func GetTasks() []OptimizeTask {
	return []OptimizeTask{
		{
			ID:          "flush_dns",
			Name:        "Flush DNS Resolver Cache",
			Description: "Clears local DNS cache to resolve hostname and connectivity issues",
			AdminOnly:   false,
			Action:      ActionFlushDNS,
		},
		{
			ID:          "icon_cache",
			Name:        "Rebuild Icon & Thumbnail Cache",
			Description: "Removes corrupted icon and thumbnail databases",
			AdminOnly:   false,
			Action:      ActionRebuildIconCache,
		},
		{
			ID:          "trim_memory",
			Name:        "Trim Process Memory Working Sets",
			Description: "Releases unused process memory back to the Windows memory manager",
			AdminOnly:   false,
			Action:      ActionTrimMemory,
		},
		{
			ID:          "dism_cleanup",
			Name:        "Windows Component Store Cleanup (DISM)",
			Description: "Cleans up superseded Windows update components",
			AdminOnly:   true,
			Action:      ActionDISMCleanup,
		},
		{
			ID:          "storage_trim",
			Name:        "SSD Storage TRIM & Defrag",
			Description: "Sends TRIM hints to the SSD to optimize write performance",
			AdminOnly:   true,
			Action:      ActionStorageTrim,
		},
	}
}

// ActionFlushDNS flushes Windows DNS resolver cache.
func ActionFlushDNS(dryRun bool) (string, error) {
	if dryRun {
		return "Would execute 'ipconfig /flushdns'", nil
	}

	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}
	return "DNS resolver cache successfully flushed", nil
}

// ActionRebuildIconCache removes IconCache.db and Explorer thumbnail caches.
func ActionRebuildIconCache(dryRun bool) (string, error) {
	localAppData := winapi.GetLocalAppDataDir()
	iconCachePath := filepath.Join(localAppData, "IconCache.db")
	explorerCacheDir := filepath.Join(localAppData, "Microsoft", "Windows", "Explorer")

	if dryRun {
		return fmt.Sprintf("Would clear '%s' and thumbnail caches in '%s'", iconCachePath, explorerCacheDir), nil
	}

	var removedCount int

	// Remove IconCache.db if exists
	if err := os.Remove(iconCachePath); err == nil {
		removedCount++
	}

	// Remove iconcache_*.db (Windows 10/11 per-user icon caches)
	matches, _ := filepath.Glob(filepath.Join(localAppData, "iconcache_*.db"))
	for _, m := range matches {
		if rerr := os.Remove(m); rerr == nil {
			removedCount++
		}
	}

	// Remove thumbcache files
	entries, err := os.ReadDir(explorerCacheDir)
	if err == nil {
		for _, entry := range entries {
			name := entry.Name()
			if len(name) >= 10 && name[:10] == "thumbcache" {
				full := filepath.Join(explorerCacheDir, name)
				if rerr := os.Remove(full); rerr == nil {
					removedCount++
				}
			}
		}
	}

	return fmt.Sprintf("Reset %d icon/thumbnail database items", removedCount), nil
}

// ActionTrimMemory calls EmptyWorkingSet on processes.
func ActionTrimMemory(dryRun bool) (string, error) {
	if dryRun {
		return "Would call EmptyWorkingSet across user processes", nil
	}

	modPsapi := windows.NewLazySystemDLL("psapi.dll")
	procEmptyWorkingSet := modPsapi.NewProc("EmptyWorkingSet")

	// Open process snapshot
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	trimmedCount := 0
	if err := windows.Process32First(snapshot, &entry); err == nil {
		for {
			pid := entry.ProcessID
			if pid > 4 {
				hProc, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_QUERY_INFORMATION, false, pid)
				if err == nil {
					ret, _, _ := procEmptyWorkingSet.Call(uintptr(hProc))
					if ret != 0 {
						trimmedCount++
					}
					windows.CloseHandle(hProc)
				}
			}
			if err := windows.Process32Next(snapshot, &entry); err != nil {
				break
			}
		}
	}

	return fmt.Sprintf("Trimmed working sets for %d active processes", trimmedCount), nil
}

// ActionDISMCleanup executes DISM component cleanup.
func ActionDISMCleanup(dryRun bool) (string, error) {
	if !winapi.IsAdmin() {
		return "", fmt.Errorf("requires Administrator privileges")
	}

	if dryRun {
		return "Would execute 'dism.exe /Online /Cleanup-Image /StartComponentCleanup'", nil
	}

	cmd := exec.Command("dism.exe", "/Online", "/Cleanup-Image", "/StartComponentCleanup")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}

	return "Windows component store cleanup completed successfully", nil
}

// ActionStorageTrim sends TRIM command to SSD.
func ActionStorageTrim(dryRun bool) (string, error) {
	if !winapi.IsAdmin() {
		return "", fmt.Errorf("requires Administrator privileges")
	}

	if dryRun {
		return "Would execute 'defrag.exe C: /O' (TRIM)", nil
	}

	cmd := exec.Command("defrag.exe", "C:", "/O")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, string(out))
	}

	return "Storage TRIM completed on drive C:", nil
}
