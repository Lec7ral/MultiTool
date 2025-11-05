package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme defines the custom color scheme for the application
type CustomTheme struct {
	baseTheme fyne.Theme
}

// Color returns the color for a specific name and variant
func (c *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 45, G: 110, B: 175, A: 255} // #2D6EAF
	case theme.ColorNameBackground:
		return color.NRGBA{R: 33, G: 33, B: 33, A: 255} // #212121
	case theme.ColorNameForeground:
		return color.NRGBA{R: 248, G: 249, B: 250, A: 255} // #F8F9FA
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 44, G: 44, B: 44, A: 255} // #2C2C2C
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 44, G: 44, B: 44, A: 255} // #2C2C2C
	default:
		return theme.DarkTheme().Color(name, variant)
	}
}

// Font returns the font for a specific style
func (c *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

// Icon returns the icon for a specific name
func (c *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

// Size returns the size for a specific name
func (c *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}

// ApplyTheme sets the custom theme for the application
func ApplyTheme(app fyne.App) {
	app.Settings().SetTheme(&CustomTheme{baseTheme: theme.DarkTheme()})
}
