package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/edwin/timelog/internal/autostart"
	"github.com/edwin/timelog/internal/config"
	"github.com/edwin/timelog/internal/scheduler"
	"github.com/edwin/timelog/internal/storage"
	"github.com/edwin/timelog/internal/ui"
)

const appID = "com.edwin.timelog"

func main() {
	a := app.NewWithID(appID)
	a.Settings().SetTheme(&ui.ModernTheme{})
	a.SetIcon(ui.IconResource())

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("cargar configuración: %v", err)
	}

	// Si en disco está marcado autostart pero el SO no lo tiene, reconciliamos.
	if cfg.Autostart && !autostart.IsEnabled() {
		_ = autostart.Enable()
	}

	store := storage.New()

	promptCtl := ui.NewPromptController(a, cfg, store)
	showPrompt := func() { fyne.Do(promptCtl.Show) }

	sch := scheduler.New(time.Duration(cfg.IntervalMinutes)*time.Minute, showPrompt)
	settingsCtl := ui.NewSettingsController(a, cfg, func() {
		snap := cfg.Snapshot()
		sch.SetInterval(time.Duration(snap.IntervalMinutes) * time.Minute)
	})

	showSettings := func() { fyne.Do(settingsCtl.Show) }

	quit := func() {
		sch.Stop()
		a.Quit()
	}

	hasTray := ui.SetupTray(a, showPrompt, showSettings, quit)

	floating := ui.FloatingWidget(a, showPrompt, showSettings, func() {
		// El botón "✕" del widget flotante también solo oculta; salir es desde el tray.
		floatingHide(a)
	})
	if hasTray {
		// Con tray disponible, mostramos el flotante pequeño como recordatorio visual.
		floating.Show()
	} else {
		// Sin tray (ej. Linux sin StatusNotifier), el flotante es la única vía.
		floating.Show()
	}

	sch.Start()

	// Lanzamos la ventana de configuración la primera vez para que el usuario
	// confirme intervalo, ubicación del CSV y autoarranque.
	if isFirstRun(cfg) {
		go func() {
			time.Sleep(300 * time.Millisecond)
			fyne.Do(settingsCtl.Show)
		}()
	}

	a.Run()
}

// floatingHide oculta todas las ventanas (usado por el botón ✕ del widget flotante).
// La app sigue viva en el system tray.
func floatingHide(a fyne.App) {
	for _, w := range a.Driver().AllWindows() {
		w.Hide()
	}
}

// isFirstRun heurística: si nunca se han registrado clientes y el intervalo
// es el default, asumimos primera ejecución y abrimos config.
func isFirstRun(cfg *config.Config) bool {
	snap := cfg.Snapshot()
	return len(snap.Clients) == 0
}
