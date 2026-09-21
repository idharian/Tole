package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modKernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procGetDiskFreeSpaceExW = modKernel32.NewProc("GetDiskFreeSpaceExW")
)

type DiskSpaceInfo struct {
	FreeBytesAvailable uint64
	TotalNumberOfBytes uint64
	TotalNumberOfFreeBytes uint64
}

// GetDiskSpace returns available and total bytes for a drive (e.g. "C:\\")
func GetDiskSpace(drivePath string) (free uint64, total uint64, err error) {
	if drivePath == "" {
		drivePath = "C:\\"
	}
	pathPtr, err := syscall.UTF16PtrFromString(drivePath)
	if err != nil {
		return 0, 0, err
	}

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64

	ret, _, callErr := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if ret == 0 {
		return 0, 0, callErr
	}

	return freeBytesAvailable, totalNumberOfBytes, nil
}
