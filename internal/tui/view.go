package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	barWidth = 24
	maxVol   = 150
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

	// Help grid — 2 columns: [key] [desc]
	helpGrid := renderHelp()

	lines := []string{
		"",
		"  " + sinkNav + "  " + sinkIdx + badges,
		"",
		"  " + leftLine,
		"  " + rightLine,
		"",
		"  " + balLine,
		"",
		helpGrid,
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
	// balance: -150..+150, track matches barWidth exactly
	total := barWidth
	center := total / 2

	// Map balance to position within half-width
	half := center
	pos := balance * half / 150
	if pos > half {
		pos = half
	}
	if pos < -half {
		pos = -half
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

func renderHelp() string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C4B5FD")).
		Bold(true).
		Width(7)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Width(16)
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Render("│")

	type entry struct{ k, d string }
	rows := [][]entry{
		{{"k / j", "L vol  +5 / -5"}, {"K / J", "R vol  +5 / -5"}},
		{{"h / l", "balance  ±3"}, {"← / →", "prev / next sink"}},
		{{"L", "lock  L = R"}, {"r", "reset to 100%"}},
		{{"m", "mute toggle"}, {"q", "quit"}},
	}

	// row width: keyW(7) + descW(16) + "  │  "(5) + keyW(7) + descW(16) = 51
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563")).Render(strings.Repeat("─", 51))
	lines := []string{"  " + sep}
	for _, row := range rows {
		left := keyStyle.Render(row[0].k) + descStyle.Render(row[0].d)
		right := keyStyle.Render(row[1].k) + descStyle.Render(row[1].d)
		lines = append(lines, "  "+left+"  "+divider+"  "+right)
	}
	return strings.Join(lines, "\n")
}
