# Research on Distributed Messaging System

## Introduction

In this article, I write about the research I have done on distributed messaging systems when I was designing PigeonMQ and the lessons I have learned from them.

## Basic Problems of Distributed Messaging Systems

### Scalability

The pressure of a distributed messaging system can be divided into two parts: the pressure of the broker and the pressure of the storage.

**Broker Pressure**: The broker pressure is the pressure of the broker to handle the messages. PigeonMQ uses the partition on message keys to distribute the messages to different brokers and we don't need to migrate the messages between brokers. This can reduce the pressure of a single broker. This idea is inspired by Apache Pulsar.

**Storage Pressure**: The storage pressure is the pressure of the storage to store the messages. PigeonMQ uses Distributed Porage as the storage layer. Distributed Porage takes a leaderless approach to achieve dynamic membership management and striping replication to achieve high scalability. This can reduce the pressure of a single storage node. This idea is inspired by Apache BookKeeper.

## Consumption Model: Push vs Pull

### Push Model

**Definition**: In the push model, the broker pushes the messages to the consumers. The broker is responsible for the message delivery.

**Pros**:

- Immediate delivery: The messages are delivered to the consumers immediately.
- Simple Consumer: The consumers don't need to poll the messages.

**Cons**:

- Consumer Overload: If the consumers are slow, the broker needs to handle the overload. Otherwise, the broker may be overwhelmed.

### Pull Model

**Definition**: In the pull model, the consumers pull the messages from the broker. The consumers are responsible for the message delivery.

**Pros**:

- Consumer Control: The consumers can control the message delivery. And overload can be avoided.

**Cons**:

- Delayed Delivery: The messages may be delayed if the consumers are slow.
