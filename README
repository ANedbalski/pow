## 1. Description
This project is an implementation of test task on Golang.

The task: is to create a TCP server that protected against DDOS attacks by implementing ProofOfWork algorithm.

## 2. Starting application
### 2.1 Requirements
+ [Go 1.21+] installed (to run tests, start server or client without Docker)
+ [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)

### 2.1 Start server and client by docker-compose:
```
make start
```

### 2.3 Build server locally:
```
make start-server
```

### 2.4 Start client locally:
```
make start-client
```

### 2.5 Launch tests:
#### 2.5.1 Launch full tests :
```
make test
```

#### 2.5.2 Launch only lightweight tests:
```
make test-short
```

## 5. Proof of Work
### 5.1 Selection of an algorithm
There is some different algorithms of Proof Work.
#### 5.1.1 Merkle tree
[Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
Advantages:
+ It can be perform updates in a specific order and to calculate root hash predictably
+ It is computationally simple and quick

Disadvantages:
+ In Merkle tree server should do too much work to validate client's solution. It is require to generate hash on each leaf


#### 5.1.2 Hashcash
[Hashcash](https://en.wikipedia.org/wiki/Hashcash)

Advantages:
+ simplicity of implementation
+ lots of documentation and articles with description
+ simplicity of validation on server side
+ possibility to dynamically manage complexity for client by changing required leading zeros count

Disadvantages:
+ Compute time depends on power of client's machine. very weak clients possibly could not solve challenge, or too powerful computers could implement DDOS-attackls. But complexity of challenge could be modified by changing the number of required zeros.
+ Pre-computing challenges in advance before DDOS-attack. Some clients could parse protocol and compute many challenges to apply all of it in one moment. It could be solved by additional validation of hashcash's params on server.

#### 5.1.3 Guided tour puzzle
 [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

Advantages:
+ can mitigate the effect of denial of service attack, because the more an attacker wants to overwhelm the server, the more puzzles it has to compute
+ enforces delay on the clients through round trip delays, so that clients' requests arrive at a rate that is sustainable by the server

Disadvantages:
+ In guided tour puzzle client should regularly request server about next parts of guide, that complicates logic of protocol.

#### 5.1.4 Cuckoo Cycle
[Cuckoo hashing](https://en.wikipedia.org/wiki/Cuckoo_hashing)

Advantages:
+ Cuckoo Cycle is immune from quantum speedup by Grover's search algorithm

Disadvantages:
+ Cuckoo hashing needs two hash functions for avoiding collision.
+ There is a possibility that collision occurs in such a way the algorithm enters into an infinite loop.

### 5.2 The algorithm selection
After comparison, I chose Hashcash.
To mitigate HashCash weaknesses I made some improvements:
+ server save **rand** value to cache and check it's existence on verify step

## 6. Assumptions and limitations:
As it is a test task, it has some limitations and aassumptions.
+ Only InMemory repository is implemented. It is possible to implement other repositories, like Redis, Postgres, etc.
+ Not all code is covered by tests. Only main logic is covered by tests.
+ Not all errors are handled by server. It is possible to implement more detailed error handling.
+ Client made to be simple. It implements only happy path and panic on errors.
+ Server vocabulary is limited. Implemented only commands required for task.


