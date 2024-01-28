# Protocol

![img](img.gif)

## Websocket

* [ProgrammingPercy Mastering WebSockets With Go - An in-depth tutorial](https://www.youtube.com/watch?v=pKpKv9MKN-E&ab_channel=ProgrammingPercy)
* [ProgrammingPercy Mastering WebSockets With Go](https://programmingpercy.tech/blog/mastering-websockets-with-go/)

* [RFC6455](https://www.rfc-editor.org/rfc/rfc6455)

* [Websocket Web Api](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

* TCP, Client-Server, Stateful
* Full duplex connection,  Scalable RealTime, Persistant connection between client and server
* Solves : ~~Http Long polling : Resourse intensive~~

* Stateful Nature: Unlike traditional HTTP, which is stateless, WebSockets are stateful. 
* This means that the server needs to maintain the connection state for each client, leading to increased memory usage and potential scalability challenges.

* WebSockets uses HTTP to send an initial request to the server. This is a regular HTTP request, but it contains a special HTTP header Connection: Upgrade. This tells the server that the client is trying to upgrade the HTTP requests TCP connection into a long-running WebSocket. If the server responds with an HTTP 101 Switching Protocols then the connection will be kept alive, making it possible for the Client and Server to send messages bidirectional, full-duplex.
*Reading and writing messages might seem like an easy task, and it is. There is however a small pitfall that many people miss. The WebSocket connection is only allowed to have one concurrent writer, we can fix this by having an unbuffered channel act as a locker.

* Use cases : Chat, Trading

## WebRTC

*  UDP, P2P
* Scalable RealTime
* Use cases : Games, video call

## WebTransport