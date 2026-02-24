package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
)

type HelpModel struct {
	screen Screen
}

func NewHelpModel(screen Screen) HelpModel {
	return HelpModel{screen: screen}
}

func (m HelpModel) Init() tea.Cmd { return nil }

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "?" || msg.String() == "esc" || msg.String() == "q" {
			return m, func() tea.Msg { return BackMsg{} }
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	var b strings.Builder

	header := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		Render("Keybindings")

	b.WriteString(header + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(style.Gold).Width(14)
	descStyle := lipgloss.NewStyle().Foreground(style.Text)

	var bindings [][]string

	common := [][]string{
		{"q / ctrl+c", "Quit"},
		{"?", "Toggle help"},
	}

	switch m.screen {
	case ScreenList:
		bindings = [][]string{
			{"j / down", "Move down"},
			{"k / up", "Move up"},
			{"g", "Go to top"},
			{"G", "Go to bottom"},
			{"enter / l", "Open entry"},
			{"a", "Add entry"},
			{"d", "Delete entry"},
			{"y", "Copy password"},
			{"/", "Search"},
			{"e", "Export vault"},
			{"p", "Change password"},
		}
	case ScreenDetail:
		bindings = [][]string{
			{"esc / h", "Back to list"},
			{"e", "Edit entry"},
			{"y", "Copy password"},
			{"r", "Reveal/hide password"},
			{"d", "Delete entry"},
		}
	case ScreenEdit:
		bindings = [][]string{
			{"tab", "Next field"},
			{"shift+tab", "Previous field"},
			{"ctrl+s", "Save"},
			{"ctrl+g", "Generate password"},
			{"ctrl+r", "Toggle gen mode"},
			{"esc", "Cancel"},
		}
	}

	bindings = append(bindings, common...)

	for _, pair := range bindings {
		b.WriteString("  " + keyStyle.Render(pair[0]) + descStyle.Render(pair[1]) + "\n")
	}

	b.WriteString("\n" + lipgloss.NewStyle().Foreground(style.Dim).Render("Press ? or esc to close"))

	return b.String()
}
