package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/icono.png
var iconoBytes []byte

//go:embed assets/logo.png
var logoBytes []byte

// IconResource devuelve el icono de la aplicación (usado para el icono de la app y del tray).
func IconResource() fyne.Resource {
	return fyne.NewStaticResource("icono.png", iconoBytes)
}

// LogoResource devuelve el logo (usado en la ventana de prompt y configuración).
func LogoResource() fyne.Resource {
	return fyne.NewStaticResource("logo.png", logoBytes)
}
