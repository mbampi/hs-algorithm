# Hirschberg–Sinclair Algorithm

## Description

This is a simple implementation of the Hirschberg–Sinclair algorithm for ring leader election in a distributed system,
based on Nancy A. Lynch's book "Distributed Algorithms".

The algorithm is implemented in Go. 
Each node is represented by a goroutine, and the communication between nodes is done using channels.

The algorithm is run with a given number of nodes, and each node is assigned a random ID.
The nodes then run the algorithm to elect a leader, and the leader is printed to the console.

## Usage

```bash

go run ./... <num_nodes>

```

## Example

To run the algorithm with 10 nodes, run the following command:
```bash

go run ./... 10

```


