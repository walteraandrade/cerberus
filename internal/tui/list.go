package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
	"github.com/walteraandrade/cerberus/internal/vault"
)

type ListModel struct {
	entries   []vault.Entry
	filtered  []vault.Entry
	cursor    int
	search    textinput.Model
	searching bool
	width     int
	height    int
}

func NewListModel(entries []vault.Entry) ListModel {
	si := textinput.New()
	si.Placeholder = "Search..."
	si.CharLimit = 128

	m := ListModel{
		entries: entries,
		search:  si,
	}
	m.applyFilter()
	return m
}

func (m *ListModel) SetEntries(entries []vault.Entry) {
	m.entries = entries
	m.applyFilter()
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

type SelectEntryMsg struct{ Entry vault.Entry }
type AddEntryMsg struct{}
type DeleteEntryMsg struct{ Entry vault.Entry }
type CopyPasswordMsg struct{ Password string }

func (m ListModel) Init() tea.Cmd { return nil }

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.searching {
			return m.updateSearch(msg)
		}
		return m.updateNav(msg)
	}

	return m, nil
}

func (m ListModel) updateNav(msg tea.KeyMsg) (ListModel, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Down):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case key.Matches(msg, Keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, Keys.Top):
		m.cursor = 0
	case key.Matches(msg, Keys.Bottom):
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		}
	case key.Matches(msg, Keys.Enter):
		if len(m.filtered) > 0 {
			return m, func() tea.Msg { return SelectEntryMsg{m.filtered[m.cursor]} }
		}
	case key.Matches(msg, Keys.Add):
		return m, func() tea.Msg { return AddEntryMsg{} }
	case key.Matches(msg, Keys.Delete):
		if len(m.filtered) > 0 {
			return m, func() tea.Msg { return DeleteEntryMsg{m.filtered[m.cursor]} }
		}
	case key.Matches(msg, Keys.Copy):
		if len(m.filtered) > 0 {
			pw := m.filtered[m.cursor].Password
			return m, func() tea.Msg { return CopyPasswordMsg{pw} }
		}
	case key.Matches(msg, Keys.Search):
		m.searching = true
		m.search.Focus()
		return m, textinput.Blink
	case key.Matches(msg, Keys.Help):
		return m, func() tea.Msg { return helpMsg{} }
	case msg.String() == "x":
		return m, func() tea.Msg { return exportMsg{} }
	case msg.String() == "p":
		return m, func() tea.Msg { return passwordChangeMsg{} }
	case key.Matches(msg, Keys.Quit):
		return m, tea.Quit
	}
	return m, nil
}

func (m ListModel) updateSearch(msg tea.KeyMsg) (ListModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.search.SetValue("")
		m.search.Blur()
		m.applyFilter()
		return m, nil
	case "enter":
		m.searching = false
		m.search.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	m.applyFilter()
	m.cursor = 0
	return m, cmd
}

func (m *ListModel) applyFilter() {
	q := strings.ToLower(m.search.Value())
	if q == "" {
		m.filtered = m.entries
		return
	}
	m.filtered = nil
	for _, e := range m.entries {
		if strings.Contains(strings.ToLower(e.Title), q) ||
			strings.Contains(strings.ToLower(e.Username), q) ||
			strings.Contains(strings.ToLower(e.URL), q) ||
			strings.Contains(strings.ToLower(e.Category), q) {
			m.filtered = append(m.filtered, e)
		}
	}
}

func (m ListModel) View() string {
	var b strings.Builder

	header := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		Render("CERBERUS")

	count := lipgloss.NewStyle().
		Foreground(style.Dim).
		Render(fmt.Sprintf(" (%d entries)", len(m.entries)))

	b.WriteString(header + count + "\n")

	if m.searching {
		b.WriteString(m.search.View() + "\n")
	}
	b.WriteString("\n")

	if len(m.filtered) == 0 {
		empty := "No entries"
		if len(m.entries) == 0 {
			empty = "Vault is empty — press 'a' to add"
		}
		b.WriteString(lipgloss.NewStyle().Foreground(style.Dim).Render(empty))
	} else {
		visible := m.visibleRange()
		for i := visible[0]; i <= visible[1]; i++ {
			e := m.filtered[i]
			cursor := "  "
			s := lipgloss.NewStyle().Foreground(style.Text)
			if i == m.cursor {
				cursor = style.Accent.Render("> ")
				s = s.Foreground(style.Gold).Bold(true)
			}

			title := s.Render(e.Title)
			meta := lipgloss.NewStyle().Foreground(style.Dim).Render(
				fmt.Sprintf(" %s", e.Username))

			if e.Category != "" {
				cat := lipgloss.NewStyle().Foreground(style.Red).Render(
					fmt.Sprintf(" [%s]", e.Category))
				meta += cat
			}

			b.WriteString(cursor + title + meta + "\n")
		}
	}

	b.WriteString("\n")
	help := "j/k:nav  enter/l:open  a:add  d:del  y:copy  /:search  x:export  p:passwd  ?:help  q:quit"
	b.WriteString(lipgloss.NewStyle().Foreground(style.Dim).Render(help))

	return b.String()
}

func (m ListModel) visibleRange() [2]int {
	maxVisible := m.height - 6
	if maxVisible < 5 {
		maxVisible = 20
	}
	total := len(m.filtered)
	if total <= maxVisible {
		return [2]int{0, total - 1}
	}

	half := maxVisible / 2
	start := m.cursor - half
	if start < 0 {
		start = 0
	}
	end := start + maxVisible - 1
	if end >= total {
		end = total - 1
		start = end - maxVisible + 1
	}
	return [2]int{start, end}
}
