package audio

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	sinkRe = regexp.MustCompile(`^Sink #(\d+)`)
	nameRe = regexp.MustCompile(`^\s+Name:\s+(.+)`)
	descRe = regexp.MustCompile(`^\s+Description:\s+(.+)`)
	volRe  = regexp.MustCompile(`front-left:\s+\d+\s+/\s+(\d+)%.*front-right:\s+\d+\s+/\s+(\d+)%`)
	muteRe = regexp.MustCompile(`^\s+Mute:\s+(yes|no)`)
)

func ListSinks() ([]Sink, error) {
	out, err := exec.Command("pactl", "list", "sinks").Output()
	if err != nil {
		return nil, fmt.Errorf("pactl list sinks: %w", err)
	}
	return parseSinks(string(out))
}

func SetVolume(sinkName string, left, right int) error {
	vol := fmt.Sprintf("%d%%/%d%%", left, right)
	if err := exec.Command("pactl", "set-sink-volume", sinkName, vol).Run(); err != nil {
		return fmt.Errorf("pactl set-sink-volume: %w", err)
	}
	return nil
}

func SetMute(sinkName string, mute bool) error {
	val := "0"
	if mute {
		val = "1"
	}
	if err := exec.Command("pactl", "set-sink-mute", sinkName, val).Run(); err != nil {
		return fmt.Errorf("pactl set-sink-mute: %w", err)
	}
	return nil
}

func parseSinks(output string) ([]Sink, error) {
	var sinks []Sink
	var cur *Sink

	for _, line := range strings.Split(output, "\n") {
		if m := sinkRe.FindStringSubmatch(line); m != nil {
			if cur != nil {
				sinks = append(sinks, *cur)
			}
			idx, _ := strconv.Atoi(m[1])
			cur = &Sink{Index: idx}
			continue
		}
		if cur == nil {
			continue
		}
		if m := nameRe.FindStringSubmatch(line); m != nil {
			cur.Name = strings.TrimSpace(m[1])
		} else if m := descRe.FindStringSubmatch(line); m != nil {
			cur.Description = strings.TrimSpace(m[1])
		} else if m := volRe.FindStringSubmatch(line); m != nil {
			cur.Volume.Left, _ = strconv.Atoi(m[1])
			cur.Volume.Right, _ = strconv.Atoi(m[2])
		} else if m := muteRe.FindStringSubmatch(line); m != nil {
			cur.Muted = m[1] == "yes"
		}
	}
	if cur != nil {
		sinks = append(sinks, *cur)
	}
	return sinks, nil
}
