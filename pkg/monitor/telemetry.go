package monitor

import (
	"os"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"tole/pkg/winapi"
)

var (
	modKernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procGlobalMemoryStatusEx = modKernel32.NewProc("GlobalMemoryStatusEx")
	procGetSystemTimes       = modKernel32.NewProc("GetSystemTimes")
	procGetTickCount64       = modKernel32.NewProc("GetTickCount64")
)

type MEMORYSTATUSEX struct {
	DwLength                uint32
	DwMemoryLoad            uint32
	UllTotalPhys            uint64
	UllAvailPhys            uint64
	UllTotalPageFile        uint64
	UllAvailPageFile        uint64
	UllTotalVirtual         uint64
	UllAvailVirtual         uint64
	UllAvailExtendedVirtual uint64
}

type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

func (ft *FILETIME) ToUint64() uint64 {
	return (uint64(ft.DwHighDateTime) << 32) | uint64(ft.DwLowDateTime)
}

type SystemStats struct {
	CPUPercent   float64
	TotalRAM     uint64
	UsedRAM      uint64
	FreeRAM      uint64
	RAMPercent   float64
	TotalDisk    uint64
	FreeDisk     uint64
	UsedDisk     uint64
	Hostname     string
	OS           string
	Arch         string
	NumCPU       int
	NumGoroutine int
	Uptime       time.Duration
}

type Monitor struct {
	mu             sync.Mutex
	lastIdle       uint64
	lastKernel     uint64
	lastUser       uint64
	lastSampleTime time.Time
}

func NewMonitor() *Monitor {
	m := &Monitor{}
	m.sampleCPU()
	return m
}

func (m *Monitor) sampleCPU() float64 {
	var idleTime, kernelTime, userTime FILETIME

	ret, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	if ret == 0 {
		return 0
	}

	idle := idleTime.ToUint64()
	kernel := kernelTime.ToUint64()
	user := userTime.ToUint64()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.lastKernel == 0 {
		m.lastIdle = idle
		m.lastKernel = kernel
		m.lastUser = user
		m.lastSampleTime = time.Now()
		return 0
	}

	idleDelta := idle - m.lastIdle
	kernelDelta := kernel - m.lastKernel
	userDelta := user - m.lastUser

	m.lastIdle = idle
	m.lastKernel = kernel
	m.lastUser = user
	m.lastSampleTime = time.Now()

	totalSystem := kernelDelta + userDelta
	if totalSystem == 0 {
		return 0
	}

	// kernelTime includes idleTime
	cpuBusy := totalSystem - idleDelta
	percent := (float64(cpuBusy) / float64(totalSystem)) * 100.0
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	return percent
}

func (m *Monitor) GetStats() SystemStats {
	cpu := m.sampleCPU()

	// Memory
	var memStatus MEMORYSTATUSEX
	memStatus.DwLength = uint32(unsafe.Sizeof(memStatus))
	procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))

	totalRAM := memStatus.UllTotalPhys
	freeRAM := memStatus.UllAvailPhys
	var usedRAM uint64
	if totalRAM >= freeRAM {
		usedRAM = totalRAM - freeRAM
	} else {
		totalRAM, freeRAM = 0, 0
	}
	var ramPercent float64
	if totalRAM > 0 {
		ramPercent = (float64(usedRAM) / float64(totalRAM)) * 100.0
	}

	// Disk
	freeDisk, totalDisk, diskErr := winapi.GetDiskSpace("C:\\")
	var usedDisk uint64
	if diskErr == nil && totalDisk >= freeDisk {
		usedDisk = totalDisk - freeDisk
	} else {
		freeDisk, totalDisk, usedDisk = 0, 0, 0
	}

	hostname, _ := os.Hostname()

	retTicks, _, _ := procGetTickCount64.Call()
	uptime := time.Duration(retTicks) * time.Millisecond

	return SystemStats{
		CPUPercent:   cpu,
		TotalRAM:     totalRAM,
		UsedRAM:      usedRAM,
		FreeRAM:      freeRAM,
		RAMPercent:   ramPercent,
		TotalDisk:    totalDisk,
		FreeDisk:     freeDisk,
		UsedDisk:     usedDisk,
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		Uptime:       uptime,
	}
}
