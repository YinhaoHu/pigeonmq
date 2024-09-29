# Proof of Concept

## Introduction

In this article, we provide a proof of concept for the PigeonMQ system. We demonstrate how to publish a message to a topic and consume a message from a topic using the PigeonMQ client library. We also demonstrate how to use the PigeonMQ-Admin-CLI to administer the PigeonMQ system.

## Publish a message to a topic

In this section, we prove that the APIs PigeonMQ provides are reasonable and can be used to publish a message to a topic.

Note that the following code is a proof of concept and may not be the final API design. For demonstration purposes, we ignore error handling.

```go
etcdServers := []string{"localhost:2379"}
client, err := pigeonmq.NewClient(etcdServers)
defer client.Close()

producerOptions := &pigeonmq.ProducerOptions{
    Topic:         "my-topic",
    Partition:     0,
    EtcdEndpoints: etcdServers,
} 
producer, err := client.CreateProducer(producerOptions)

// Produce 3 messages to the topic.
for i := 0; i < 3; i++ {
    message := &pigeonmq.Message{
        Key:     fmt.Sprintf("key-%d", i),
        Payload: []byte(fmt.Sprintf("payload-%d", i)),
    }
    _, err := producer.Produce(message)
}
producer.Close()

```

## Consume a message from a topic

In this section, we prove that the APIs PigeonMQ provides are reasonable and can be used to consume a message from a topic.

For demonstration purposes, we ignore error handling.

```go
etcdServers := []string{"localhost:2379"}
client, err := pigeonmq.NewClient(etcdServers)
defer client.Close()

consumerOptions := &pigeonmq.ConsumerOptions{
    Topic:         "my-topic",
    Partition:     0,
    EtcdEndpoints: etcdServers,
}
consumer, err := client.CreateConsumer(consumerOptions)

// Consume messages from the topic.
for {
    message, err := consumer.Consume()
    if err != nil {
        break
    }
    fmt.Printf("key: %s, payload: %s\n", message.Key, string(message.Payload))
}
```

## Admin CLI

In this section, we prove that the use of PigeonMQ-Admin-CLI provides a reasonable way to administer the PigeonMQ system.

```bash
pigeonmq-admin topic create my-topic
pigeonmq-admin topic list
pigeonmq-admin topic delete my-topic
```
