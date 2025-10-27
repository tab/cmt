package commit

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the key bindings for the commit TUI
type KeyMap struct {
	Accept      key.Binding
	Edit        key.Binding
	Regenerate  key.Binding
	ToggleLogs  key.Binding
	ToggleFocus key.Binding
	Quit        key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Accept: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "accept"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Regenerate: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "regenerate"),
		),
		ToggleLogs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "toggle logs"),
		),
		ToggleFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch pane"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Accept, k.Edit, k.Regenerate, k.ToggleFocus, k.ToggleLogs, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Accept, k.Edit, k.Regenerate, k.ToggleFocus, k.ToggleLogs, k.Quit},
	}
}
