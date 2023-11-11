package hs

// Way represents the direction of a message.
type Way int

const (
	// In represents a message going back to the process that sent it.
	In Way = iota
	// Out represents a message going away from the process that sent it.
	Out
)

// String returns a string representation of a Way.
func (w Way) String() string {
	if w == In {
		return "In"
	}
	return "Out"
}
