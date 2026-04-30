package ui

import "charm.land/lipgloss/v2"

var (
	// High Contrast Black & Blue Palette
	colorBlue      = lipgloss.Color("#00AFFF") // Electric Blue
	colorBlueDim   = lipgloss.Color("#005F87") // Deep Blue for borders/zebra
	colorBlack     = lipgloss.Color("#000000") // True Black
	colorGrayDark  = lipgloss.Color("#121212") // Near Black for zebra stripes
	colorGrayLight = lipgloss.Color("#585858") // Subtext
	colorWhite     = lipgloss.Color("#FFFFFF") // Primary Text

	baseStyle = lipgloss.NewStyle().Foreground(colorWhite)

	headerStyle = baseStyle.
			Bold(true).
			Foreground(colorBlue).
			MarginBottom(1)

	// Tab Row Style - Full width blue bar
	tabRowStyle = lipgloss.NewStyle().
			Background(colorBlue).
			Foreground(colorBlack)

	tabStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Foreground(colorBlack).
			Background(colorBlue)

	activeTabStyle = tabStyle.
			Bold(true).
			Background(colorWhite).
			Foreground(colorBlack)

	paneStyle = lipgloss.NewStyle().
			Padding(0, 1)

	activePaneStyle = paneStyle

	titleStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true).
			Padding(0, 1)

	footerStyle = baseStyle.
			Foreground(colorGrayLight).
			MarginTop(1)

	// Table specific styles
	tableTitleStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true)

	columnHeaderStyle = lipgloss.NewStyle().
				Foreground(colorWhite).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(colorBlueDim).
				Padding(0, 1)

	rowStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(colorWhite)

	rowAltStyle = rowStyle.
			Background(colorGrayDark).
			Foreground(colorGrayLight)

	selectedRowStyle = rowStyle.
				Background(colorBlueDim).
				Foreground(colorWhite).
				Bold(true)
)
