package ui

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/edwin/timelog/internal/autostart"
	"github.com/edwin/timelog/internal/config"
)

// SettingsController controla la ventana de configuración.
type SettingsController struct {
	app          fyne.App
	cfg          *config.Config
	onChange     func()
	mu           sync.Mutex
	active       fyne.Window
}

func NewSettingsController(app fyne.App, cfg *config.Config, onChange func()) *SettingsController {
	return &SettingsController{app: app, cfg: cfg, onChange: onChange}
}

func (s *SettingsController) Show() {
	s.mu.Lock()
	if s.active != nil {
		s.active.RequestFocus()
		s.active.Show()
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	snap := s.cfg.Snapshot()

	w := s.app.NewWindow("Configuración · TimeLog")
	w.SetIcon(LogoResource())

	logo := canvas.NewImageFromResource(LogoResource())
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(80, 80))

	title := canvas.NewText("Preferencias", themeForeground())
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 20
	subtitle := canvas.NewText("Personaliza el recordatorio y la ubicación de tus registros", themePlaceholder())
	subtitle.TextSize = 12

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(strconv.Itoa(snap.IntervalMinutes))
	intervalEntry.Validator = func(v string) error {
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || n <= 0 {
			return fmt.Errorf("debe ser un número entero mayor que 0")
		}
		return nil
	}

	pathEntry := widget.NewEntry()
	pathEntry.SetText(snap.CSVPath)

	browseBtn := widget.NewButton("Elegir…", func() {
		dlg := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
			if err != nil || uri == nil {
				return
			}
			defer uri.Close()
			pathEntry.SetText(uri.URI().Path())
		}, w)
		dlg.SetFileName("timelog.csv")
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		dlg.Resize(fyne.NewSize(720, 520))
		dlg.Show()
	})

	autostartCheck := widget.NewCheck("Iniciar al iniciar sesión", nil)
	autostartCheck.SetChecked(snap.Autostart || autostart.IsEnabled())

	status := widget.NewLabel("")

	saveBtn := widget.NewButton("Guardar cambios", func() {
		if err := intervalEntry.Validate(); err != nil {
			status.SetText(err.Error())
			return
		}
		minutes, _ := strconv.Atoi(strings.TrimSpace(intervalEntry.Text))
		newPath := strings.TrimSpace(pathEntry.Text)
		if newPath == "" {
			status.SetText("La ruta del CSV no puede estar vacía.")
			return
		}
		wantAuto := autostartCheck.Checked

		err := s.cfg.Update(func(cc *config.Config) {
			cc.IntervalMinutes = minutes
			cc.CSVPath = newPath
			cc.Autostart = wantAuto
		})
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Aplicar autostart al sistema operativo
		if wantAuto && !autostart.IsEnabled() {
			if err := autostart.Enable(); err != nil {
				dialog.ShowError(fmt.Errorf("no se pudo activar el autoarranque: %w", err), w)
			}
		} else if !wantAuto && autostart.IsEnabled() {
			if err := autostart.Disable(); err != nil {
				dialog.ShowError(fmt.Errorf("no se pudo desactivar el autoarranque: %w", err), w)
			}
		}

		if s.onChange != nil {
			s.onChange()
		}
		status.SetText("Cambios guardados ✓")
		go func() {
			time.Sleep(2 * time.Second)
			fyne.Do(func() { status.SetText("") })
		}()
	})
	saveBtn.Importance = widget.HighImportance

	form := container.NewVBox(
		labeledRow("Cada cuánto preguntar (minutos)", intervalEntry),
		labeledRow("Archivo CSV de registro", container.NewBorder(nil, nil, nil, browseBtn, pathEntry)),
		container.NewPadded(autostartCheck),
	)

	header := container.NewBorder(nil, nil, container.NewPadded(logo), nil,
		container.NewVBox(title, subtitle),
	)

	bottom := container.NewBorder(nil, nil, status, saveBtn)

	content := container.NewBorder(
		container.NewPadded(header),
		container.NewPadded(bottom),
		nil, nil,
		container.NewPadded(form),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 380))
	w.CenterOnScreen()
	w.SetCloseIntercept(func() {
		s.mu.Lock()
		s.active = nil
		s.mu.Unlock()
		w.Hide()
	})

	s.mu.Lock()
	s.active = w
	s.mu.Unlock()

	w.Show()
}

func labeledRow(label string, control fyne.CanvasObject) *fyne.Container {
	lbl := widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return container.NewVBox(lbl, control)
}
