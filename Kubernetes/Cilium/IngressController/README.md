# [Cilium Ingress Controller](https://isovalent.com/labs/cilium-ingress-controller/)

💡 Getting Started with Cilium Ingress Controller
Your lab environment is currently being set up. Stay tuned!

In the meantime, let's review what we will go through today:

How to deploy Cilium Service Mesh
Traffic Management & Load-Balancing with Ingress
Enable L7 Observability
Configuration
Please click the arrow on the right to proceed.

↔️ What is a Service Mesh?
With the introduction of distributed applications, additional visibility, connectivity, and security requirements have surfaced.

Application components communicate over untrusted networks across cloud and premises boundaries, load-balancing is required to understand application protocols, resiliency is becoming crucial, and security must evolve to a model where sender and receiver can authenticate each other’s identity. In the early days of distributed applications, these requirements were resolved by directly embedding the required logic into the applications.

A service mesh extracts these features out of the application and offers them as part of the infrastructure for all applications to use and thus no longer requires to change each application.

![cilium](https://play.instruqt.com/assets/tracks/459t2bcbwe0p/00fa76b301d60f122eb9882c4ea10e87/assets/cilium_service_mesh.png)

💬 Resilient Connectivity
Service to service communication must be possible across boundaries such as clouds, clusters, and premises. Communication must be resilient and fault tolerant.

🌐 L7 Traffic Management
Load balancing, rate limiting, and resiliency must be L7-aware (HTTP, REST, gRPC, WebSocket, …).

🪪 Identity-based Security
Relying on network identifiers to achieve security is no longer sufficient, both the sending and receiving services must be able to authenticate each other based on identities instead of a network identifier.

🛰️ Observability & Tracing
Observability in the form of tracing and metrics is critical to understanding, monitoring, and troubleshooting application stability, performance, and availability.

🪟 Transparency
The functionality must be available to applications in a transparent manner, i.e. without requiring to change application code.

🤹🏼‍♂️ Embedded Envoy Proxy
Cilium already uses Envoy for L7 policy and observability for some protocols, and this same component is used as the sidecar proxy in many popular Service Mesh implementations.

So it's a natural step to extend Cilium to offer more of the features commonly associated with Service Mesh —though contrary to other solutions, without the need for any pod sidecars.

Instead, this Envoy proxy is embedded with Cilium, which means that only one Envoy container is required per node.

🐝 eBPF acceleration
In a typical Service Mesh, all network packets need to pass through a sidecar proxy container on their path to or from the application container in a Pod.

This means each packet traverses the TCP/IP stack three times before even leaving the Pod.

![cilium](https://cilium.io/static/45e89165c9b995b3a8b8581e4a97711c/10c02/sidecarless.webp)

In Cilium Service Mesh, we’re moving that proxy container onto the host and kernel so that sidecars for each application pod are no longer required.

Because eBPF allows us to intercept packets at the socket as well as at the network interface, Cilium can dramatically shorten the overall path for each packet.

[![Watch the video](https://img.youtube.com/vi/lZskwr3uXn8/maxresdefault.jpg)](https://youtu.be/lZskwr3uXn8)


## Welcome to the Ingress Controller Lab

Let's have a look at our lab environment and see if Cilium has been installed correctly. The following command will wait for Cilium to be up and running and report its status:

```sh
cilium status --wait
In this lab's environment, Cilium has been installed using the cilium CLI and the following flags:

--set ingressController.enabled=true
Verify that Cilium was started with the Ingress Controller feature:

cilium config view | grep ingress-controller
```

## 🌐 Ingress Resources
Kubernetes provides a standard Ingress resource type to configure L7 load-balancing and traffic management. Cilium automatically implements Ingress in the cluster and performs the configured load-balancing configuration.

In most clusters, applying Ingress resources requires to install an Ingress Controller, using e.g. Nginx, Traefik, or Contour.

In this lab, we are using Cilium to manage Ingress resources without a need for an external controller, so you won't need to choose an Ingress Controller provider, or have to care about keeping it up-to-date!

Let's start with HTTP ingress resources!

📚 The bookinfo Application
In this challenge, we will use bookinfo as a sample application.

This demo set of microservices provided by the Istio project consists of several deployments and services:

🔍 details
⭐ ratings
✍ reviews
📕 productpage
We will use several of these services as bases for our Ingresses.

🔛 We need a Load Balancer
The Cilium Service Mesh Ingress Controller requires the ability to create LoadBalancer Kubernetes services.

Since we are using Kind on a Virtual Machine, we do not benefit from an underlying Cloud Provider's load balancer integration.

For this lab, we will use Cilium's own LoadBalancer capabilities to provide IP Address Management (IPAM) and Layer 2 announcement of IP addresses assigned to LoadBalancer services.

## Let's deploy the sample application in the cluster:

kubectl apply -f /opt/bookinfo.yml
Check that the application is properly deployed:

kubectl get pods
You should see multiple pods being deployed in the default namespace. Wait until they are Running.

Notice that with Cilium Service Mesh there is no Envoy sidecar created alongside each of the demo app microservices. With a sidecar implementation the output would show 2/2 READY: one for the microservice and one for the Envoy sidecar.


🚪Deploy the ingress
Inspect the basic Ingress example provided in the current directory:

yq basic-ingress.yaml
First, note that the ingressClassName field uses the value cilium. This instructs Kubernetes to use Cilium as the Ingress controller for this resource.

This Kubernetes manifest will create an Ingress resource setting two backends each pointing to a service, based on the requested path:

/details routes to the details service, on port 9080
/ routes to the productpage service, on port 9080
The pathType is Prefix, allowing for subpath matching.

Apply this manifest:

kubectl apply -f basic-ingress.yaml
List services:

kubectl get svc
You will see a LoadBalancer service named cilium-ingress-basic-ingress which was created for the ingress. Cilium LB-IPAM will automatically provision an IP address for it.

The same external IP address is also associated to the Ingress:

kubectl get ingress
Let's retrieve this IP address:

INGRESS_IP=$(kubectl get ingress basic-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $INGRESS_IP

🌐 Make HTTP requests
Check that you can make HTTP requests to that external address:

curl -so /dev/null -w "%{http_code}\n" http://$INGRESS_IP/
The / path takes you to the home page for the bookinfo application, and this should return an HTTP code 200.

From outside the cluster you can also make requests directly to the details service using the path /details (again, you should see a code 200 reply):

curl -so /dev/null -w "%{http_code}\n" http://$INGRESS_IP/details/1
However, you can't directly access other URL paths that weren't defined in basic-ingress.yaml. For example, while you could get JSON data from a request to <address>/details/1, you will get a 404 error if you make a request to <address>/ratings as it is not explicitly mapped in the Ingress resource:

curl -so /dev/null -w "%{http_code}\n" http://$INGRESS_IP/ratings

📚 View the BookInfo application
In order to view the Ingress in our lab, let's assign a static NodePort to it:

kubectl patch svc cilium-ingress-basic-ingress --patch '{"spec": {"type": "LoadBalancer", "ports": [ { "name": "http", "port": 80, "protocol": "TCP", "targetPort": 80, "nodePort": 32100 } ] } }'
You can now click on the 🔗 Bookinfo tab to visit the application. This will generate some requests to the application's microservices.

Next, we'll see how to observe the network flows generated by our application!

## 🔭 Service Mesh Observability
Observability is a central part of a Service Mesh.

It allows to transparently visualize every network flow in your applications without a need to instrumentalize them.

🛰️ Hubble
Hubble is an optional component of Cilium, which brings observability to the network stack.

Hubble provides flows to visualize:

every request going through Cilium
using Cilium identities
providing visibility on L3/L4/L7

hubble observe --namespace default
hubble observe --namespace default -o jsonpb | jq
kubectl annotate pod -l app=productpage --overwrite io.cilium.proxy-visibility="<Ingress/9080/TCP/HTTP>"
kubectl apply -f https://docs.isovalent.com/public/http-ingress-visibility.yaml
hubble observe --protocol http --label app=reviews --port 9080