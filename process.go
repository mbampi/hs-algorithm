package main

import (
	"fmt"
	"log"
	"math"
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

// Direction represents the direction of a received message.
type Direction int

const (
	// Left represents a message coming from the left neighbor.
	Left Direction = iota
	// Right represents a message coming from the right neighbor.
	Right
)

func (p *Process) start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Start the election
	go p.electionPhase(0)

	end := false
	for !end {
		select {
		case msg := <-p.left:
			// Process the message from the left neighbor
			numMessages.increment()
			end = p.processMessage(msg, Left)
		case msg := <-p.right:
			// Process the message from the right neighbor
			numMessages.increment()
			end = p.processMessage(msg, Right)
		}
	}
}

// processMessage processes a message received from a neighbor.
func (p *Process) processMessage(msg Message, incomingDirection Direction) bool {
	if msg.uid == p.uid {
		if msg.way == Out {
			fmt.Printf("Process %d: I am the leader!\n", p.uid)
			return true
		} else {
			if incomingDirection == Left {
				log.Printf("Process %d: Left message completed the round trip.\n", p.uid)
				p.hasLeftToken = true
			} else {
				log.Printf("Process %d: Right message completed the round trip.\n", p.uid)
				p.hasRightToken = true
			}

			if p.hasLeftToken && p.hasRightToken {
				p.phase++
				log.Printf("Process %d: Both messages completed the round trip. Starting phase %d.\n", p.uid, p.phase)
				p.hasLeftToken = false
				p.hasRightToken = false
				go p.electionPhase(p.phase)
			}
		}
	} else { // msg.uid != p.uid
		if msg.way == Out {
			log.Printf("Process %d: Received Out message from neighbor. msg=%+v\n", p.uid, msg)
			if msg.uid < p.uid {
				// ignore message
				log.Printf("Process %d: ignoring message with uid %d. msg=%+v\n", p.uid, msg.uid, msg)
			} else {
				msg.hops--
				if msg.hops == 0 {
					msg.way = In
					if incomingDirection == Left {
						log.Printf("Process %d: Sending In message to left neighbor. msg=%+v\n", p.uid, msg)
						p.leftNeighbour <- msg
					} else {
						log.Printf("Process %d: Sending In message to right neighbor. msg=%+v\n", p.uid, msg)
						p.rightNeighbour <- msg
					}
				} else {
					if incomingDirection == Left {
						log.Printf("Process %d: Sending Out message to right neighbor. msg=%+v\n", p.uid, msg)
						p.rightNeighbour <- msg
					} else {
						log.Printf("Process %d: Sending Out message to left neighbor. msg=%+v\n", p.uid, msg)
						p.leftNeighbour <- msg
					}
				}
			}
		} else { // msg.way == In
			log.Printf("Process %d: Received In message from neighbor. msg=%+v\n", p.uid, msg)
			if incomingDirection == Left {
				log.Printf("Process %d: Sending In message to right neighbor. msg=%+v\\n", p.uid, msg)
				p.rightNeighbour <- msg
			} else {
				log.Printf("Process %d: Sending In message to left neighbor. msg=%+v\\n", p.uid, msg)
				p.leftNeighbour <- msg
			}
		}
	}

	return false
}

// electionPhase starts an election phase.
func (p *Process) electionPhase(phase int) {
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
