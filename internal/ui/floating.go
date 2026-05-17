package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// FloatingWidget crea la ventana pequeña con el icono, que el usuario puede
// dejar a un lado de la pantalla. La ventana NO tiene botón de cierre real:
// al intentar cerrarse, simplemente se oculta y queda accesible desde el tray.
//
// Notas multiplataforma: el comportamiento "siempre encima" depende del
// gestor de ventanas. En macOS y la mayoría de WMs Linux el tamaño compacto
// y la ausencia de decoración la mantiene como overlay liviano; en Windows
// la ventana respeta el orden Z estándar. El tray es la fuente principal de
// presencia persistente.
func FloatingWidget(app fyne.App, onClick, onSettings, onQuit func()) fyne.Window {
	w := app.NewWindow("TimeLog")
	w.SetIcon(IconResource())
	w.SetPadded(false)

	icon := canvas.NewImageFromResource(IconResource())
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(64, 64))

	iconBtn := widget.NewButton("", onClick)
	iconBtn.SetIcon(IconResource())
	iconBtn.Importance = widget.LowImportance

	settingsBtn := widget.NewButton("⚙", onSettings)
	settingsBtn.Importance = widget.LowImportance

	quitBtn := widget.NewButton("✕", onQuit)
	quitBtn.Importance = widget.LowImportance

	row := container.NewHBox(iconBtn, settingsBtn, quitBtn)
	w.SetContent(container.NewPadded(row))
	w.Resize(fyne.NewSize(180, 80))
	w.CenterOnScreen()

	// El "botón de cierre" del SO solo oculta la ventana — la app sigue viva en el tray.
	w.SetCloseIntercept(func() { w.Hide() })

	return w
}
