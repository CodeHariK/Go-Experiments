# RabbitMQ

* [https://training.cloudamqp.com/](https://training.cloudamqp.com/)

Message Queue are used to manage flood of request with unpredictable processing and response time. 
Loose coupling. Good with multiple chained recievers. Works while recievers are temporarily unavailable . Microservices.

PubSub          : Message can be consumed multiple times. Consumer has to keep record of what has been already consumed.
Message Queue   : Once a message has been consumed, it is removed from queue.
ReqRes          : Stateless, Simple, Scalable (load balancing), High Coupling Complicated Architecture, 
                  Bad with multiple chained recievers, recievers have to be always running
                  
* [Complete Config File](https://github.com/rabbitmq/rabbitmq-server/blob/main/deps/rabbit/docs/rabbitmq.conf.example)

* [AMQP 0-9-1 Model Explained](https://www.rabbitmq.com/tutorials/amqp-concepts.html)

* [AMQP 0-9-1 Complete Reference Guide](https://www.rabbitmq.com/amqp-0-9-1-reference.html)

* [VMware RabbitMQ for Kubernetes](https://docs.vmware.com/en/VMware-RabbitMQ-for-Kubernetes/1/rmq/configure.html)

* [Documentation](https://www.rabbitmq.com/documentation.html)

* [Go RabbitMQ Client Library](https://pkg.go.dev/github.com/rabbitmq/amqp091-go)

* [RabbitMQ AWS](https://github.com/rabbitmq/rabbitmq-server/tree/main/deps/rabbitmq_aws)

## Tutorial

* [Hello World](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)

## RabbitMQ

* Server Pushes the message to Consumer.
* Connection : (TCP/IP socket connection) A connection in AMQP 0.9.1 is a network connection between your application and the AMQP broker.
* Channel : (Multiplexing, Same Tcp Connection) Logical Connection in Connection, ~Mini Connections.
* Exchange : An exchange is responsible for routing messages to different queues with the help of header attributes, bindings, and routing keys.
* Queue : A queue acts as a buffer that stores messages that are consumed later.
* Binding : A binding is a relation between a queue and an exchange consisting of a set of rules that the exchange uses (among other things) to route messages to queues.
* Virtual Hosts : (Multi Tenancy) Virtual hosts (vhost) provide a way to segregate applications in the broker. Different users can have different access privileges to different vhost. Queues and exchanges is created so they only exist in one vhost.

### Exchanges

Direct, Topic, Fanout, Headers

### Nodes and Clustering

A RabbitMQ cluster is a group of one or several RabbitMQ servers, where each RabbitMQ server is sharing users, virtual hosts, queues, exchanges, bindings, runtime parameters, and other distributed states. When setting up a cluster, you set up another RabbitMQ node and join that node into the default node. There is no need to configure any users or virtual hosts on the second node. Just join the cluster and the configuration will automatically synchronize with the existing RabbitMQ instance, including users, virtual hosts, exchanges, queues, and policies.

### Quorum Queues

Quorum queues ensure that the cluster is up to date by agreeing on the contents of a queue.
All communication is routed to the queue leader, which means the queue leader locality has an effect on the latency and bandwidth requirement of the messages. In quorum queues, the leader and replication are consensus-driven, which means they agree on the state of the queue and its contents. Quorum queues will only confirm when the majority of its nodes are available, which thereby avoids data loss. If the majority of nodes agree on the contents of a queue, the data is valid. Otherwise, the system attempts to bring all queues up to date.

### RabbitMQ Streams

There are two components of stream functionality, stream queues, and the stream protocol. Stream queues are persistent and replicated, and work similarly to an append-only log. Stream queues can be used with traditional AMQP clients, i.e. they don’t have to use the stream protocol. Messages in stream queues are not removed when consumed, instead, they can be consumed over and over again using an offset or a timestamp, similar to functionalities that can be found in Apache Kafka. RabbitMQ Streams also comes with its very own stream protocol, which has shown to be much faster than AMQP in RabbitMQ.

Streams model an append-only log that's immutable. In this context, this means messages written to a Stream can't be erased, they can only be read. To read messages from a Stream in RabbitMQ, one or more consumers subscribe to it and read the same message as many times as they want. Additionally, Streams are always persistent and replicated.

Queues are limited in the following scenarios:

They deliver the same message to multiple consumers by binding a dedicated queue for each consumer. Clearly, this could create a scalability problem.
They erase read messages making it impossible to re-read(replay) them or grab a specific message in the queue.
They perform poorly when dealing with millions of messages because they are optimized to gravitate toward an empty state.

The use cases where streams shine include:

Fan-out architectures: Where many consumers need to read the same message
Replay & time-travel: Where consumers need to read and reread a message from any point in the stream.
Large Volumes of Messages: Streams are great for use cases where large volumes of messages need to be persisted, because they store messages on the file system.
High Throughput: RabbitMQ Streams process relatively higher volumes of messages per second.

#### Fan-out Architectures
A fan-out architecture is where multiple consumers read the same message. As mentioned earlier, implementing this sort of architecture with queues isn’t optimal. Having to add queues for every added consumer is resource intensive, which gets worse when dealing with queues that need to persist data.

Streams in RabbitMQ make implementing fan-out architectures a breeze. Because consumers read messages from a Stream in a non-destructive manner, a message will always be there for the next consumer to access it. In essence, to implement a fan-out architecture, just declare a RabbitMQ Stream and bind as many consumers as needed.

#### Speed
In a comparison Quorum queues handled about 40,000 messages per second. Streams, on the other hand, handled around 64,000 messages per second when used with AMQP protocol, and over 1 million messages per second when used with native Stream protocol.
Stream protocol doesn’t handle things like routing messages, de-queueing messages, etc., it technically does less work, and this translates to higher performance.

### Policies

rabbitmqctl list_policies -p vhostname
rabbitmqctl set_policy <name> <pattern> <definition>
rabbitmqctl set_policy Policy2 '.*' '{"message-ttl": 60}'
rabbitmqctl set_policy <name> <pattern> <definition> -p <vhost> --apply-o <queues|exchanges> --priority 8
rabbitmqctl set_policy Policy5 'queue_A' '{"message-ttl": 60}' -p zyxhvjpx --apply-to queues --priority 10'
rabbitmqctl clear_policy <name>

### Consumer acknowledgements

https://www.rabbitmq.com/docs/confirms

Consumer acknowledgements (acks) provide notice if and when messages between RabbitMQ and the consumer have been sent to or received and handled by the intended destination. The response of a consumer acknowledgement setting can also trigger actions that may re-queue and resend messages if needed. 

### Prefetch

prefetch value is used to specify how many messages are being sent at the same time.

Messages in RabbitMQ are pushed from the broker to the consumers. The RabbitMQ default prefetch setting gives clients an unlimited buffer, meaning that RabbitMQ, by default, sends as many messages as it can to any consumer that appears ready to accept them. It is, therefore, possible to have more than one message "in-flight" on a channel at any given moment.

Messages are cached by the RabbitMQ client library (in the consumer) until processed. All pre-fetched messages are invisible to other consumers and are listed as unacked messages in the RabbitMQ management interface.

An unlimited buffer of messages sent from the broker to the consumer could lead to a window of many unacknowledged messages. Prefetching in RabbitMQ simply allows you to set a limit of the number of unacked (not handled) messages.

There are two prefetch options available, channel prefetch count and consumer prefetch count.

The channel prefetch value defines the maximum number of unacknowledged deliveries that are permitted on a channel. Setting a limit on this buffer caps the number of received messages before the broker waits for an acknowledgment.

Because a single channel may consume from multiple queues, coordination between them is required to ensure that they don’t pass the limit. This can be a slow process especially when consuming across a cluster, and it is not the recommended approach.

The best practice is to set a consumer prefetch by setting a limit on the number of unacked messages at the client.

Optimum Consumer Prefetch Count
A larger prefetch count generally improves the rate of message delivery. The broker does not need to wait for acknowledgments as often and the communication between the broker and consumers decreases. Still, smaller prefetch values can be ideal for distributing messages across larger systems. Smaller values maintain an even rate of message consumption. A value of one helps ensure equal message distribution.

A prefetch count that is set too small may hurt performance since RabbitMQ might end up in a state where the broker is waiting to get permission to send more messages. The image below illustrates a long idling time. In the figure, we have a QoS prefetch setting of one (1). This means that RabbitMQ won't send out the next message until after the round trip completes (deliver, process, acknowledge). Round-trip time in this figure is in total 125ms with a processing time of only 5ms.

A large prefetch count, on the other hand, could take lots of messages off the queue and deliver all of them to one single consumer, keeping the other consumers in an idling state, as illustrated in the figure below:

If you have one single or only a few consumers processing messages quickly, we recommend prefetching many messages at once to keep your client as busy as possible. If you have about the same processing time all the time and network behavior remains the same, simply take the total round trip time and divide by the processing time on the client for each message to get an estimated prefetch value.

In a situation with many consumers and short processing time, we recommend a lower prefetch value. A value that is too low will keep the consumers idling a lot since they need to wait for messages to arrive. A value that is too high may keep one consumer busy while other consumers are being kept in an idling state.

If you have many consumers and/or long processing time, we recommend setting the prefetch count to one (1) so that messages are evenly distributed among all your consumers.

Avoid the usual mistake of having an unlimited prefetch, where one client receives all messages and runs out of memory and crashes, causing all the messages to be re-delivered.

### Dead Lettering

Some messages become undeliverable or unhandled even when received by the broker. This can happen when the amount of time the message has spent in a queue exceeds the time to live, TTL, when the queue reaches its capacity, or when a message is negatively acknowledged by the consumer. Such a message is called a dead message.

But don’t worry, there is no need to lose messages entirely. Setting up a RabbitMQ dead letter exchange and a dead letter queue allows orphaned messages to be stored and processed. Dead letter handling is a key part of almost any messaging system, and RabbitMQ is no different.

There are three identified situations where a message becomes undeliverable after reaching RabbitMQ:

A message is negatively acknowledged by the consumer
The TTL of a message expires
The queue reaches capacity

You can create a redelivery mechanism using dead letter queues by attaching the original exchange to the dead letter queue. However, it is possible to create an endless cycle of redelivered messages which could clog dead letter queues if left unchecked.

Create a separate dead letter queue that stores messages before pushing them to an exchange or routing them back to the original queue.

To further improve the system, include an incremented property in the message body indicating the number of times the message was received. This requires handling dead letters in a separate consumer but allows you to eventually drop messages or push them to storage.

### Alternate Exchange
No matter how careful you are, mistakes can happen. For example, a client may accidentally or maliciously route messages using non-existent routing keys. To avoid complications from lost information, collecting unroutable messages in a RabbitMQ alternate exchange is an easy, safe backup.

### Queue Performance Optimization

Use Quorom Queues
Keep Queues Short
Enable Lazy Queues for Long Queues
Limit Queue Size With Max-length or TTL
Consider the Number of Queues
Split Queues Over Different Cores
Don’t Set Names on Temporary Queues
Auto-delete Unused Queues
Set Limited Use on Priority Queues
Consider the Payload

Don’t Use Too Many Connections or Channels
Don’t Share Channels Between Threads
Don’t Open and Close Connections or Channels Repeatedly
Separate Connections for Publisher and Consumer
Consider VPC Peering

Acknowledge Messages According to Use Case
Limit Unacknowledged Messages
Use Persistent Messages and Durable Queues To Ensure You Don't Loose Any Messages
Use Transient Messages and Non-durable Queues for High Performance
Consider Prefetch Size
Consider the Exchange Type
Disable Plugins That You Are Not Using
Don’t Set Management Statistics Rate as Detailed

