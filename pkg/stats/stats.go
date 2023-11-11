package stats

import "sync"

// NumMessages is the number of messages sent between processes.
var NumMessages MessagesCount

type MessagesCount struct {
	num int
	sync.Mutex
}

// Increment increments the number of messages.
// It is safe to call Increment concurrently.
func (n *MessagesCount) Increment() {
	n.Lock()
	defer n.Unlock()
	n.num++
}

// Get returns the number of messages.
func (n *MessagesCount) Get() int {
	n.Lock()
	defer n.Unlock()
	return n.num
}
