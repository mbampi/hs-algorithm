package stats

import "sync"

type MessagesCount struct {
	num int
	sync.Mutex
}

func (n *MessagesCount) Increment() {
	n.Lock()
	defer n.Unlock()
	n.num++
}

func (n *MessagesCount) Get() int {
	n.Lock()
	defer n.Unlock()
	return n.num
}

var NumMessages MessagesCount
