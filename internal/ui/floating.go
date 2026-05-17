package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FloatingWidget crea la mini-ventana siempre visible con el botón principal
// "+ Nueva tarea" — el usuario puede registrar una entrada en cualquier
// momento sin esperar al recordatorio periódico. La X del SO solo oculta;
// la app sigue viva en el system tray.
func FloatingWidget(app fyne.App, onPrompt, onSettings, onHide func()) fyne.Window {
	w := app.NewWindow("TimeLog")
	w.SetIcon(IconResource())

	// Marca de la app: ícono + nombre
	logo := canvas.NewImageFromResource(IconResource())
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(32, 32))

	brand := canvas.NewText("TimeLog", themeForeground())
	brand.TextStyle = fyne.TextStyle{Bold: true}
	brand.TextSize = 14

	header := container.NewHBox(logo, brand)

	// Acción principal — registrar una tarea ahora, a voluntad
	addBtn := widget.NewButtonWithIcon("Nueva tarea", theme.ContentAddIcon(), onPrompt)
	addBtn.Importance = widget.HighImportance

	// Acciones secundarias en una fila compacta
	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), onSettings)
	settingsBtn.Importance = widget.LowImportance

	hideBtn := widget.NewButtonWithIcon("", theme.WindowMinimizeIcon(), onHide)
	hideBtn.Importance = widget.LowImportance

	actions := container.NewBorder(nil, nil, nil, container.NewHBox(settingsBtn, hideBtn))

	content := container.NewPadded(container.NewVBox(
		header,
		addBtn,
		actions,
	))
	w.SetContent(content)

	const widW, widH = 220, 150
	w.Resize(fyne.NewSize(widW, widH))
	w.CenterOnScreen()

	// La X del SO solo oculta — para salir realmente usar el tray.
	w.SetCloseIntercept(func() { w.Hide() })

	// En Windows: HWND_TOPMOST + anclar en esquina superior derecha.
	enableAlwaysOnTop(w)
	placeTopRight(w, widW, widH)

	return w
}
