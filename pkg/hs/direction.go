package hs

// Direction represents the direction of a received message.
type Direction int

const (
	// Left represents a message coming from the left neighbor.
	Left Direction = iota
	// Right represents a message coming from the right neighbor.
	Right
)

// String returns a string representation of a direction.
func (d Direction) String() string {
	if d == Left {
		return "Left"
	}
	return "Right"
}

// oppositeDirection returns the opposite direction.
func oppositeDirection(dir Direction) Direction {
	if dir == Left {
		return Right
	}
	return Left
}
