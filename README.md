# Hirschberg–Sinclair Algorithm

## Description

This is a simple implementation of the Hirschberg–Sinclair algorithm for ring leader election in a distributed system,
based on Nancy A. Lynch's book "Distributed Algorithms".

The algorithm is implemented in Go. 
Each node is represented by a goroutine, and the communication between nodes is done using channels.

The algorithm is run with a given number of nodes, and each node is assigned a random ID.
The nodes then run the algorithm to elect a leader, and the leader is printed to the console.

Note: is important to note that this is implemented asynchonously, different from the book, where the algorithm is implemented synchronously. It can cause some differences in the number of messages sent by each node, and eventually can fail the execution.

## Usage

```bash

go run cmd/main.go <num_nodes>

```

## Example

To run the algorithm with 10 nodes, run the following command:
```bash

go run cmd/main.go 10

```

## Modification

On this repository you will find two variations of the original algorithm. They are separeted in different branches.

### main

This is the original algorithm, as described in the book.

### early-finding

This is a variation of the original algorithm, where the algorithm stops as soon as one (any) node finds out that who is the leader.

### all-should-know

This is a variation of the original algorithm, where all nodes should know who is the leader before the algorithm stops.

### early-all-should-know

This is a variation of the early-finding algorithm, where all nodes should know who is the leader before the algorithm stops.

### comparisons

Comparison of the number of messages sent by each algorithm, for different number of nodes.
The number of messages is the sum of messages sent by all nodes.
As the algorithm is randomized, the number of messages is the average of 5 runs.

| Algorithm             | 5 nodes | 10 nodes | 20 nodes | 40 nodes | 80 nodes | 160 nodes |
| --------------------- | ------- | -------- | -------- | -------- | -------- | --------- |
| main                  | 50      | 132      | 324      | 765      | 1800     | 4143      |
| early-finding         | 30      | 90       | 218      | 587      | 1471     | 3338      |
| all-should-know       | 55      | 136      | 336      | 793      | 1863     | 4329      |
| early-all-should-know | 34      | 95       | 254      | 624      | 1446     | 3583      |
