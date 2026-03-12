# 🐜 Lem-in: Digital Ant Farm Optimization

A high-performance pathfinder written in **Go** that moves **N ants** through a network of rooms using optimized **network flow** algorithms. Given a colony description (rooms, tunnels, number of ants), the program computes the fastest possible schedule and prints each ant's movement turn by turn.

---

## ✨ Features

| Feature | Description |
|---|---|
| **Edmonds-Karp Algorithm** | BFS-based maximum flow to find the optimal set of paths. |
| **Node Splitting** | Vertex capacities enforced via virtual in/out node halves. |
| **Residual Graph BFS** | Augmenting paths discovered through forward & backward edges. |
| **Greedy Ant Distribution** | Ants assigned to minimize per-path completion time. |
| **Early Termination** | Stops augmenting when adding more paths increases total turns. |

---

## 🧠 Technical Deep Dive

### The Problem

Move **N** ants from a **start** room to an **end** room through a network of tunnels. Each room (except start and end) can hold at most **one ant per turn**. The objective is to minimize the total number of turns.

### Why Network Flow?

A greedy shortest-path approach fails when multiple short paths share bottleneck nodes. Network flow naturally discovers **vertex-disjoint** paths that avoid such conflicts.

### Node Splitting (Vertex Capacities)

Standard max-flow algorithms only handle **edge capacities**. To enforce **vertex capacity = 1** on interior nodes, we use **node splitting**:

```
Original:          Split:
   ┌───┐           ┌────┐     ┌─────┐
──▶│ v │──▶   ──▶  │ v_in│────▶│v_out │──▶
   └───┘           └────┘     └─────┘
                    capacity=1 on internal edge
```

In our implementation, we don't physically split nodes. Instead, the BFS tracks a **virtual side** for each visited node:

- **side = 0** — The node is free (no flow through it). Treated as unsplit.
- **side = 'i'** — The "in" half of a used interior node (entered via a forward edge).
- **side = 'o'** — The "out" half of a used interior node (entered via a backward/cancel edge).

This guarantees that two augmenting paths never share the same interior node.

### Residual Edges Prevent Greedy Bottlenecks

When the BFS finds a new augmenting path, it may **cancel** flow on previously used edges (traversing backward edges in the residual graph). This allows the algorithm to "undo" a suboptimal earlier path and redistribute flow for a globally better solution.

```
Iteration 1: S → A → C → E  (greedy shortest path)
Iteration 2: S → B → C̃ → A → E  (backward edge through C cancels old flow)
Result:       S → A → E  and  S → B → C → E  (two disjoint paths)
```

### Complexity

The Edmonds-Karp algorithm runs in **O(V · E²)** time, where **V** is the number of rooms and **E** is the number of tunnels. In practice, ant farm networks are small enough that this runs near-instantly.

---

## 📁 Project Structure

```
lem-in/
├── main.go          Entry point — reads file, orchestrates pipeline
├── parser.go        3-phase state machine to parse input format
├── struct.go        Core data types (Farm, Room, BFS state types)
├── solver.go        Edmonds-Karp solver with node-splitting BFS
├── simulation.go    Greedy ant distribution & turn-by-turn output
├── go.mod           Go module definition
└── test/            Sample input files
    ├── test0.txt
    ├── test1.txt
    └── test_diamond.txt
```

---

## 🚀 How to Run

### Prerequisites

- **Go 1.20+** installed ([download](https://go.dev/dl/))

### Build & Run

```bash
# Clone the repository
git clone https://github.com/bouzerda0/lem-in.git
cd lem-in

# Run directly
go run . test/test_diamond.txt

# Or build a binary first
go build -o lem-in .
./lem-in test/test_diamond.txt
```

### Input Format

```
<number_of_ants>
<room_definitions>     # name x y
<link_definitions>     # room1-room2
```

Special directives `##start` and `##end` mark the starting and ending rooms respectively.

### Example

**Input** (`test/test_diamond.txt`):
```
4
##start
S 0 0
A 1 0
B 1 1
C 2 0
##end
E 3 0
S-A
S-B
A-C
B-C
C-E
A-E
```

**Output**:
```
4
##start
S 0 0
A 1 0
B 1 1
C 2 0
##end
E 3 0
S-A
S-B
A-C
B-C
C-E
A-E

L1-A L2-B
L1-E L2-C L3-A L4-B
L3-E L4-C
```

---

## 🧪 Testing

Run against the provided test files:

```bash
go run . test/test0.txt
go run . test/test_diamond.txt
```

Verify the build is clean:

```bash
go vet ./...
go build ./...
```

---

## 👤 Author

**bouzerda0** — [Zone 01 Oujda](https://zone01oujda.ma/)

---

## 📝 License

This project was developed as part of the Zone 01 curriculum.
