package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
)

type ConfirmModel struct {
	message string
	onYes   tea.Msg
}

func NewConfirmModel(message string, onYes tea.Msg) ConfirmModel {
	return ConfirmModel{message: message, onYes: onYes}
}

type ConfirmYesMsg struct{ Inner tea.Msg }
type ConfirmNoMsg struct{}

func (m ConfirmModel) Init() tea.Cmd { return nil }

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "Y", "enter":
			inner := m.onYes
			return m, func() tea.Msg { return ConfirmYesMsg{inner} }
		case "n", "N", "esc":
			return m, func() tea.Msg { return ConfirmNoMsg{} }
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	msg := lipgloss.NewStyle().Foreground(style.Text).Render(m.message)
	prompt := lipgloss.NewStyle().Foreground(style.Gold).Render(" (y/n)")
	return fmt.Sprintf("\n%s%s", msg, prompt)
}
