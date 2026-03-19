package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	barWidth = 24
	maxVol   = 150
	balWidth = 11 // half-width of balance bar
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A78BFA"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6D28D9")).
			Padding(0, 2)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Width(3).
			Foreground(lipgloss.Color("#C4B5FD"))

	volNumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Width(4)

	filledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED"))

	overStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151"))

	sinkNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E5E7EB"))

	lockStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F59E0B"))

	muteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EF4444"))

	balCenterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).Bold(true)

	balMarkerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).Bold(true)

	balTrackStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4B5563"))
)

func (m Model) View() string {
	if m.err != nil {
		return errStyle.Render(fmt.Sprintf("error: %v\n\npress q to quit.", m.err))
	}

	if len(m.sinks) == 0 {
		return boxStyle.Render(titleStyle.Render("volta") + "\n\n  loading sinks…\n")
	}

	sink := m.sinks[m.sinkIdx]

	// Truncate long descriptions
	desc := sink.Description
	if len(desc) > 36 {
		desc = desc[:33] + "…"
	}

	// Sink nav line
	sinkNav := fmt.Sprintf("◀ %s ▶", sinkNameStyle.Render(desc))
	sinkIdx := helpStyle.Render(fmt.Sprintf("(%d / %d)", m.sinkIdx+1, len(m.sinks)))

	// Status badges
	badges := ""
	if sink.Muted {
		badges += "  " + muteStyle.Render("◉ MUTED")
	}
	if m.locked {
		badges += "  " + lockStyle.Render("⚿ LOCKED")
	}

	// Volume bars
	leftBar := renderBar(sink.Volume.Left)
	rightBar := renderBar(sink.Volume.Right)

	leftLine := fmt.Sprintf("%s %s %s%%",
		labelStyle.Render("L"),
		leftBar,
		volNumStyle.Render(fmt.Sprintf("%d", sink.Volume.Left)),
	)
	rightLine := fmt.Sprintf("%s %s %s%%",
		labelStyle.Render("R"),
		rightBar,
		volNumStyle.Render(fmt.Sprintf("%d", sink.Volume.Right)),
	)

	// Balance
	balLine := renderBalance(sink.Balance())

	// Help text
	help1 := helpStyle.Render("↑↓ L vol   K/J R vol   ←→ balance   tab next")
	help2 := helpStyle.Render("L lock   r reset   m mute   q quit")

	lines := []string{
		"",
		"  " + sinkNav + "  " + sinkIdx + badges,
		"",
		"  " + leftLine,
		"  " + rightLine,
		"",
		"  " + balLine,
		"",
		"  " + help1,
		"  " + help2,
		"",
	}

	header := titleStyle.Render("volta")
	body := strings.Join(lines, "\n")

	return boxStyle.Render(header + body)
}

func renderBar(vol int) string {
	// 0-100% fills first barWidth blocks, 101-150% overflows in red
	normalFill := vol * barWidth / 100
	if normalFill > barWidth {
		normalFill = barWidth
	}
	if normalFill < 0 {
		normalFill = 0
	}

	overFill := 0
	if vol > 100 {
		overFill = (vol - 100) * barWidth / 50
		if overFill > barWidth-normalFill {
			overFill = barWidth - normalFill
		}
	}

	empty := barWidth - normalFill - overFill

	return filledStyle.Render(strings.Repeat("█", normalFill)) +
		overStyle.Render(strings.Repeat("█", overFill)) +
		emptyStyle.Render(strings.Repeat("░", empty))
}

func renderBalance(balance int) string {
	// balance: -150..+150, display on 2*balWidth+1 track
	total := 2*balWidth + 1
	center := balWidth

	// Map balance to position
	pos := balance * balWidth / 150
	if pos > balWidth {
		pos = balWidth
	}
	if pos < -balWidth {
		pos = -balWidth
	}
	markerPos := center + pos

	bar := make([]string, total)
	for i := 0; i < total; i++ {
		switch {
		case i == markerPos && i == center:
			bar[i] = balCenterStyle.Render("●")
		case i == markerPos:
			bar[i] = balMarkerStyle.Render("●")
		case i == center:
			bar[i] = balTrackStyle.Render("┼")
		default:
			bar[i] = balTrackStyle.Render("─")
		}
	}

	var balLabel string
	switch {
	case balance == 0:
		balLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render("L = R")
	case balance > 0:
		balLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA")).Render(fmt.Sprintf("+%d R", balance))
	default:
		balLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA")).Render(fmt.Sprintf("%d L", -balance))
	}

	return fmt.Sprintf("%s %s  %s",
		labelStyle.Render("BAL"),
		strings.Join(bar, ""),
		balLabel,
	)
}
