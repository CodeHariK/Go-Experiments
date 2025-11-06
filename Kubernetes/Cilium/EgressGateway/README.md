# [Cilium Egress Gateway](https://isovalent.com/labs/cilium-egress-gateway/)

## ⎈ Kubernetes Networking
Kubernetes changes the way we think about networking. In an ideal Kubernetes world, the network would be entirely flat and all routing and security between the applications would be controlled by the Pod network, using Network Policies.

## ➡️ Reaching out
In many Enterprise environments, though, the applications hosted on Kubernetes need to communicate with workloads living outside the Kubernetes cluster, which are subject to connectivity constraints and security enforcement. Because of the nature of these networks, traditional firewalling usually relies on static IP addresses (or at least IP ranges). This can make it difficult to integrate a Kubernetes cluster, which has a varying —and at times dynamic— number of nodes into such a network.

## ⬢ Cilium Egress Gateway
Cilium’s Egress Gateway feature changes this, by allowing you to specify which nodes should be used by a pod in order to reach the outside world. Traffic from these Pods will be Source NATed to the IP address of the node and will reach the external firewall with a predictable IP, enabling the firewall to enforce the right policy on the pod.

![cilium](https://play.instruqt.com/assets/tracks/wm9tp6yqnexf/1e2b29dd1c755301e5d5432343d027a0/assets/egressip.png)

[![Watch the video](https://img.youtube.com/vi/vaA6wPFxZ4Q/maxresdefault.jpg)](https://youtu.be/vaA6wPFxZ4Q)

