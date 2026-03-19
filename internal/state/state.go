package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	LastSink string `json:"last_sink"`
}

func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "volta", "state.json"), nil
}

func Load() State {
	p, err := path()
	if err != nil {
		return State{}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return State{}
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}
	}
	return s
}

func Save(s State) {
	p, err := path()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	_ = os.WriteFile(p, data, 0644)
}
