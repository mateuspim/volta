package tui

import (
	"github.com/pym/volta/internal/audio"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type sinksLoadedMsg []audio.Sink
type volumeAppliedMsg struct{}

type Model struct {
	sinks   []audio.Sink
	sinkIdx int
	locked  bool
	err     error
	width   int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return loadSinks
}

func loadSinks() tea.Msg {
	sinks, err := audio.ListSinks()
	if err != nil {
		return errMsg(err)
	}
	return sinksLoadedMsg(sinks)
}

func cmdSetVolume(name string, left, right int) tea.Cmd {
	return func() tea.Msg {
		_ = audio.SetVolume(name, left, right)
		return volumeAppliedMsg{}
	}
}

func cmdSetMute(name string, mute bool) tea.Cmd {
	return func() tea.Msg {
		_ = audio.SetMute(name, mute)
		return volumeAppliedMsg{}
	}
}
