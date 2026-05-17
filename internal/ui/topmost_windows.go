//go:build windows

package ui

import (
	"syscall"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
)

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	procFindWindowW           = user32.NewProc("FindWindowW")
	procSetWindowPos          = user32.NewProc("SetWindowPos")
	procEnumWindows           = user32.NewProc("EnumWindows")
	procGetWindowTextW        = user32.NewProc("GetWindowTextW")
	procGetWindowThreadProcID = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible       = user32.NewProc("IsWindowVisible")
	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentProcessID   = kernel32.NewProc("GetCurrentProcessId")
)

const (
	hwndTopmost   = ^uintptr(0)     // (HWND)-1
	hwndNoTopmost = ^uintptr(0) - 1 // (HWND)-2
	swpNoMove     = 0x0002
	swpNoSize     = 0x0001
	swpShowWindow = 0x0040
	swpNoActivate = 0x0010
)

func enableAlwaysOnTop(w fyne.Window) {
	// Buscamos el HWND de la ventana de Fyne reintentado durante un par de
	// segundos: Fyne crea la ventana asíncronamente, así que el primer
	// intento puede no encontrarla.
	go func() {
		title := w.Title()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			hwnd := findWindowByTitleInProcess(title)
			if hwnd != 0 {
				procSetWindowPos.Call(
					hwnd, hwndTopmost,
					0, 0, 0, 0,
					swpNoMove|swpNoSize|swpShowWindow|swpNoActivate,
				)
				return
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()
}

func disableAlwaysOnTop(w fyne.Window) {
	hwnd := findWindowByTitleInProcess(w.Title())
	if hwnd == 0 {
		return
	}
	procSetWindowPos.Call(
		hwnd, hwndNoTopmost,
		0, 0, 0, 0,
		swpNoMove|swpNoSize|swpShowWindow|swpNoActivate,
	)
}

// findWindowByTitleInProcess intenta primero un FindWindow rápido por título
// (clase nula). Si encuentra más de una ventana con ese título en otros
// procesos, hace fallback a EnumWindows filtrando por el PID actual.
func findWindowByTitleInProcess(title string) uintptr {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd != 0 && belongsToCurrentProcess(hwnd) {
		return hwnd
	}
	return enumFind(title)
}

func belongsToCurrentProcess(hwnd uintptr) bool {
	pid, _, _ := procGetCurrentProcessID.Call()
	var winPID uint32
	procGetWindowThreadProcID.Call(hwnd, uintptr(unsafe.Pointer(&winPID)))
	return uintptr(winPID) == pid
}

func enumFind(title string) uintptr {
	currentPID, _, _ := procGetCurrentProcessID.Call()
	var found uintptr
	wantTitleUtf16 := utf16FromString(title)

	cb := syscall.NewCallback(func(hwnd uintptr, _ uintptr) uintptr {
		var winPID uint32
		procGetWindowThreadProcID.Call(hwnd, uintptr(unsafe.Pointer(&winPID)))
		if uintptr(winPID) != currentPID {
			return 1
		}
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}
		buf := make([]uint16, 256)
		n, _, _ := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 256)
		if int(n) == 0 {
			return 1
		}
		if utf16Equal(buf[:n], wantTitleUtf16) {
			found = hwnd
			return 0
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)
	return found
}

func utf16FromString(s string) []uint16 {
	p, _ := syscall.UTF16FromString(s)
	if len(p) > 0 && p[len(p)-1] == 0 {
		p = p[:len(p)-1]
	}
	return p
}

func utf16Equal(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
