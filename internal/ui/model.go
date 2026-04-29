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
	tabAccount
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

// Model represents the top-level Bubble Tea model for the dashboard.
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
}

func NewModel() Model {
	return Model{
		activeTab:     tabOverview,
		tabs:          []string{"Overview", "Containers", "Images", "Network"},
		isLoading:     true,
		scrollOffsets: make(map[int]int),
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
			if m.scrollOffsets[m.activeTab] > 0 {
				m.scrollOffsets[m.activeTab]--
			}
		case "down", "j":
			m.scrollOffsets[m.activeTab]++
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Initializing...") // Wait for first WindowSizeMsg
	}

	title := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, titleStyle.Render("Dockly"))

	tabs := m.renderedTabs(m.width)

	footerText := "↑/k up • ↓/j down • →/l/tab next • ←/h/shift+tab prev • q quit"
	footer := footerStyle.Width(m.width).Align(lipgloss.Center).Render(footerText)

	// Calculate available height: Total - Title - Tabs - Footer - Padding
	availableHeight := m.height - lipgloss.Height(title) - lipgloss.Height(tabs) - lipgloss.Height(footer) - 2
	if availableHeight < 5 {
		return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Terminal size too small"))
	}

	content := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.renderedContent(m.width, availableHeight))

	ui := lipgloss.JoinVertical(lipgloss.Top,
		title,
		tabs,
		content,
		footer,
	)

	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui))
}
