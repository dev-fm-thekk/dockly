package ui

import "charm.land/lipgloss/v2"

// Styling Definitions
var (
	colorPrimary   = lipgloss.Color("#bd93f9") // Purple
	colorSecondary = lipgloss.Color("#50fa7b") // Green
	colorAccent    = lipgloss.Color("#ff79c6") // Pink
	colorText      = lipgloss.Color("#f8f8f2") // White/Light Gray
	colorSubtext   = lipgloss.Color("#6272a4") // Gray

	baseStyle = lipgloss.NewStyle().Foreground(colorText)

	headerStyle = baseStyle.
			Bold(true).
			Foreground(colorSecondary).
			MarginBottom(1)

	tabStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(colorSubtext).
			Padding(0, 2).
			Foreground(colorSubtext)

	activeTabStyle = tabStyle.
			BorderForeground(colorPrimary).
			Foreground(colorPrimary).
			Bold(true)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(colorSubtext).
			Padding(1, 2)

	activePaneStyle = paneStyle.
			BorderForeground(colorAccent)

	titleStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Padding(1, 4).
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(colorPrimary)

	footerStyle = baseStyle.
			Foreground(colorSubtext).
			MarginTop(1)

	columnHeaderStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(colorSubtext).
				Padding(0, 1)

	rowStyle = lipgloss.NewStyle().
			Padding(0, 1)
)
