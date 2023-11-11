package hs

import "fmt"

// Message represents a message sent between processes, containing the UID of the
// process that sent the message, the number of hops the message has to take
// before completing the round trip, and the direction of the message.
// Called "token" in the book.
type Message struct {
	uid  int // Unique ID of the node that sent the message.
	hops int // Number of hops the message has to take.
	way  Way // In or Out.
}

// String returns a string representation of a message.
func (m *Message) String() string {
	return fmt.Sprintf("uid=%d; hops=%d, way=%s", m.uid, m.hops, m.way.String())
}
