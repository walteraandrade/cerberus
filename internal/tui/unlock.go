package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/style"
)

type UnlockModel struct {
	input    textinput.Model
	err      string
	creating bool
	confirm  textinput.Model
	stage    int // 0=password, 1=confirm (create only)
}

func NewUnlockModel(vaultExists bool) UnlockModel {
	ti := textinput.New()
	ti.Placeholder = "Master password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.Focus()
	ti.CharLimit = 256

	ci := textinput.New()
	ci.Placeholder = "Confirm password"
	ci.EchoMode = textinput.EchoPassword
	ci.EchoCharacter = '•'
	ci.CharLimit = 256

	return UnlockModel{
		input:    ti,
		creating: !vaultExists,
		confirm:  ci,
	}
}

type UnlockMsg struct {
	Password string
	Create   bool
}

type UnlockErrMsg struct {
	Err string
}

func (m UnlockModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UnlockModel) Update(msg tea.Msg) (UnlockModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.handleEnter()
		case "ctrl+c":
			return m, tea.Quit
		}
	case UnlockErrMsg:
		m.err = msg.Err
		m.input.SetValue("")
		m.input.Focus()
		m.stage = 0
		return m, nil
	}

	var cmd tea.Cmd
	if m.stage == 0 {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.confirm, cmd = m.confirm.Update(msg)
	}
	return m, cmd
}

func (m UnlockModel) handleEnter() (UnlockModel, tea.Cmd) {
	if m.creating {
		if m.stage == 0 {
			if m.input.Value() == "" {
				m.err = "password cannot be empty"
				return m, nil
			}
			m.stage = 1
			m.err = ""
			m.input.Blur()
			m.confirm.Focus()
			return m, textinput.Blink
		}
		if m.input.Value() != m.confirm.Value() {
			m.err = "passwords don't match"
			m.confirm.SetValue("")
			return m, nil
		}
	} else {
		if m.input.Value() == "" {
			m.err = "password cannot be empty"
			return m, nil
		}
	}
	m.err = ""
	return m, func() tea.Msg {
		return UnlockMsg{Password: m.input.Value(), Create: m.creating}
	}
}

func (m UnlockModel) View() string {
	var s string

	title := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		MarginBottom(1).
		Render("CERBERUS")

	subtitle := lipgloss.NewStyle().
		Foreground(style.Gold).
		Render("Password Vault")

	s += title + "\n" + subtitle + "\n\n"

	if m.creating {
		s += lipgloss.NewStyle().Foreground(style.Dim).Render("Creating new vault") + "\n\n"
		s += m.input.View() + "\n"
		if m.stage == 1 {
			s += m.confirm.View() + "\n"
		}
	} else {
		s += m.input.View() + "\n"
	}

	if m.err != "" {
		s += "\n" + style.Error.Render(m.err)
	}

	s += "\n\n" + lipgloss.NewStyle().Foreground(style.Dim).Render("enter: submit • ctrl+c: quit")

	return s
}
