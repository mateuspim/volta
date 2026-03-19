package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pym/volta/internal/audio"
	"github.com/pym/volta/internal/state"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case sinksLoadedMsg:
		m.sinks = []audio.Sink(msg)
		m.sinkIdx = 0
		if saved := state.Load().LastSink; saved != "" {
			for i, s := range m.sinks {
				if s.Name == saved {
					m.sinkIdx = i
					break
				}
			}
		}
	case errMsg:
		m.err = msg
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		if len(m.sinks) > 0 {
			state.Save(state.State{LastSink: m.sinks[m.sinkIdx].Name})
		}
		return m, tea.Quit
	}

	if len(m.sinks) == 0 {
		return m, nil
	}

	sink := m.sinks[m.sinkIdx]

	switch {
	case key.Matches(msg, keys.SinkNext):
		m.sinkIdx = (m.sinkIdx + 1) % len(m.sinks)
		state.Save(state.State{LastSink: m.sinks[m.sinkIdx].Name})

	case key.Matches(msg, keys.SinkPrev):
		m.sinkIdx = (m.sinkIdx - 1 + len(m.sinks)) % len(m.sinks)
		state.Save(state.State{LastSink: m.sinks[m.sinkIdx].Name})

	case key.Matches(msg, keys.Lock):
		m.locked = !m.locked
		if m.locked {
			avg := (sink.Volume.Left + sink.Volume.Right) / 2
			sink.Volume.Left = avg
			sink.Volume.Right = avg
			m.sinks[m.sinkIdx] = sink
			return m, cmdSetVolume(sink.Name, avg, avg)
		}

	case key.Matches(msg, keys.Reset):
		sink.Volume.Left = 100
		sink.Volume.Right = 100
		m.locked = false
		m.sinks[m.sinkIdx] = sink
		return m, cmdSetVolume(sink.Name, 100, 100)

	case key.Matches(msg, keys.Mute):
		sink.Muted = !sink.Muted
		m.sinks[m.sinkIdx] = sink
		return m, cmdSetMute(sink.Name, sink.Muted)

	case key.Matches(msg, keys.Up):
		l := clamp(sink.Volume.Left+5, 0, 150)
		r := sink.Volume.Right
		if m.locked {
			r = l
		}
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)

	case key.Matches(msg, keys.Down):
		l := clamp(sink.Volume.Left-5, 0, 150)
		r := sink.Volume.Right
		if m.locked {
			r = l
		}
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)

	case key.Matches(msg, keys.RightUp):
		r := clamp(sink.Volume.Right+5, 0, 150)
		l := sink.Volume.Left
		if m.locked {
			l = r
		}
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)

	case key.Matches(msg, keys.RightDown):
		r := clamp(sink.Volume.Right-5, 0, 150)
		l := sink.Volume.Left
		if m.locked {
			l = r
		}
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)

	case key.Matches(msg, keys.Left):
		l := clamp(sink.Volume.Left+3, 0, 150)
		r := clamp(sink.Volume.Right-3, 0, 150)
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)

	case key.Matches(msg, keys.Right):
		l := clamp(sink.Volume.Left-3, 0, 150)
		r := clamp(sink.Volume.Right+3, 0, 150)
		m.sinks[m.sinkIdx].Volume.Left = l
		m.sinks[m.sinkIdx].Volume.Right = r
		return m, cmdSetVolume(sink.Name, l, r)
	}

	return m, nil
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
