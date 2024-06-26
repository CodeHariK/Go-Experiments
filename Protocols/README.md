# Protocol

![img](img.gif)

Inter-process communications are usually implemented using message passing with a synchronous request-response style or asynchronous event-driven styles. In the synchronous communication style, the client process sends a request message to the server process over the network and waits for a response message. In asynchronous event-driven messaging, processes communicate with asynchronous message passing by using an intermediary known as an event broker.

## HTTP/2
High-performance binary message protocol with support for bidirectional messaging.

* [Complete Redis, Websockets, Pub Subs and Message queues bootcamp](https://www.youtube.com/watch?v=IJkYipYNEtI)

## Websocket

* [ProgrammingPercy Mastering WebSockets With Go - An in-depth tutorial](https://www.youtube.com/watch?v=pKpKv9MKN-E&ab_channel=ProgrammingPercy)
* [ProgrammingPercy Mastering WebSockets With Go](https://programmingpercy.tech/blog/mastering-websockets-with-go/)
* [Going Infinite, handling 1 millions websockets connections in Go](https://www.youtube.com/watch?v=LI1YTFMi8W4)
* [How to scale WebSockets to millions of connections](https://www.youtube.com/watch?v=vXJsJ52vwAA)

* [RFC6455](https://www.rfc-editor.org/rfc/rfc6455)

* [Websocket Web Api](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

* TCP, Client-Server, Stateful
* Full duplex connection,  Scalable RealTime, Persistant connection between client and server
* Solves : Http Long polling (Resourse intensive)

* Stateful Nature: Unlike traditional HTTP, which is stateless, WebSockets are stateful. 
* This means that the server needs to maintain the connection state for each client, leading to increased memory usage and potential scalability challenges.

* WebSockets uses HTTP to send an initial request to the server. This is a regular HTTP request, but it contains a special HTTP header Connection: Upgrade. This tells the server that the client is trying to upgrade the HTTP requests TCP connection into a long-running WebSocket. If the server responds with an HTTP 101 Switching Protocols then the connection will be kept alive, making it possible for the Client and Server to send messages bidirectional, full-duplex.
* Reading and writing messages might seem like an easy task, and it is. There is however a small pitfall that many people miss. The WebSocket connection is only allowed to have one concurrent writer, we can fix this by having an unbuffered channel act as a locker.
* The HeartBeats — Ping & Pong : WebSockets allow both the Server and the Client to send a Ping frame. The Ping is used to check if the other part of the connection is still alive. But not only do we check if we other connection is alive, but we also keep it alive. A WebSocket that is idle will / can close because it has been idle for too long, Ping & Pong allows us to easily keep the channel alive avoiding long-running connections with low traffic to close unexpectedly. Whenever a Ping has been sent, the other party has to respond with a Pong. If no response is sent, you can assume that the other party is no longer alive. It’s logical right, you don’t keep talking to somebody that doesn’t respond.
* One rule of security is to always expect malicious usage. If people can, they will. So one thing that is good to always do is to limit the maximum size of a message that is allowed to be processed on the server.

* Use cases : Chat, Trading

## WebRTC

* UDP, P2P
* Scalable RealTime
* Use cases : Games, video call

## REST

* HTTP
* Inefficient text based message protocols
* Lacks strongly typed interfaces between apps

## GRPC

* HTTP/2
* Protocol buffer based binary protocol
* Request–Response and Streaming
* Authentication, encryption, resiliency (deadlines and timeouts), metadata exchange, compression, load balancing, service discovery
* Unary RPC, Server-Streaming RPC, Client-Streaming RPC, Bidirectional-Streaming RPC

* https://protobuf.dev/programming-guides/encoding/

Interprocess communication technology that allows you to connect, invoke, operate, and
debug distributed heterogeneous applications as easily as making a local function call.

RESTful services are quite bulky, inefficient, and error-prone for building inter-process communication. It is often required to have a highly scalable, loosely coupled inter-process communication technology that is more efficient than RESTful services. This is where gRPC, a modern inter-process communication style for building distributed applications and microservices, comes into the picture (we’ll compare and contrast gRPC with RESTful communication later in this chapter). gRPC primarily uses a synchronous request-response style for communication but can operate in fully asynchronous or streaming mode once the initial communication is established.
