package tui

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/walteraandrade/cerberus/internal/generator"
	"github.com/walteraandrade/cerberus/internal/style"
	"github.com/walteraandrade/cerberus/internal/vault"
)

const (
	fieldTitle = iota
	fieldURL
	fieldUsername
	fieldPassword
	fieldCategory
	fieldNotes
	fieldCount
)

type EditModel struct {
	inputs  []textinput.Model
	focus   int
	editing bool
	entryID string
	genMode int // 0=password, 1=passphrase
}

type SaveEntryMsg struct{ Entry vault.Entry }

func NewEditModel(entry *vault.Entry) EditModel {
	labels := []string{"Title", "URL", "Username", "Password", "Category", "Notes"}
	inputs := make([]textinput.Model, fieldCount)

	for i, label := range labels {
		ti := textinput.New()
		ti.Placeholder = label
		ti.CharLimit = 512
		if i == fieldPassword {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		inputs[i] = ti
	}

	m := EditModel{inputs: inputs}

	if entry != nil {
		m.editing = true
		m.entryID = entry.ID
		inputs[fieldTitle].SetValue(entry.Title)
		inputs[fieldURL].SetValue(entry.URL)
		inputs[fieldUsername].SetValue(entry.Username)
		inputs[fieldPassword].SetValue(entry.Password)
		inputs[fieldCategory].SetValue(entry.Category)
		inputs[fieldNotes].SetValue(entry.Notes)
	}

	inputs[0].Focus()
	m.inputs = inputs
	return m
}

func (m EditModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, Keys.Generate):
			return m.handleGenerate()
		case msg.String() == "ctrl+r":
			m.genMode = (m.genMode + 1) % 2
			return m, nil
		case key.Matches(msg, Keys.Tab):
			return m.nextField(), nil
		case key.Matches(msg, Keys.ShiftTab):
			return m.prevField(), nil
		case msg.String() == "enter":
			if m.focus == fieldCount-1 {
				return m, m.save()
			}
			return m.nextField(), nil
		case msg.String() == "ctrl+s":
			return m, m.save()
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m EditModel) handleGenerate() (EditModel, tea.Cmd) {
	var pw string
	var err error
	if m.genMode == 0 {
		pw, err = generator.Password(generator.DefaultPasswordOpts())
	} else {
		pw, err = generator.Passphrase(generator.DefaultPassphraseOpts())
	}
	if err != nil {
		return m, nil
	}
	m.inputs[fieldPassword].SetValue(pw)
	return m, nil
}

func (m EditModel) save() tea.Cmd {
	title := strings.TrimSpace(m.inputs[fieldTitle].Value())
	if title == "" {
		return nil
	}

	id := m.entryID
	if id == "" {
		b := make([]byte, 8)
		rand.Read(b)
		id = hex.EncodeToString(b)
	}

	now := time.Now()
	entry := vault.Entry{
		ID:        id,
		Title:     title,
		URL:       strings.TrimSpace(m.inputs[fieldURL].Value()),
		Username:  strings.TrimSpace(m.inputs[fieldUsername].Value()),
		Password:  m.inputs[fieldPassword].Value(),
		Category:  strings.TrimSpace(m.inputs[fieldCategory].Value()),
		Notes:     strings.TrimSpace(m.inputs[fieldNotes].Value()),
		UpdatedAt: now,
	}
	if !m.editing {
		entry.CreatedAt = now
	}
	return func() tea.Msg { return SaveEntryMsg{entry} }
}

func (m EditModel) nextField() EditModel {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % fieldCount
	m.inputs[m.focus].Focus()
	return m
}

func (m EditModel) prevField() EditModel {
	m.inputs[m.focus].Blur()
	m.focus = (m.focus - 1 + fieldCount) % fieldCount
	m.inputs[m.focus].Focus()
	return m
}

func (m EditModel) View() string {
	var b strings.Builder

	action := "Add Entry"
	if m.editing {
		action = "Edit Entry"
	}
	header := lipgloss.NewStyle().
		Foreground(style.Red).
		Bold(true).
		Render(action)

	b.WriteString(header + "\n\n")

	labels := []string{"Title", "URL", "Username", "Password", "Category", "Notes"}
	labelStyle := lipgloss.NewStyle().Foreground(style.Gold).Width(12)

	for i, input := range m.inputs {
		prefix := "  "
		if i == m.focus {
			prefix = style.Accent.Render("> ")
		}
		b.WriteString(fmt.Sprintf("%s%s %s\n", prefix, labelStyle.Render(labels[i]), input.View()))
	}

	b.WriteString("\n")
	genLabel := "password"
	if m.genMode == 1 {
		genLabel = "passphrase"
	}
	help := fmt.Sprintf("tab:next  ctrl+s:save  ctrl+g:generate (%s)  ctrl+r:toggle mode  esc:cancel", genLabel)
	b.WriteString(lipgloss.NewStyle().Foreground(style.Dim).Render(help))

	return b.String()
}
