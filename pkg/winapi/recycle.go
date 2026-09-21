package winapi

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modShell32               = windows.NewLazySystemDLL("shell32.dll")
	procSHQueryRecycleBinW   = modShell32.NewProc("SHQueryRecycleBinW")
	procSHEmptyRecycleBinW   = modShell32.NewProc("SHEmptyRecycleBinW")
)

type SHQUERYRBINFO struct {
	CbSize      uint32
	I64Size     int64
	I64NumItems int64
}

const (
	SHERB_NOCONFIRMATION = 0x00000001
	SHERB_NOPROGRESSUI   = 0x00000002
	SHERB_NOSOUND        = 0x00000004
)

// GetRecycleBinInfo returns the total size in bytes and number of items in the Recycle Bin.
func GetRecycleBinInfo() (int64, int64, error) {
	var rbInfo SHQUERYRBINFO
	rbInfo.CbSize = uint32(unsafe.Sizeof(rbInfo))

	ret, _, err := procSHQueryRecycleBinW.Call(
		0, // null root path checks all drives
		uintptr(unsafe.Pointer(&rbInfo)),
	)
	if ret != 0 {
		return 0, 0, err
	}

	return rbInfo.I64Size, rbInfo.I64NumItems, nil
}

// EmptyRecycleBin empties the Windows Recycle Bin on all drives.
func EmptyRecycleBin(noSound bool) error {
	var flags uintptr = SHERB_NOCONFIRMATION | SHERB_NOPROGRESSUI
	if noSound {
		flags |= SHERB_NOSOUND
	}

	ret, _, err := procSHEmptyRecycleBinW.Call(
		0, // hwnd
		0, // pszRootPath (null means all drives)
		flags,
	)
	// S_OK is 0
	if ret != 0 {
		// S_FALSE (0x00000001) or other codes
		if ret == 1 {
			return nil
		}
		return err
	}
	return nil
}

// MoveToRecycleBin moves a file or directory to the Recycle Bin using SHFileOperationW
func MoveToRecycleBin(path string) error {
	modShell := windows.NewLazySystemDLL("shell32.dll")
	procSHFileOperation := modShell.NewProc("SHFileOperationW")

	const (
		FO_DELETE          = 3
		FOF_ALLOWUNDO      = 0x0040
		FOF_NOCONFIRMATION = 0x0010
		FOF_SILENT         = 0x0004
	)

	type SHFILEOPSTRUCTW struct {
		Hwnd                  uintptr
		WFunc                 uint32
		PFrom                 *uint16
		PTo                   *uint16
		FFlags                uint16
		FAnyOperationsAborted int32
		HNameMappings         uintptr
		LpszProgressTitle     *uint16
	}

	// pFrom must be double null-terminated
	utf16Path, err := syscall.UTF16FromString(path)
	if err != nil {
		return err
	}
	utf16Path = append(utf16Path, 0) // double null terminator

	op := SHFILEOPSTRUCTW{
		WFunc:  FO_DELETE,
		PFrom:  &utf16Path[0],
		FFlags: FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_SILENT,
	}

	ret, _, _ := procSHFileOperation.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}
