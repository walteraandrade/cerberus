package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
)

type PasswordChangeModel struct {
	inputs []textinput.Model
	focus  int // 0=current, 1=new, 2=confirm
	err    string
}

type PasswordChangedMsg struct {
	OldPassword string
	NewPassword string
}

func NewPasswordChangeModel() PasswordChangeModel {
	labels := []string{"Current password", "New password", "Confirm new password"}
	inputs := make([]textinput.Model, 3)
	for i, l := range labels {
		ti := textinput.New()
		ti.Placeholder = l
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
		ti.CharLimit = 256
		inputs[i] = ti
	}
	inputs[0].Focus()
	return PasswordChangeModel{inputs: inputs}
}

func (m PasswordChangeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PasswordChangeModel) Update(msg tea.Msg) (PasswordChangeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			return m.nextField(), nil
		case "shift+tab":
			return m.prevField(), nil
		case "enter":
			if m.focus < 2 {
				return m.nextField(), nil
			}
			return m.submit()
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m PasswordChangeModel) submit() (PasswordChangeModel, tea.Cmd) {
	current := m.inputs[0].Value()
	newPw := m.inputs[1].Value()
	confirm := m.inputs[2].Value()

	if current == "" || newPw == "" {
		m.err = "all fields required"
		return m, nil
	}
	if newPw != confirm {
		m.err = "new passwords don't match"
		m.inputs[2].SetValue("")
		return m, nil
	}
	if current == newPw {
		m.err = "new password must differ"
		return m, nil
	}

	m.err = ""
	old, new_ := current, newPw
	return m, func() tea.Msg {
		return PasswordChangedMsg{OldPassword: old, NewPassword: new_}
	}
}

func (m PasswordChangeModel) nextField() PasswordChangeModel {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % 3
	m.inputs[m.focus].Focus()
	return m
}

func (m PasswordChangeModel) prevField() PasswordChangeModel {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus - 1 + 3) % 3
	m.inputs[m.focus].Focus()
	return m
}

func (m PasswordChangeModel) View() string {
	header := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		Render("Change Master Password")

	s := header + "\n\n"

	labels := []string{"Current", "New", "Confirm"}
	labelStyle := lipgloss.NewStyle().Foreground(style.Gold).Width(12)

	for i, input := range m.inputs {
		prefix := "  "
		if i == m.focus {
			prefix = style.Accent.Render("> ")
		}
		s += prefix + labelStyle.Render(labels[i]) + " " + input.View() + "\n"
	}

	if m.err != "" {
		s += "\n" + style.Error.Render(m.err)
	}

	s += "\n\n" + lipgloss.NewStyle().Foreground(style.Dim).Render("tab:next  enter:submit  esc:cancel")
	return s
}
