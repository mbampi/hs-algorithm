package hs

import (
	"fmt"
	"hsalgorithm/pkg/stats"
	"log"
	"math"
	"math/rand"
	"sync"
)

// Process represents a process in the ring. It contains the UID of the process,
// the channels to allow neighbors to communicate with it, and the
// channels to communicate with the left and right neighbors.
// Called "node" in the book.
type Process struct {
	uid            int            // Unique ID of the process.
	left           chan Message   // Channel to receive messages from the left neighbor.
	right          chan Message   // Channel to receive messages from the right neighbor.
	rightNeighbour chan<- Message // Channel to send messages to the right neighbor.
	leftNeighbour  chan<- Message // Channel to send messages to the left neighbor.

	hasLeftToken  bool // Whether the process has received the its own left token back.
	hasRightToken bool // Whether the process has received the its own right token back.
	phase         int  // Current phase of the election.
}

// CreateProcesses creates n processes and connects them in a ring,
// in a random order.
func CreateProcesses(numProcesses int) []Process {
	uids := randomUIDs(numProcesses)

	// Create processes
	processes := make([]Process, numProcesses)
	for i := 0; i < numProcesses; i++ {
		processes[i] = Process{
			uid:   uids[i],
			left:  make(chan Message, 1),
			right: make(chan Message, 1),
		}
	}

	for i := 0; i < numProcesses; i++ {
		leftIndex := (i - 1 + numProcesses) % numProcesses
		rightIndex := (i + 1) % numProcesses

		log.Printf("Process %d: Connecting to left neighbour %d and right neighbour %d\n", processes[i].uid, processes[leftIndex].uid, processes[rightIndex].uid)
		processes[i].leftNeighbour = processes[leftIndex].right
		processes[i].rightNeighbour = processes[rightIndex].left
	}

	return processes
}

// randomUIDs returns a slice of n unique random integers.
func randomUIDs(n int) []int {
	uids := make([]int, n)
	for i := 0; i < n; i++ {
		uids[i] = i
	}

	// Shuffle the uids
	for i := range uids {
		j := rand.Intn(i + 1)
		uids[i], uids[j] = uids[j], uids[i]
	}

	return uids
}

func (p *Process) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	p.startElectionPhase(0)

	for {
		select {
		// Handles the message from the left neighbor
		case msg := <-p.left:
			if p.handleMessage(msg, Left) {
				return
			}
		// Handles the message from the right neighbor
		case msg := <-p.right:
			if p.handleMessage(msg, Right) {
				return
			}
		}
	}
}

// handleMessage processes a message received from a neighbor.
func (p *Process) handleMessage(msg Message, incomingDirection Direction) bool {
	stats.NumMessages.Increment()

	// If the message UID is the same as this process and it's an Out message, declare leadership
	if msg.uid == p.uid && msg.way == Out {
		fmt.Printf("Process %d: I am the leader!\n", p.uid)
		return true
	}

	// If the message UID is the same as this process and it's an In message, process it
	if msg.uid == p.uid {
		p.handleReturningToken(incomingDirection)
		return false
	}

	// If the message is from another process and it's an Out message, process it
	if msg.way == Out {
		p.handleOutMessage(msg, incomingDirection)
	}

	// If the message is from another process and it's an In message, forward it
	if msg.way == In {
		log.Printf("Process %d: Received In message from neighbor. %s\n", p.uid, msg.String())
		p.forwardMessage(msg, incomingDirection)
	}

	return false
}

// handleReturningToken handles a token that is coming back to the original process.
func (p *Process) handleReturningToken(incomingDirection Direction) {
	log.Printf("Process %d: %s token came back.\n", p.uid, incomingDirection.String())
	if incomingDirection == Left {
		p.hasLeftToken = true
	} else {
		p.hasRightToken = true
	}

	if p.hasLeftToken && p.hasRightToken {
		p.phase++
		p.hasLeftToken = false
		p.hasRightToken = false

		log.Printf("Process %d: Both messages completed the round trip. Starting phase %d.\n", p.uid, p.phase)
		go p.startElectionPhase(p.phase)
	}
}

// handleOutMessage handles Out messages from another process
func (p *Process) handleOutMessage(msg Message, incomingDirection Direction) {
	log.Printf("Process %d: Received Out message from neighbor. %s\n", p.uid, msg.String())
	if msg.uid < p.uid {
		// ignore message
		log.Printf("Process %d: ignoring message with uid %d. %s\n", p.uid, msg.uid, msg.String())
	} else {
		msg.hops--
		if msg.hops == 0 {
			msg.way = In
			p.forwardMessage(msg, oppositeDirection(incomingDirection))
		} else {
			p.forwardMessage(msg, incomingDirection)
		}
	}
}

// forwardMessage forwards a message to the next process in the appropriate direction
func (p *Process) forwardMessage(msg Message, dir Direction) {
	log.Printf("Process %d: Sending message to %s neighbor. %s\n", p.uid, msg.String(), dir.String())

	if dir == Left {
		p.rightNeighbour <- msg
	} else {
		p.leftNeighbour <- msg
	}
}

// startElectionPhase starts an election phase.
func (p *Process) startElectionPhase(phase int) {
	log.Printf("Process %d: Starting election phase %d.\n", p.uid, phase)

	p.leftNeighbour <- Message{
		uid:  p.uid,
		hops: int(math.Pow(2, float64(phase))),
		way:  Out,
	}
	p.rightNeighbour <- Message{
		uid:  p.uid,
		hops: int(math.Pow(2, float64(phase))),
		way:  Out,
	}
}
