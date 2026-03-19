package audio

type Volume struct {
	Left  int
	Right int
}

type Sink struct {
	Index       int
	Name        string
	Description string
	Volume      Volume
	Muted       bool
}

func (s Sink) Balance() int {
	return s.Volume.Right - s.Volume.Left
}
