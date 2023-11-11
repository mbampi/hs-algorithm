package main

import (
	"fmt"
	"hsalgorithm/pkg/hs"
	"hsalgorithm/pkg/stats"
	"log"
	"os"
	"strconv"
)

func main() {

	numProcesses := getNumProcesses()

	ring := hs.NewRing(numProcesses)

	ring.Run()

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
