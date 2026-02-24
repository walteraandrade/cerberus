package style

import "github.com/charmbracelet/lipgloss"

var (
	Red  = lipgloss.Color("#E53935")
	Gold = lipgloss.Color("#FDD835")
	Dim  = lipgloss.Color("#666666")
	Text = lipgloss.Color("#EEEEEE")

	Title = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)

	Accent = lipgloss.NewStyle().
		Foreground(Gold)

	Subtle = lipgloss.NewStyle().
		Foreground(Dim)

	Selected = lipgloss.NewStyle().
			Foreground(Gold).
			Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5252")).
		Bold(true)
)
