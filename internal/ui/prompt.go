package ui

import (
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/edwin/timelog/internal/config"
	"github.com/edwin/timelog/internal/storage"
)

const (
	promptInitialW = 460
	promptInitialH = 380
	growStep       = 30 * time.Second
	growRatio      = 1.0 / 3.0
)

// PromptResult representa la respuesta del usuario.
type PromptResult struct {
	Saved       bool
	Client      string
	Activity    string
	Description string
}

// PromptController gestiona la ventana de captura.
type PromptController struct {
	app    fyne.App
	cfg    *config.Config
	store  *storage.CSVStore
	mu     sync.Mutex
	active fyne.Window
}

func NewPromptController(app fyne.App, cfg *config.Config, store *storage.CSVStore) *PromptController {
	return &PromptController{app: app, cfg: cfg, store: store}
}

// Show abre la ventana de prompt si no hay otra activa. Es thread-safe y debe llamarse desde la goroutine principal de Fyne.
func (p *PromptController) Show() {
	p.mu.Lock()
	if p.active != nil {
		p.active.RequestFocus()
		p.active.Show()
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	snap := p.cfg.Snapshot()

	w := p.app.NewWindow("¿Qué estás haciendo?")
	w.SetIcon(LogoResource())

	logo := canvas.NewImageFromResource(LogoResource())
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(96, 96))

	title := canvas.NewText("Registro de actividad", themeForeground())
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = 22

	subtitle := canvas.NewText(time.Now().Format("Lunes 02 Ene 2006 · 15:04"), themePlaceholder())
	subtitle.TextSize = 12
	localizeSubtitle(subtitle)

	clientEntry := widget.NewEntry()
	clientEntry.SetPlaceHolder("Nombre del cliente")
	clientSelect := widget.NewSelect(append([]string{""}, snap.Clients...), func(v string) {
		if v != "" {
			clientEntry.SetText(v)
		}
	})
	clientSelect.PlaceHolder = "Recientes…"

	activitySelect := widget.NewSelect(snap.Activities, func(string) {})
	activitySelect.PlaceHolder = "Tipo de actividad"

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("¿Qué estás haciendo? (descripción breve)")
	descEntry.Wrapping = fyne.TextWrapWord
	descEntry.SetMinRowsVisible(4)

	form := container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("Cliente"), nil,
			container.NewBorder(nil, nil, nil, clientSelect, clientEntry),
		),
		container.NewBorder(nil, nil, widget.NewLabel("Actividad"), nil, activitySelect),
		widget.NewLabel("Descripción"),
		descEntry,
	)

	status := widget.NewLabel("")

	var saveBtn, laterBtn *widget.Button
	saveBtn = widget.NewButton("Guardar", func() {
		client := strings.TrimSpace(clientEntry.Text)
		activity := activitySelect.Selected
		desc := strings.TrimSpace(descEntry.Text)
		if client == "" || activity == "" || desc == "" {
			status.SetText("Cliente, actividad y descripción son obligatorios.")
			return
		}
		entry := storage.Entry{
			Timestamp:   time.Now(),
			Client:      client,
			Activity:    activity,
			Description: desc,
		}
		if err := p.store.Append(snap.CSVPath, entry); err != nil {
			dialog.ShowError(err, w)
			return
		}
		_ = p.cfg.AddClient(client)
		p.close(w)
	})
	saveBtn.Importance = widget.HighImportance
	laterBtn = widget.NewButton("Recordar luego", func() {
		p.close(w)
	})

	header := container.NewBorder(nil, nil, container.NewPadded(logo), nil,
		container.NewVBox(title, subtitle),
	)

	content := container.NewBorder(
		container.NewPadded(header),
		container.NewPadded(container.NewBorder(nil, nil, nil,
			container.NewHBox(laterBtn, saveBtn), status,
		)),
		nil, nil,
		container.NewPadded(form),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(promptInitialW, promptInitialH))
	w.CenterOnScreen()
	w.SetCloseIntercept(func() { p.close(w) })

	p.mu.Lock()
	p.active = w
	p.mu.Unlock()

	w.Show()
	w.RequestFocus()
	w.Canvas().Focus(clientEntry)

	go p.growLoop(w)
}

func (p *PromptController) close(w fyne.Window) {
	p.mu.Lock()
	if p.active == w {
		p.active = nil
	}
	p.mu.Unlock()
	w.Hide()
	w.Close()
}

// growLoop hace crecer la ventana progresivamente hasta ~1/3 de pantalla
// si el usuario no la cierra. La ventana toma foco periódicamente.
func (p *PromptController) growLoop(w fyne.Window) {
	maxW, maxH := targetMaxSize()
	curW := float32(promptInitialW)
	curH := float32(promptInitialH)

	step := float32(0)
	for {
		time.Sleep(growStep)

		p.mu.Lock()
		alive := p.active == w
		p.mu.Unlock()
		if !alive {
			return
		}

		step++
		// crece ~12% por paso, con techo en maxW/maxH (~1/3 de pantalla)
		factor := float32(1 + 0.12*float64(step))
		nW := float32(promptInitialW) * factor
		nH := float32(promptInitialH) * factor
		if nW > maxW {
			nW = maxW
		}
		if nH > maxH {
			nH = maxH
		}

		if nW == curW && nH == curH {
			// ya tope, sigue tomando foco
			w.RequestFocus()
			continue
		}
		curW, curH = nW, nH
		w.Resize(fyne.NewSize(curW, curH))
		w.CenterOnScreen()
		w.RequestFocus()
	}
}
