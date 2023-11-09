package main

import (
	"fmt"
	"hsalgorithm/pkg/hs"
	"hsalgorithm/pkg/stats"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {

	numProcesses := getNumProcesses()

	processes := hs.CreateProcesses(numProcesses)

	wg := sync.WaitGroup{}
	wg.Add(1) // because it should finish when leader knows it is the leader

	// Start processes
	for i := 0; i < numProcesses; i++ {
		log.Printf("Starting process %d.", i)
		go processes[i].Run(&wg)
	}

	wg.Wait()

	fmt.Printf("Number of messages: %d\n", stats.NumMessages.Get())
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
