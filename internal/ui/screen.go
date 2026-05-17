package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// targetMaxSize estima el tamaño máximo (~1/3 del área de pantalla) al que
// puede crecer la ventana de prompt. Fyne no expone el tamaño físico de la
// pantalla en su API portátil, así que asumimos un monitor común de 1920x1080
// y permitimos crecer hasta cubrir aproximadamente 1/3 de su área manteniendo
// proporción cercana al diseño original.
func targetMaxSize() (float32, float32) {
	const (
		assumedW = float64(1920)
		assumedH = float64(1080)
	)
	area := assumedW * assumedH / 3
	ratio := float64(promptInitialW) / float64(promptInitialH)
	h := math.Sqrt(area / ratio)
	w := ratio * h
	return float32(w), float32(h)
}

func themeForeground() color.Color {
	return color.NRGBA{R: 0xea, G: 0xed, B: 0xf3, A: 0xff}
}

func themePlaceholder() color.Color {
	return color.NRGBA{R: 0x9a, G: 0xa0, B: 0xae, A: 0xff}
}

func localizeSubtitle(t *canvas.Text) {
	// Fyne usa locale del sistema; si el texto luce inglés, traducimos el día.
	repl := map[string]string{
		"Monday":    "Lunes",
		"Tuesday":   "Martes",
		"Wednesday": "Miércoles",
		"Thursday":  "Jueves",
		"Friday":    "Viernes",
		"Saturday":  "Sábado",
		"Sunday":    "Domingo",
	}
	for k, v := range repl {
		if len(t.Text) >= len(k) && t.Text[:len(k)] == k {
			t.Text = v + t.Text[len(k):]
			break
		}
	}
}

// minSize devuelve el menor de dos valores fyne.Size componente a componente.
func minSize(a, b fyne.Size) fyne.Size {
	return fyne.NewSize(minFloat(a.Width, b.Width), minFloat(a.Height, b.Height))
}

func minFloat(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
