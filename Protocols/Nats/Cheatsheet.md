# Nats

* nats context ls
* nats context save cluster_sys
* nats context edit cluster_sys
* nats context select cluster_sys

* nats server ping
* nats server info
* nats server ls

* nats-server -c server.conf
* nats-server -c server.conf -name na_1

* nats bench hello --sub=10

* nats stream ls
* nats sub --stream orders --new

* nats pub hello.world --count=-1 --sleep=10ms
* nats pub hello.world {{.Count}} --count 1000
* nats reply hello.world "Hi"
* nats req hello.world ""
* nats sub hello.world