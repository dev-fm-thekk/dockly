package main

import (
	"fmt"
	"os"

	"dockly/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
