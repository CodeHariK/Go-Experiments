# System Design

## Cache

* https://www.bytesizedpieces.com/posts/cache-types
* https://redis.io/blog/cache-eviction-strategies/

## Delegation

What does not need to be realtime, should not be realtime.

Message brokers acts as buffers for messages, that will be eventually processed by subscribers.

Queue : Expiration, Delay, Atmostonce, Atleastonce, exactlyonce, waitsuntilmessageretrived

## Database

### SQL
* Constraints : Unique, not null, primary key, foriegn key
* ACID
* Atomicity : Each transaction is single unit of execution
* Consistency : Database goes from one valid state to another
* Durability : Post commit data is recoverable even after system failure
* Isolation : Concurrent execution = Sequential execution

## Distibuted System

+ Horizontal Scalibility, Reliability, Performance
- Latency, Observbility, Consensus, Clock Synchronization, Distributed Deadlocks

* Components fail indpendent of each other, but failure cascade
* Cicuit breaker : break the circuit, so that failure won't propogated or overload system and system can self heal.

* L4 vs L7 load balancer
* L4 works on packets, does not know anything above it.
* Round robin, weighted rr, least connection, ip hash, url hash
* HealthCheck, Connection Polling, Proxy, Request Tracing, Stats, Discovery, Header Manipulation
