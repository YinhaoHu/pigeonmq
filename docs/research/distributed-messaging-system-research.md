# Research on Distributed Messaging System

## Introduction

In this article, I write about the research I have done on distributed messaging systems when I was designing PigeonMQ and the lessons I have learned from them.

## Basic Problems of Distributed Messaging Systems

### Scalability

The pressure of a distributed messaging system can be divided into two parts: the pressure of the broker and the pressure of the storage.

**Broker Pressure**: The broker pressure is the pressure of the broker to handle the messages. PigeonMQ uses the partition on message keys to distribute the messages to different brokers and we don't need to migrate the messages between brokers. This can reduce the pressure of a single broker. This idea is inspired by Apache Pulsar.

**Storage Pressure**: The storage pressure is the pressure of the storage to store the messages. PigeonMQ uses Distributed Porage as the storage layer. Distributed Porage takes a leaderless approach to achieve dynamic membership management and striping replication to achieve high scalability. This can reduce the pressure of a single storage node. This idea is inspired by Apache BookKeeper.
