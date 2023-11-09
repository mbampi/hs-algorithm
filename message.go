package main

// Message represents a message sent between processes, containing the UID of the
// process that sent the message, the number of hops the message has to take
// before completing the round trip, and the direction of the message.
// Called "token" in the book.
type Message struct {
	uid  int // Unique ID of the node that sent the message.
	hops int // Number of hops the message has to take.
	way  Way // 0 for left, 1 for right.
}

// Way represents the direction of a message.
type Way int

const (
	// In represents a message going back to the process that sent it.
	In Way = iota
	// Out represents a message going away from the process that sent it.
	Out
)
