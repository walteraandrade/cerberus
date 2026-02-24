package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
	"github.com/walteraandrade/cerberus/internal/vault"
)

type DetailModel struct {
	entry    vault.Entry
	revealed bool
}

func NewDetailModel(entry vault.Entry) DetailModel {
	return DetailModel{entry: entry}
}

type EditEntryMsg struct{ Entry vault.Entry }
type BackMsg struct{}

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			return m, func() tea.Msg { return BackMsg{} }
		case key.Matches(msg, Keys.Edit):
			return m, func() tea.Msg { return EditEntryMsg{m.entry} }
		case key.Matches(msg, Keys.Copy):
			pw := m.entry.Password
			return m, func() tea.Msg { return CopyPasswordMsg{pw} }
		case key.Matches(msg, Keys.Reveal):
			m.revealed = !m.revealed
		case key.Matches(msg, Keys.Delete):
			return m, func() tea.Msg { return DeleteEntryMsg{m.entry} }
		case key.Matches(msg, Keys.Help):
			return m, func() tea.Msg { return helpMsg{} }
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DetailModel) View() string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		Render(m.entry.Title)

	b.WriteString(title + "\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(style.Gold).Width(12)
	valStyle := lipgloss.NewStyle().Foreground(style.Text)

	fields := []struct{ label, value string }{
		{"URL", m.entry.URL},
		{"Username", m.entry.Username},
		{"Password", m.passwordDisplay()},
		{"Category", m.entry.Category},
		{"Notes", m.entry.Notes},
		{"Created", m.entry.CreatedAt.Format("2006-01-02 15:04")},
		{"Modified", m.entry.UpdatedAt.Format("2006-01-02 15:04")},
	}

	for _, f := range fields {
		if f.value == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("%s %s\n",
			labelStyle.Render(f.label),
			valStyle.Render(f.value)))
	}

	b.WriteString("\n")
	help := "esc/h:back  e:edit  y:copy  r:reveal  d:del  q:quit"
	b.WriteString(lipgloss.NewStyle().Foreground(style.Dim).Render(help))

	return b.String()
}

func (m DetailModel) passwordDisplay() string {
	if m.revealed {
		return m.entry.Password
	}
	return strings.Repeat("•", min(len(m.entry.Password), 20))
}
