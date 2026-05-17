//go:build windows

package ui

import (
	"time"

	"fyne.io/fyne/v2"
)

var procGetSystemMetrics = user32.NewProc("GetSystemMetrics")

const (
	smCxScreen = 0
	smCyScreen = 1
)

// placeTopRight mueve la ventana a la esquina superior-derecha del monitor
// primario, con un pequeño margen, y la fija como topmost.
func placeTopRight(w fyne.Window, width, height int) {
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			hwnd := findWindowByTitleInProcess(w.Title())
			if hwnd != 0 {
				cx, _, _ := procGetSystemMetrics.Call(smCxScreen)
				margin := uintptr(16)
				x := cx - uintptr(width) - margin
				y := margin
				procSetWindowPos.Call(
					hwnd, hwndTopmost,
					x, y, uintptr(width), uintptr(height),
					swpShowWindow|swpNoActivate,
				)
				return
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()
}
