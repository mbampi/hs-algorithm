package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

func main() {

	numProcesses := getNumProcesses()

	processes := createProcesses(numProcesses)

	wg := sync.WaitGroup{}
	wg.Add(1) // because it should finish when leader knows it is the leader

	// Start processes
	for i := 0; i < numProcesses; i++ {
		log.Printf("Starting process %d.", i)
		go processes[i].start(&wg)
	}

	wg.Wait()

	fmt.Printf("Number of messages: %d\n", numMessages.num)
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

// getNumProcesses returns the number of processes to create,
// specified as the first command line argument.
func getNumProcesses() int {
	if len(os.Args) < 2 {
		log.Fatal("Please specify the number of processes.")
	}

	numProcesses, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Please specify the number of processes.")
	}

	return numProcesses
}
