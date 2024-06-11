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
