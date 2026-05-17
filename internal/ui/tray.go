package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// SetupTray instala el menú del system tray (menubar en macOS, system tray en
// Windows/Linux). Devuelve true si el driver soporta tray.
func SetupTray(app fyne.App, onPrompt, onSettings, onQuit func()) bool {
	desk, ok := app.(desktop.App)
	if !ok {
		return false
	}

	addItem := fyne.NewMenuItem("+ Nueva tarea", onPrompt)
	addItem.Icon = IconResource()
	menu := fyne.NewMenu("TimeLog",
		addItem,
		fyne.NewMenuItem("Configuración…", onSettings),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Salir", onQuit),
	)
	desk.SetSystemTrayMenu(menu)
	desk.SetSystemTrayIcon(IconResource())
	return true
}
