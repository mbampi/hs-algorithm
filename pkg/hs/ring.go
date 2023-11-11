package hs

import (
	"log"
	"sync"
)

type Ring struct {
	processes []Process
}

// NewRing creates a ring of processes.
func NewRing(numProcesses int) *Ring {
	processes := createProcesses(numProcesses)
	return &Ring{
		processes: processes,
	}
}

// Start starts the ring.
func (r *Ring) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1) // because it should finish when leader knows it is the leader

	// Start processes

	for i := 0; i < len(r.processes); i++ {
		log.Printf("Starting process %d.", i)
		go r.processes[i].Run(&wg)
	}

	wg.Wait()
}

// createProcesses creates n processes and connects them in a ring,
// in a random order.
func createProcesses(numProcesses int) []Process {
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

	// Connect processes in a ring
	for i := 0; i < numProcesses; i++ {
		leftIndex := (i - 1 + numProcesses) % numProcesses
		rightIndex := (i + 1) % numProcesses

		log.Printf("Process %d: Connecting to left neighbour %d and right neighbour %d\n", processes[i].uid, processes[leftIndex].uid, processes[rightIndex].uid)
		processes[i].leftNeighbour = processes[leftIndex].right
		processes[i].rightNeighbour = processes[rightIndex].left
	}

	return processes
}
