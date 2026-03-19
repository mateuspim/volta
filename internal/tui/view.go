package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const maxVol = 150

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

	// Fix box to terminal width minus 2 so right border is always visible.
	// inner = boxWidth - border(2) - padding(4)
	// bar   = inner - indent(2) - label(3) - spaces(2) - vol(4) - %(1) = inner - 12
	boxWidth := m.width - 2
	if boxWidth < 40 {
		boxWidth = 40
	}
	inner := boxWidth - 6
	bw := inner - 12
	if bw < 10 {
		bw = 10
	}

	leftBar := renderBar(sink.Volume.Left, bw)
	rightBar := renderBar(sink.Volume.Right, bw)

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

	// Balance — two lines: track + label below right-aligned under vol%
	balBlock := renderBalance(sink.Balance(), bw)

	// Help grid — 2 columns: [key] [desc]
	helpGrid := renderHelp(bw)

	lines := []string{
		"",
		"  " + sinkNav + "  " + sinkIdx + badges,
		"",
		"  " + leftLine,
		"  " + rightLine,
		"",
		balBlock,
		helpGrid,
		"",
	}

	header := titleStyle.Render("volta")
	body := strings.Join(lines, "\n")

	return boxStyle.Width(boxWidth).Render(header + body)
}

func renderBar(vol, bw int) string {
	normalFill := vol * bw / 100
	if normalFill > bw {
		normalFill = bw
	}
	if normalFill < 0 {
		normalFill = 0
	}

	overFill := 0
	if vol > 100 {
		overFill = (vol - 100) * bw / 50
		if overFill > bw-normalFill {
			overFill = bw - normalFill
		}
	}

	empty := bw - normalFill - overFill

	return filledStyle.Render(strings.Repeat("█", normalFill)) +
		overStyle.Render(strings.Repeat("█", overFill)) +
		emptyStyle.Render(strings.Repeat("░", empty))
}

func renderBalance(balance, bw int) string {
	total := bw
	center := total / 2

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

	// label right-aligned in Width(6) appended directly after track:
	// "BAL "(4) + track(bw) + label(6) = bw+10 — matches slider line width
	var bareLabel string
	var labelColor lipgloss.Color
	switch {
	case balance == 0:
		bareLabel = "L = R"
		labelColor = "#10B981"
	case balance > 0:
		bareLabel = fmt.Sprintf("+%d R", balance)
		labelColor = "#A78BFA"
	default:
		bareLabel = fmt.Sprintf("%d L", -balance)
		labelColor = "#A78BFA"
	}
	styledLabel := lipgloss.NewStyle().Foreground(labelColor).Width(6).Align(lipgloss.Right).Render(bareLabel)

	return "  " + fmt.Sprintf("%s %s%s",
		labelStyle.Render("BAL"),
		strings.Join(bar, ""),
		styledLabel,
	)
}

func renderHelp(bw int) string {
	const keyW = 7
	// match slider content width (bw+10) after "  " indent:
	// keyW + descW + "  │  "(5) + keyW + descW = bw+10 → descW = (bw-9)/2
	descW := (bw - 9) / 2
	if descW < 12 {
		descW = 12
	}
	sepW := bw + 10

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C4B5FD")).
		Bold(true).
		Width(keyW)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Width(descW)
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

	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563")).Render(strings.Repeat("─", sepW))
	lines := []string{"  " + sep}
	for _, row := range rows {
		left := keyStyle.Render(row[0].k) + descStyle.Render(row[0].d)
		right := keyStyle.Render(row[1].k) + descStyle.Render(row[1].d)
		lines = append(lines, "  "+left+"  "+divider+"  "+right)
	}
	return strings.Join(lines, "\n")
}
