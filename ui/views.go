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
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

func (m Model) renderedTabs(width int) string {
	var renderedTabs []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(t))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, row)
}

func (m Model) renderedContent(width, height int) string {
	var content string

	// Usable height inside the pane borders and padding
	// Header takes 1 line, spacer 1 line.
	// Total overhead: 2 (borders) + 2 (padding) + 1 (header) + 1 (spacer) = 6 lines
	maxVisibleItems := height - 10 // table rows are 1 line each

	switch m.activeTab {
	case tabOverview:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}
		statusLabel := "System Status: ONLINE"
		if !m.metrics.Status {
			statusLabel = "System Status: OFFLINE"
		}
		memUsedGb := float64(m.metrics.Memory.Used) / (1024 * 1024 * 1024)
		memTotalGb := float64(m.metrics.Memory.Total) / (1024 * 1024 * 1024)
		memPerc := float64(0)
		if memTotalGb > 0 {
			memPerc = (memUsedGb / memTotalGb) * 100.0
		}
		cpuBar := createProgressBar(m.metrics.CPU.Percentage, 20)
		memBar := createProgressBar(memPerc, 20)

		stats := lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render(statusLabel),
			"",
			fmt.Sprintf("CPU Usage:    %s %.2f%%", cpuBar, m.metrics.CPU.Percentage),
			fmt.Sprintf("Memory Usage: %s %.2fGB / %.2fGB", memBar, memUsedGb, memTotalGb),
			"",
			fmt.Sprintf("Active Containers: %d", m.metrics.ActiveContainers),
			fmt.Sprintf("Total Images:      %d", m.metrics.Images),
			"",
			lipgloss.NewStyle().Foreground(colorSubtext).Render("Last updated: "+time.Now().Format("15:04:05")),
		)
		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(stats)

	case tabContainers:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}

		headers := []string{"STATUS", "NAME", "IMAGE", "UPTIME", "STATE"}
		var rows [][]string
		for _, c := range m.containers {
			statusIcon := "🔴"
			if strings.Contains(strings.ToLower(c.State), "running") {
				statusIcon = "🟢"
			}
			rows = append(rows, []string{
				statusIcon,
				c.Name,
				c.Image,
				c.Status,
				c.State,
			})
		}

		// Calculate scroll
		offset := m.scrollOffsets[tabContainers]
		if offset > len(rows)-maxVisibleItems && len(rows) > maxVisibleItems {
			offset = len(rows) - maxVisibleItems
		}
		if offset < 0 {
			offset = 0
		}
		visibleRows := rows[offset:min(offset+maxVisibleItems, len(rows))]

		t := table.New().
			Headers(headers...).
			Rows(visibleRows...).
			BaseStyle(rowStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 { // header
					return columnHeaderStyle
				}
				return rowStyle
			})

		scrollInfo := ""
		if len(rows) > maxVisibleItems {
			scrollInfo = fmt.Sprintf(" (%d-%d of %d)", offset+1, min(offset+maxVisibleItems, len(rows)), len(rows))
		}

		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render("Active Containers"+scrollInfo),
				t.String(),
			),
		)

	case tabImages:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}

		headers := []string{"REPOSITORY", "TAG", "ID", "SIZE", "CREATED"}
		var rows [][]string
		for _, img := range m.images {
			sizeMb := fmt.Sprintf("%.2f MB", float64(img.Size)/(1024*1024))
			createdTime := time.Unix(img.Created, 0).Format("2006-01-02")
			rows = append(rows, []string{
				img.Repository,
				img.Tag,
				img.ID,
				sizeMb,
				createdTime,
			})
		}

		// Calculate scroll
		offset := m.scrollOffsets[tabImages]
		if offset > len(rows)-maxVisibleItems && len(rows) > maxVisibleItems {
			offset = len(rows) - maxVisibleItems
		}
		if offset < 0 {
			offset = 0
		}
		visibleRows := rows[offset:min(offset+maxVisibleItems, len(rows))]

		t := table.New().
			Headers(headers...).
			Rows(visibleRows...).
			BaseStyle(rowStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 { // header
					return columnHeaderStyle
				}
				return rowStyle
			})

		scrollInfo := ""
		if len(rows) > maxVisibleItems {
			scrollInfo = fmt.Sprintf(" (%d-%d of %d)", offset+1, min(offset+maxVisibleItems, len(rows)), len(rows))
		}

		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render("Local Images"+scrollInfo),
				t.String(),
			),
		)

	case tabSettings:
		settings := lipgloss.JoinVertical(lipgloss.Left,
			headerStyle.Render("Dashboard Preferences"),
			"",
			"[x] Enable Dark Mode",
			"[ ] Auto-refresh stats",
			"[x] Show exiting containers",
			"",
			"Press 'space' to toggle (Placeholder)",
		)
		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(settings)
	}

	return content
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
