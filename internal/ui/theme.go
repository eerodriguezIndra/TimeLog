package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ModernTheme proporciona una paleta moderna oscura con acentos azules.
type ModernTheme struct{}

var _ fyne.Theme = (*ModernTheme)(nil)

func (m *ModernTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x15, G: 0x17, B: 0x21, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xea, G: 0xed, B: 0xf3, A: 0xff}
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x1f, G: 0x23, B: 0x2e, A: 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x2a, G: 0x2f, B: 0x3d, A: 0xff}
	case theme.ColorNamePrimary, theme.ColorNameFocus:
		return color.NRGBA{R: 0x4e, G: 0x8c, B: 0xff, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x4a, G: 0x4e, B: 0x5a, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x70, G: 0x76, B: 0x84, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x80}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x2a, G: 0x2f, B: 0x3d, A: 0xff}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x4e, G: 0x8c, B: 0xff, A: 0x88}
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (m *ModernTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m *ModernTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *ModernTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 1
	case theme.SizeNameScrollBar:
		return 12
	}
	return theme.DefaultTheme().Size(name)
}
