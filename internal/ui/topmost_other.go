//go:build !windows

package ui

import "fyne.io/fyne/v2"

// En plataformas no-Windows no aplicamos always-on-top a nivel SO desde
// código portable. En macOS el system tray (menubar) cumple el rol de
// presencia persistente; en Linux depende del WM.

func enableAlwaysOnTop(_ fyne.Window)  {}
func disableAlwaysOnTop(_ fyne.Window) {}
