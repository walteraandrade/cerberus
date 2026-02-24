package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Top      key.Binding
	Bottom   key.Binding
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Add      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Copy     key.Binding
	Search   key.Binding
	Help     key.Binding
	Generate key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Reveal   key.Binding
	Confirm  key.Binding
	Cancel   key.Binding
}

var Keys = KeyMap{
	Up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/up", "up")),
	Down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/down", "down")),
	Top:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "top")),
	Bottom:   key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "bottom")),
	Enter:    key.NewBinding(key.WithKeys("enter", "l"), key.WithHelp("enter/l", "open")),
	Back:     key.NewBinding(key.WithKeys("esc", "h"), key.WithHelp("esc/h", "back")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Add:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Edit:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Copy:     key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy password")),
	Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Generate: key.NewBinding(key.WithKeys("ctrl+g"), key.WithHelp("ctrl+g", "generate")),
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev field")),
	Reveal:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reveal")),
	Confirm:  key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
	Cancel:   key.NewBinding(key.WithKeys("n", "esc"), key.WithHelp("n/esc", "cancel")),
}
