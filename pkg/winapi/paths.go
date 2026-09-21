package winapi

import (
	"os"
	"path/filepath"
)

// Common Windows directories
func GetUserTempDir() string {
	if t := os.Getenv("TEMP"); t != "" {
		return t
	}
	if t := os.Getenv("TMP"); t != "" {
		return t
	}
	return filepath.Join(GetUserProfileDir(), "AppData", "Local", "Temp")
}

func GetSystemTempDir() string {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	return filepath.Join(windir, "Temp")
}

func GetWindowsPrefetchDir() string {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	return filepath.Join(windir, "Prefetch")
}

func GetSoftwareDistributionDownloadDir() string {
	windir := os.Getenv("WINDIR")
	if windir == "" {
		windir = `C:\Windows`
	}
	return filepath.Join(windir, "SoftwareDistribution", "Download")
}

func GetLocalAppDataDir() string {
	if dir := os.Getenv("LOCALAPPDATA"); dir != "" {
		return dir
	}
	return filepath.Join(GetUserProfileDir(), "AppData", "Local")
}

func GetRoamingAppDataDir() string {
	if dir := os.Getenv("APPDATA"); dir != "" {
		return dir
	}
	return filepath.Join(GetUserProfileDir(), "AppData", "Roaming")
}

func GetUserProfileDir() string {
	if dir := os.Getenv("USERPROFILE"); dir != "" {
		return dir
	}
	return "C:\\Users\\Default"
}

func GetProgramDataDir() string {
	if dir := os.Getenv("ProgramData"); dir != "" {
		return dir
	}
	return `C:\ProgramData`
}
