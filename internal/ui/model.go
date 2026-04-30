package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	dockerapi "dockly/internal/docker-api"
)

const (
	tabOverview = iota
	tabContainers
	tabImages
	tabNetwork
)

type apiDataMsg struct {
	metrics    dockerapi.Metrics
	containers []dockerapi.Container
	images     []dockerapi.Image
	network    []dockerapi.NetworkConfig
}

func fetchData() tea.Msg {
	return apiDataMsg{
		metrics:    dockerapi.FetchMetrics(),
		containers: dockerapi.FetchContainers(),
		images:     dockerapi.FetchImages(),
		network:    dockerapi.FetchNetConfigAll(),
	}
}

type Model struct {
	width     int
	height    int
	activeTab int
	tabs      []string

	metrics    dockerapi.Metrics
	containers []dockerapi.Container
	images     []dockerapi.Image
	network    []dockerapi.NetworkConfig
	isLoading  bool

	scrollOffsets map[int]int
	cursor        map[int]int // Track current selected row per tab
}

func NewModel() Model {
	return Model{
		activeTab:     tabOverview,
		tabs:          []string{"OVERVIEW", "CONTAINERS", "IMAGES", "NETWORKS"},
		isLoading:     true,
		scrollOffsets: make(map[int]int),
		cursor:        make(map[int]int),
	}
}

func (m Model) Init() tea.Cmd {
	return fetchData
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case apiDataMsg:
		m.metrics = msg.metrics
		m.containers = msg.containers
		m.images = msg.images
		m.network = msg.network
		m.isLoading = false
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		case "shift+tab", "left", "h":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		case "up", "k":
			limit := 0
			switch m.activeTab {
			case tabContainers:
				limit = len(m.containers)
			case tabImages:
				limit = len(m.images)
			case tabNetwork:
				limit = len(m.network)
			}

			if limit > 0 {
				if m.cursor[m.activeTab] == 0 {
					m.cursor[m.activeTab] = limit - 1 // Wrap to bottom
				} else {
					m.cursor[m.activeTab]--
				}
			}

		case "down", "j":
			limit := 0
			switch m.activeTab {
			case tabContainers:
				limit = len(m.containers)
			case tabImages:
				limit = len(m.images)
			case tabNetwork:
				limit = len(m.network)
			}

			if limit > 0 {
				if m.cursor[m.activeTab] == limit-1 {
					m.cursor[m.activeTab] = 0 // Wrap to top
				} else {
					m.cursor[m.activeTab]++
				}
			}
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Initializing...")
	}

	tabs := m.renderedTabs(m.width)
	footerText := " tab • h/l navigation  |  j/k select  |  enter details  |  q quit "
	footer := footerStyle.Width(m.width).Align(lipgloss.Right).Render(footerText)

	contentHeight := m.height - lipgloss.Height(tabs) - lipgloss.Height(footer) - 1
	content := m.renderedContent(m.width, contentHeight)

	ui := lipgloss.JoinVertical(lipgloss.Left,
		tabs,
		content,
		footer,
	)

	return tea.NewView(ui)
}
