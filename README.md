# PigeonMQ(V1.0.0)

## Overview

PiegonMQ is a scalable, reliable and performant messaging system insipired by Apache Pulsar.

Scaiability is achieved by separating the storage layer from the service layer. The storage layer is designed to be horizontally scalable which will automatically share the loads of the overloaded Poras controlled by Broker. The service layer is designed to be stateless and horizontally scalable by partitioning the topics and partitions across multiple Brokers.

Reliability is achieved by using the Leaderless Replication approach in the storage layer. PigeonMQ ensures that the successfully produced messages satisfies the Quorum Ack condition and eventually satisfied the Quorum Write condition if the world is happy.

High performance is achieved by using the multi level cache mechanism, batching and optimized append-only storage system.

This is a study project of mine to learn how to design and implement a distributed messaging system. I try to keep the design simple enough but grasp the core concepts of a distributed messaging system. Scailability, reliability and performance are the main goals of the system, but I intentionally ignore the security and maintenance-friendly for operators which are definitely important in a real-world system.

![PigeonMQ Architecture](./docs/diagrams/output/pigeonmq-architecture.drawio.png)

## Documentation

The documentations are available in the `docs` directory including the design, frequently asked questions, and more.

## Quick Start

## Status

PigeonMQ is currently in the development phase. The first version will be released soon. The design is complete and the implementation is in progress.
