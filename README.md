# volta

A terminal UI for per-channel (L/R) volume and balance control using PulseAudio / PipeWire (`pactl`).

```
╭──────────────────────────────────────────────────────────────────╮
│  volta                                                           │
│    ◀ EDIFIER W830NB ▶  (3 / 6)                                   │
│                                                                  │
│    L   ██████████████████████████████████████░░░░░░░░░░ 80  %    │
│    R   ████████████████████████████████████████████████ 100 %    │
│                                                                  │
│    BAL ────────────────────────┼──●────────────────────  +20 R   │
│                                                                  │
│    ──────────────────────────────────────────────────────────    │
│    k / j  L vol  +5 / -5       │  K / J  R vol  +5 / -5          │
│    h / l  balance  ±3          │  ← / →  prev / next sink        │
│    L      lock  L = R          │  r      reset to 100%           │
│    m      mute toggle          │  q      quit                    │
│                                                                  │
╰──────────────────────────────────────────────────────────────────╯
```

## Features

- **Per-channel control** — adjust left and right volume independently
- **Balance bar** — visual representation of the stereo balance
- **Lock mode** — link L and R channels so they move together
- **Sink cycling** — navigate all available audio sinks with arrow keys
- **Mute toggle** — mute/unmute any sink instantly
- **Persistent state** — remembers the last used sink across sessions
- **Responsive layout** — all bars and controls adapt to terminal width
- **PipeWire / PulseAudio** — works with any `pactl`-compatible backend

## Requirements

- Go 1.22+
- PulseAudio or PipeWire with `pactl` available

## Installation

**Pre-built binary (latest):**

```bash
# amd64
curl -L https://github.com/mateuspim/volta/releases/download/latest/volta-linux-amd64 -o volta
chmod +x volta && sudo mv volta /usr/local/bin/

# arm64
curl -L https://github.com/mateuspim/volta/releases/download/latest/volta-linux-arm64 -o volta
chmod +x volta && sudo mv volta /usr/local/bin/
```

**Specific version:**

```bash
curl -L https://github.com/mateuspim/volta/releases/download/v1.0.0/volta-linux-amd64 -o volta
chmod +x volta && sudo mv volta /usr/local/bin/
```

**From source:**

```bash
git clone git@github.com:mateuspim/volta.git
cd volta
make build
# binary at ./bin/volta
```

## Controls

| Key | Action |
|---|---|
| `k` / `j` | Left channel volume +5 / -5 |
| `K` / `J` | Right channel volume +5 / -5 |
| `h` / `l` | Shift balance left / right by 3 |
| `←` / `→` | Previous / next sink |
| `tab` / `shift+tab` | Previous / next sink |
| `L` | Toggle lock L = R |
| `r` | Reset both channels to 100% |
| `m` | Toggle mute |
| `q` / `ctrl+c` | Quit |

> Arrow keys cycle sinks. `hjkl` / `KJ` control volume and balance exclusively.

## Volume range

Bars support 0–150%. The 100–150% range renders in red as an over-gain indicator.

## State

Last used sink is saved to `~/.config/volta/state.json` and restored on the next launch. If the sink is unavailable (e.g. Bluetooth headphones unplugged), volta falls back to the first available sink.

## Project structure

```
volta/
├── cmd/volta/          # entrypoint
├── internal/
│   ├── audio/          # pactl wrapper (list sinks, set volume, mute)
│   ├── tui/            # bubbletea model, update loop, rendering
│   └── state/          # persistent sink state (~/.config/volta/)
├── .github/workflows/  # CI: build + test on dev, release binaries on tags
└── Makefile
```

## CI / Releases

- Push to `dev` or PR to `master` → build + vet + test
- Tag `v*` on `master` → release with `linux/amd64` and `linux/arm64` binaries

## Tech stack

| | |
|---|---|
| TUI framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| Styling | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| Components | [Bubbles](https://github.com/charmbracelet/bubbles) |
| Audio backend | `pactl` (PulseAudio / PipeWire) |

## License

MIT
