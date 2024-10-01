# Porage

Porage is the storage layer component of PigeonMQ, written in Golang. It follows a read-write separation architecture and adopts an external coordination consensus algorithm. This design ensures that Porage nodes, called **Poras**, do not communicate directly but are coordinated by the PigeonMQ Client, enabling lock-free operations using Golang's channel mechanism.

Porage's design is inspired by **Apache BookKeeper** but introduces its own optimized workflow, aiming to provide a robust storage solution for PigeonMQ.

## Key Features

- **Read and Write Separation**: Efficiently handles both read and write operations without direct locking, enhancing performance and scalability.
- **Crash Recovery**: Automatically recovers ledger data after crashes or restarts.
- **Channel-Based Coordination**: Pora nodes use channels for communication, enabling lock-free operations within a node.
- **Custom Journal, Index, and MemTable**: Implements a high-performance journal and indexing system for ledger entries.

## Core Concepts

### Ledger

A **Ledger** in Porage is a collection of entries, uniquely identified by a **LedgerID**. It can either be in an **Open** or **Closed** state. Open ledgers accept write requests, while closed ledgers only allow read operations.

### Data Structures

- **LedgerEntry**: Composed of `LeaderID`, `EntryID`, and `Payload`. In this version, error correction codes are not needed.
- **JournalEntry**: Contains `SequenceID` and `LedgerEntry`, used for write-ahead logging (WAL).
- **Private Workspace**: Lock-free workspace for each ledger, consisting of:
  - **PrivateMemTable**: Stores entries in memory.
  - **PrivateIndexFile**: Index file to track entry locations in the log.
  - **PrivateEntryLogger**: Sequential log storing the ledger entries.

## Write Path in Porage

1. **Network Goroutines**: Multiple goroutines handle incoming write requests.
2. **Journal Channel**: Acts as a buffer, allowing sequential processing of write requests.
3. **Journal Goroutine**: Writes data to the journal and routes it to the appropriate ledger channels.
4. **Ledger Channels**: Data is sent to specific ledgers for storage.
5. **Private Workspace**: Each ledger has its own workspace for efficient, lock-free storage operations.

## Read Path in Porage

- Entries are first looked up in the **PrivateMemTable** (memory).
- If not found, they are located using the **PrivateIndexFile** (disk index).

## Components

- **Journal**: Write-Ahead Log (WAL) for storing ledger entries in segments.
- **MemTable**: In-memory map for fast access to ledger entries.
- **Index**: Efficient mapping from `EntryID` to the position in the `EntryLogger`.
- **EntryLogger**: Sequential file for persistent entry storage.
