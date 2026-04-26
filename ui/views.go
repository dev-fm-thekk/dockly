package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
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

	switch m.activeTab {
	case tabOverview:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}
		statusLabel := "System Status: OFFLINE"
		if m.metrics.Status {
			statusLabel = "System Status: ONLINE"
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
		)
		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(stats)

	case tabContainers:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}
		var list []string
		for _, c := range m.containers {
			statusIcon := "🔴"
			if strings.Contains(strings.ToLower(c.State), "running") {
				statusIcon = "🟢"
			}
			ports := ""
			if ports == "" {
				ports = "None"
			}
			list = append(list, fmt.Sprintf("%s %-20s | %-16s | Ports: %s", statusIcon, c.Name, c.Status, ports))
		}
		if len(list) == 0 {
			list = append(list, "No Containers Found.")
		}
		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render("Active Containers"),
				strings.Join(list, "\n\n"),
			),
		)

	case tabImages:
		if m.isLoading {
			content = activePaneStyle.Width(width - 4).Height(height - 4).Align(lipgloss.Center).Render("\n\nLoading Docker API Data...")
			break
		}
		var list []string
		for _, img := range m.images {
			sizeMb := float64(img.Size) / (1024 * 1024)
			createdTime := time.Unix(img.Created, 0).Format(time.RFC822)
			list = append(list, fmt.Sprintf("%-25s : %-15s | %7.2f MB | %s", img.Repository, img.Tag, sizeMb, createdTime))
		}
		if len(list) == 0 {
			list = append(list, "No Local Images Found.")
		}
		content = activePaneStyle.Width(width - 4).Height(height - 4).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render("Local Images"),
				strings.Join(list, "\n\n"),
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
