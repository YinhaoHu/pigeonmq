# Frequently Asked Questions

## What is the goal of this project?

This is a study project of mine to learn how to design and implement a distributed messaging system. I try to keep the design simple enough but grasp the core concepts of a distributed messaging system.

## Why use Go instead of Rust or C++?

To focus more on the core of the distributed messaging system and less on the language itself. Go is a good choice for this project because it has a richer ecosystem and is easier to code compared to Rust or C++.

## Why use etcd instead of ZooKeeper?

**`TODO`**: Add more content here.

## Why we cannot add Poras to a ledger to scale the storage?

Background: The storage layer of PigeonMQ is scailed by creating a new ledger for a topic partition instead of adding Poras to a ledger. This question and answer can be mapped to Apache Pulsar.

We can answer this question by taking an example as intuition.

Firstly, let's say we have 3 Poras for a DistributedLedger and we have 2 entries in each of the Poras. The number represents the entry ID of the entry. Qw=3.

```text
Pora1   Pora2   Pora3
1       1       1
```

Now, we add a new Pora to the DistributedLedger and Qw does not change. We append two new entries to the DistributedLedger.

```text
Pora1   Pora2   Pora3   Pora4
1       1       1       
        2       2       2
3       3       3
```

Later, we remove Pora4 from the cluster for some reason.

```text
Pora1   Pora2   Pora3
1       1       1
        2       2 
3       3       3
```

Now, we find that the entry with ID 2 does not satisfy the Qw=3. We have some options to go from here:

- We can just do nothing but the entry with ID 2 is not durable as guaranteed. This might be unacceptable for some use cases.
- We can replicate the entry with ID 2 to another Pora to satisfy the Qw=3 by appending it to Pora1. But this would violate the entry order in the ledger.
- We can replicate the entry with ID 2 to another Pora to satisfy the Qw=3 by writing it before entry with ID 3. But this would make the performace worse because we need to read and write the entries in the Pora.

Obviously, none of the options are acceptable for a distributed messaging system like PigeonMQ. Therefore, we cannot add Poras to a ledger to scale the storage as Apache BookKeeper does.
