package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	RightUp   key.Binding
	RightDown key.Binding
	Left      key.Binding
	Right     key.Binding
	SinkNext  key.Binding
	SinkPrev  key.Binding
	Lock      key.Binding
	Reset     key.Binding
	Mute      key.Binding
	Quit      key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "L vol +5"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "L vol -5"),
	),
	RightUp: key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "R vol +5"),
	),
	RightDown: key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "R vol -5"),
	),
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "balance -3"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "balance +3"),
	),
	SinkNext: key.NewBinding(
		key.WithKeys("tab", "right"),
		key.WithHelp("→/tab", "next sink"),
	),
	SinkPrev: key.NewBinding(
		key.WithKeys("shift+tab", "left"),
		key.WithHelp("←/shift+tab", "prev sink"),
	),
	Lock: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "toggle lock L=R"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset 100%"),
	),
	Mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle mute"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
