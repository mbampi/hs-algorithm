package main

import "sync"

type NumMessages struct {
	num int
	sync.Mutex
}

func (n *NumMessages) increment() {
	n.Lock()
	defer n.Unlock()
	n.num++
}

var numMessages NumMessages
