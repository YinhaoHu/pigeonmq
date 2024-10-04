# Algorithms

This document explains the algorithms used in PigeonMQ.

## Cluster Start

1. Each Pora needs to register its information in `etcd`.
2. Each Broker needs to register its information in `etcd`.

## Runtime Cron Jobs

1. Each Pora needs to renew its load information in `etcd` periodically.
2. Each Broker needs to renew its load information in `etcd` periodically.
3. Each Broker needs to do a retention job to remove the old messages by consulting `etcd` and delete them in the Poras.

## Create a Topic Partition

Let's call the node that runs PigeonMQ-Admin-CLI as the `Admin Node`.

1. `Admin Node` chooses a Broker based on the node selection algorithm.
2. `Admin Node` sends a request to the Broker to create a topic partition.
3. Broker creates a topic partition in `etcd`.
4. Broker chooses Poras based on the node selection algorithm and sends a request to the Poras.

**Success Return Condition**: Broker creates a topic partition in `etcd`.

**Note**: PigeonMQ does not provide the ability to reconfigure the emsemble size and quorum size after the topic partition is created. This is a design decision to simplify the system. If the user wants to change the ensemble size or quorum size, the user should delete the topic partition and create a new one.

### Failure Handling

1. What if the Broker fails after step 3?
    - Other Brokers can detect the failure of the Broker and a new broker will take the ownership of that partition based on the node selection algorithm. Then, the new broker should see the number of Poras that the previous broker has chosen. If the number does not match the Ensemble size, the new broker should choose the remaining Poras based on the node selection algorithm.

## Produce a Message

This section explains the process of producing a message to a topic partition. Here is the process:

1. Producer sends a write request to the Broker which owns the partition.
2. Broker forwards the write request to the Qw Poras and waits for Qa Poras to respond.
3. If Qa Poras respond, Broker sends success to the Producer.

**Exactly-Once Not Guaranteed**: PigeonMQ does not guarantee exactly-once semantics. It is the responsibility of the Producer to handle the duplication of messages.

Duplicate messages can occur when Producer retries the message because of the timeout. But the message is already written to the Poras. The response from Broker is lost.

Idempotence of `AppendEntryWithID` should be provided by Porage.

### Failure Handling

1. What if the Broker fails after step 1 or Producer can not respond in time?

    Producer just needs to retry the write request to the Broker which owns the partition after a timeout. The available Broker may come back online and the partition is still available.

2. What if the Broker failed to receive the response from Qa Poras in time?

    Broker can directly respond to the Producer with a failure message. Producer can retry the write request to the Broker which owns the partition.

## Consume a Message

This section explains the process of consuming a message from a topic partition. Here is the process:

1. Consumer sends a read request to the Broker which owns the partition.
2. Broker checks there is newer `LastConfirmedMessage` than the `LastReadMessage` of the Consumer.
3. If there is no newer message, Broker waits for the message to arrive.
4. Broker checks whether the message is available in memory.
   - If the message is available in memory, Broker sends the message to the Consumer.
   - If the message is not available in memory, Broker sends the read request to the Qr Poras and waits for the responses.

PigeonMQ uses the `Pull Model` for message consumption. The consumer should store the `LastReadMessage` and send the `LastReadMessage` to the Broker in the read request.

Currently, PigeonMQ does not provide consumer fault tolerance. But this would be a good feature to add in the future. Basic idea is to store the `LastReadMessage` in `etcd` with watch mechanism and lease.

### Failure Handling

1. What if the Broker fails after step 1 or Consumer can not respond in time?

    Consumer just needs to retry the read request to the Broker which owns the partition after a timeout. The available Broker may come back online and the partition is still available.

2. What if the Broker failed to receive the response from Qr Poras in time?

    Broker can directly respond to the Consumer with a failure message. Consumer can retry the read request to the Broker which owns the partition.

## Delete a Topic Partition

This section explains the process of deleting a topic partition. Here is the process:

1. `Admin Node` sends a request to the Broker to delete a topic partition.
2. Broker sets the partition status to `DELETING` in `etcd`.
3. Broker sends a request to the Poras to delete the topic partition ledgers.
4. If all Poras respond in time, Broker deletes the partition in `etcd`. Otherwise, Broker responds with a failure message indicating that the partition is not deleted but status is `DELETING`.

**Success Return Condition**: Broker deletes the partition in `etcd`.

### Failure Handling

1. What if the Broker fails after step 2?

    Other Brokers can detect the failure of the Broker and a new broker will take the ownership of that partition based on the node selection algorithm. Then, the new broker should see the partition status in `etcd` and if the status is `DELETING`, the new broker should delete the partition in `etcd`.

2. What if the Broker failed to receive the response from Poras in time?

    Broker can directly respond to the `Admin Node` with a failure message. `Admin Node` can retry the delete request to the Broker which owns the partition. This marking-deleting design simplifies the process and there will be no messy state in the system.
    >Messy State: When a new broker finds that there are no enough poras, it will select new Poras or some Poras have deleted ledgers and some have not.

## Crash Detection

### Among the Poras or Brokers

Each Pora/Broker watches the key prefix related to all the partitions. If the key is a lease key related to a Pora/Broker and the event type is `DELETE`, the Pora/Broker is considered as crashed. Then a new Pora/Broker is chosen to replace the crashed one based on the node selection algorithm.

## Node Selection

### Selection for Broker and Poras in creating a topic partition

**`Least-Loaded`**: Currentlty, PigeonMQ uses `least-loaded` as the node selection algorithm for choosing a Broker to create a topic partition and choosing the Poras for a topic partition. This decision is made for simplicity and study purposes.

Maybe PigeonMQ can provide a way to customize the node selection algorithm in the future including `round-robin`, `random`, `user-defined`, etc.

### Selection for new Broker and Pora in case of Broker failure

**`Least-Loaded + Lock`**: All Brokers firstly get the load information of all Brokers and then choose the least-loaded Broker as the new Broker. If a broker thinks it is the winner, it atomically got the ownership of the partition.
