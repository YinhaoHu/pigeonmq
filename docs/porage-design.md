# Porage design Version 1.0

## Introduction

Porage is the storage layer component of PigeonMQ written in Golang. A single Porage node is called Pora. Porage adopts an external coordination consensus algorithm in terms of distribution, which means that Poras do not communicate with each other, but are coordinated by the Client. Storage adopts a read-write separation architecture. Thanks to Channel, Pora can basically achieve lock-free.

In terms of design, Porage is similar to Apache BookKeeper.

In this design about a single pora, we need to reach only the following goals in functionality:

- Read and write of a ledger.
- Recovery of a ledger after a crash or restart.

## Glossary

The core data structures of Pora are:

**LedgerEntry:** It consists of three parts: LeaderID, EntryID, and Payload. Under the current assumption of PigeonMQ, we do not need an error correction code to check whether a single stored data is damaged.

**JournalEntry:** It consists of two parts: SequenceID and LedgerEntry.

The workers in Pora are:

**NetworkWorker:** There are multiple NetworkWorkers in the entire Pora, which are responsible for receiving write requests from the network and sending them to the JournalWorker.

**JournalWorker**: There is only one JournalWorker in the entire Pora, which is responsible for receiving write requests, GroupCommit to disk, and sending the entry to be written to the corresponding LedgerChannel.

**LedgerWorker:** Receives the entry from the LedgerChannel it is responsible for, and then writes it to its own PrivateMemTable, PrivateIndexFile, and PrivateEntryLogger.

**PrivateWorkspace:** An independent workspace of a Ledger, which is the key foundation for lock-free. A PrivateWorkspace includes PrivateMemtable, PrivateIndexFile and PrivateEntryLogger.

**PrivateMemTable:** Stores the entries in memory of a Ledger.

**PrivateIndexFIle:** An IndexFile of a Ledger, used to find the location of an item in the Private Entry Logger.

**PrivateEntryLogger:** An EntryLogger of a Ledger， which stores the entries of a Ledger sequentially.

## Key Concepts

**Ledger:** A Ledger is a collection of entries. Each Ledger has a unique ID.
A ledger can be in the state of Close or Open. When a ledger is in the state of Open, it can receive write requests. When a ledger is in the state of Close, it can only be read.

## Read and write logic

In Porage, the write path is:

Based on the diagram provided for the write path in `Porage`, here’s a description of the write path:

### Write Path in `Porage`

1. **Network Goroutines**:
   - Multiple network goroutines (e.g., Network Goroutine 1, Network Goroutine 2, ..., Network Goroutine N) handle incoming write requests. Each goroutine receives data and sends it to the `Journal Channel`.

2. **Journal Channel**:
   - The data from all the network goroutines is funneled into the `Journal Channel`. This channel acts as a buffer, allowing the `Journal Goroutine` to process the incoming data sequentially.

3. **Journal Goroutine**:
   - The `Journal Goroutine` is responsible for handling the data from the `Journal Channel`. It processes each entry and then routes the data to specific `Ledger Channels`.
   - Journal should be trimmed based on the progress of the flush of the Ledger. But in this version(1.0), JournalTrim is ignored for simplicity.

4. **Ledger Channels**:
   - For each write operation, the `Journal Goroutine` directs the data to a specific `Ledger Channel` (e.g., Ledger i Channel, Ledger j Channel, ..., Ledger y Channel). These channels correspond to different ledgers within the system, organizing the data appropriately.

5. **Private Workspace**:
   - Within each node, there is a `Private Workspace` where the actual data storage operations occur. The workspace contains several workers:
     - **Local MemTable Worker**: This worker handles in-memory data structures (`MemTables`) for fast access and organizes data before it is written to persistent storage.
     - **Local EntryLogger Worker**: This worker is responsible for logging entries to disk, ensuring data durability.
     - **Local Index Worker**: This worker updates indices to allow for efficient data retrieval.
   - The `Local MemTable Worker` also receives notification channel named `NotifyChannel`, which is created every time by a `Network Goroutine`. This is used to notify the corresponding `Network Goroutine` that the write operation has been completed.

This sequence of operations effectively captures the write path within a `Pora` node in `Porage`. Data flows from the network goroutines, through the journaling process, into specific ledgers, and is eventually handled by different workers within the node for storage and indexing.

In Porage, the read path is:

1. Find Entry from PrivateMemTable.

2. If there is no Entry in PrivateMemTable, find Entry from PrivateIndexFile.

## Component

Note that Porage assumes that the data will not be damaged during the storage process. And as long as the data is written to the disk, all of the positions of the entries will be remembered in the IndexFile which means there is no need to recovery the index information from the EntryLogger.

### Journal

All WALs written by Ledger to Entry are written in a single Journal logical file, which is divided into several physical files. We call such a physical file a JournalSegment. Only one JournalSegment is open at any time.

Porage implements Journal itself.

### MemTable

It is basically a map of EntryID to a list of entries. It is used to store the entries in memory of a Ledger.

### Index

It is basically a map of EntryID to the position of the entry in the EntryLogger. It is used to find the location of an item in the EntryLogger efficiently.

Porage implements Index with the help of [Badger](https://github.com/dgraph-io/badger).

### EntryLogger

It is a file that stores the entries of a Ledger sequentially. It is used to store the entries of a Ledger. Pre-allocation mechanism is used to avoid the fragmentation of the file.

## Algorithms

### Recovery

This section explains the basic ideas of recovery algorithm in a single Pora.

**Requirement**:

- The index can be installed only if the entry is confirmed to be flushed to disk.

**Process**:

1. Use the Index to find the confirmed last persistent entry in the EntryLogger. We name it as `LastPersistentEntry`.
2. Read the journal and apply from that entry. Use `WriteAt`.
3. Recover all the ledgers.
4. Delete the journal file.
5. Pora starts to serve.

### Trim Journal

This section explains the basic ideas of trimming the journal in a single Pora.

**Requirement**:

- A periodically wake-up worker to trim the journal which is named `Trimmer`.
- Ledger worker should provide an API to get the `lastFlushTime`.

**Process**:

1. `Trimmer` wakes up.
2. Get all the `lastFlushTime` of all the ledgers.
3. Compare all  `lastFlushTime`s with the timestamp in the name of the journal file.
4. Delete the journal file if the timestamp is smaller than the minimum timestamp of the last confirmed persistent entries.
5. Go back to sleep.

## Code Guideline

For every component in Porage, we have basically the following files:

- `api.go`: Contains the functions exposed to `ledger`.
- `startup.go`: Set the configurations of the component.
- `worker.go`: The main logic of the component.
- `storage.go`: The data structure of the component.

Package `ledger` is the main package of Porage. Package structrue is as follows:

``` text
ledger
├──entry logger
├──index
├──memtable
├──journal
```
