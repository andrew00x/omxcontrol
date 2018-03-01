package omxcontrol

type Status int

const (
	Unknown Status = iota
	Playing
	Paused
)

var statusValues = [...]string{"Unknown", "Playing", "Paused"}

func (s Status) String() string {
	return statusValues[s]
}
