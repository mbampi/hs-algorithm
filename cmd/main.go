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

	// read flag to determine whether to run the algorithm or the test
	if len(os.Args) > 2 {
		if os.Args[2] == "test" {
			avgMessages := testAlgorithm(numProcesses)
			fmt.Printf("Average number of messages: %d\n", avgMessages)
			return
		}
	}

	runAlgorithm(numProcesses)
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

// runAlgorithm runs the algorithm once and prints the number of messages.
func runAlgorithm(numProcesses int) {
	ring := hs.NewRing(numProcesses)

	ring.Run()

	fmt.Printf("Number of messages: %d\n", stats.NumMessages.Get())
}

// testAlgorithm runs the algorithm 5 times and returns the average number of messages.
func testAlgorithm(numProcesses int) int {
	numTests := 5
	totalMessages := 0

	for i := 0; i < numTests; i++ {
		ring := hs.NewRing(numProcesses)
		ring.Run()

		totalMessages += stats.NumMessages.Get()
		stats.NumMessages.Reset()
	}

	return totalMessages / numTests
}
