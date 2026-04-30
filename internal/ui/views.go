package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func createProgressBar(percentage float64, length int) string {
	filled := int((percentage / 100.0) * float64(length))
	if filled > length {
		filled = length
	} else if filled < 0 {
		filled = 0
	}
	empty := length - filled
	return "[" + lipgloss.NewStyle().Foreground(colorBlue).Render(strings.Repeat("█", filled)) + strings.Repeat("░", empty) + "]"
}

func (m Model) renderedTabs(width int) string {
	var renderedTabs []string
	for i, t := range m.tabs {
		style := tabStyle
		if i == m.activeTab {
			style = activeTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, tabRowStyle.Width(width).Render(row))
}

func (m Model) renderedContent(width, height int) string {
	var content string
	maxVisibleItems := height - 6
	if maxVisibleItems < 1 {
		maxVisibleItems = 1
	}

	renderTable := func(title string, headers []string, rows [][]string, tabIndex int) string {
		numRows := len(rows)
		cursor := m.cursor[tabIndex]

		// Handle Automatic Scrolling: the selection drives the view
		offset := m.scrollOffsets[tabIndex]
		if cursor < offset {
			offset = cursor
		} else if cursor >= offset+maxVisibleItems {
			offset = cursor - maxVisibleItems + 1
		}

		// Update model's stored offset
		m.scrollOffsets[tabIndex] = offset

		end := offset + maxVisibleItems
		if end > numRows {
			end = numRows
		}

		var visibleRows [][]string
		if numRows > 0 && offset < numRows {
			visibleRows = rows[offset:end]
		}

		t := table.New().
			Headers(headers...).
			Rows(visibleRows...).
			Width(width - 4).
			BaseStyle(rowStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				// absoluteIndex translates the table's 1-based data rows into 0-based array indices.
				// If row is 1 (first data row) and offset is 0: (1 - 1) + 0 = 0 (First Element).

				// Now we compare our calculated data index against the Model's cursor.
				if row == cursor {
					return selectedRowStyle
				}

				// Apply zebra striping based on the absolute position in the full list.
				if row%2 == 0 {
					return rowAltStyle
				}
				return rowStyle
			})

		scrollInfo := ""
		if numRows > 0 {
			// Shows current position as 1-based (e.g., 1/10 instead of 0/10)
			scrollInfo = fmt.Sprintf(" %d/%d ", cursor+1, numRows)
		}

		titleBar := tableTitleStyle.Width(width - 4).Render(
			lipgloss.JoinHorizontal(lipgloss.Center,
				strings.Repeat("─", max(2, (width-len(title)-len(scrollInfo)-10)/2)),
				" "+title+scrollInfo+" ",
				strings.Repeat("─", max(2, (width-len(title)-len(scrollInfo)-10)/2)),
			),
		)

		return activePaneStyle.Width(width - 2).Height(height).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleBar,
				t.String(),
			),
		)
	}

	switch m.activeTab {
	case tabOverview:
		if m.isLoading {
			content = activePaneStyle.Width(width - 2).Height(height).Align(lipgloss.Center).Render("\n\nLoading...")
			break
		}

		memUsedGb := float64(m.metrics.Memory.Used) / (1024 * 1024 * 1024)
		memTotalGb := float64(m.metrics.Memory.Total) / (1024 * 1024 * 1024)
		memPerc := (memUsedGb / memTotalGb) * 100.0

		// Explicit white style to override terminal defaults
		whiteText := lipgloss.NewStyle().Foreground(colorWhite)

		// Updated bar helper to handle the blue/white combo strictly
		createBlueWhiteBar := func(percentage float64, length int) string {
			filledCount := int((percentage / 100.0) * float64(length))
			if filledCount > length {
				filledCount = length
			}
			if filledCount < 0 {
				filledCount = 0
			}

			filled := lipgloss.NewStyle().Foreground(colorBlue).Render(strings.Repeat("█", filledCount))
			empty := whiteText.Render(strings.Repeat("░", length-filledCount))

			return whiteText.Render("[") + filled + empty + whiteText.Render("]")
		}

		stats := lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render("SYSTEM OVERVIEW"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Left,
				whiteText.Render("CPU Usage:    "),
				createBlueWhiteBar(m.metrics.CPU.Percentage, 30),
				whiteText.Render(fmt.Sprintf(" %.2f%%", m.metrics.CPU.Percentage)),
			),
			lipgloss.JoinHorizontal(lipgloss.Left,
				whiteText.Render("Memory Usage: "),
				createBlueWhiteBar(memPerc, 30),
				whiteText.Render(fmt.Sprintf(" %.2fGB / %.2fGB", memUsedGb, memTotalGb)),
			),
			"",
			lipgloss.JoinHorizontal(lipgloss.Left,
				whiteText.Render("Active Containers: "),
				whiteText.Render(fmt.Sprintf("%d", m.metrics.ActiveContainers)),
			),
			lipgloss.JoinHorizontal(lipgloss.Left,
				whiteText.Render("Total Images:      "),
				whiteText.Render(fmt.Sprintf("%d", m.metrics.Images)),
			),
			"",
			lipgloss.NewStyle().Foreground(colorGrayLight).Render("Last updated: "+time.Now().Format("15:04:05")),
		)
		content = activePaneStyle.Width(width - 2).Height(height).Render(stats)

	case tabContainers:
		headers := []string{"STATUS", "NAME", "IMAGE", "STATE"}
		var rows [][]string
		for _, c := range m.containers {
			status := "exited"
			if strings.Contains(strings.ToLower(c.State), "running") {
				status = "running"
			}
			rows = append(rows, []string{status, c.Name, c.Image, c.State})
		}
		content = renderTable("CONTAINERS", headers, rows, tabContainers)

	case tabImages:
		headers := []string{"REPOSITORY", "TAG", "SIZE", "CREATED"}
		var rows [][]string
		for _, img := range m.images {
			sizeMb := fmt.Sprintf("%.2f MB", float64(img.Size)/(1024*1024))
			created := time.Unix(img.Created, 0).Format("2006-01-02")
			rows = append(rows, []string{img.Repository, img.Tag, sizeMb, created})
		}
		content = renderTable("IMAGES", headers, rows, tabImages)

	case tabNetwork:
		headers := []string{"ID", "NAME", "HOST", "PORT"}
		var rows [][]string
		for _, n := range m.network {
			rows = append(rows, []string{n.Container.ID[:12], n.Container.Name, n.Host, n.Port})
		}
		content = renderTable("NETWORKS", headers, rows, tabNetwork)
	}

	return content
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
