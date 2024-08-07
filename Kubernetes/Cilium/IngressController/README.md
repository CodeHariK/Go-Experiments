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

## 🌐 Deploying a gRPC Ingress
While HTTP is still the king of protocols in the web, gRPC is increasingly used, in particular for its low latency and high throughput capabilities.

Let's see how we can deploy a Kubernetes Ingress for a gRPC application using Cilium!

🚀 Deploy a gRPC Application
In this challenge, we will deploy a sample gRPC application, which consists of multiple services such as:

📧 email
🛒 checkout and cart
💡 recommendation
👨‍💻 frontend
💳 payment
🚚 shipping
💱 currency
📦 productcatalog
In this challenge, we will set up a gRPC Ingress with two path prefixes:

/hipstershop.ProductCatalogService pointing to the productcatalog service
/hipstershop.CurrencyService pointing to the currency service

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: grpc-ingress
  namespace: default
spec:
  ingressClassName: cilium
  rules:
    - http:
        paths:
          - backend:
              service:
                name: productcatalogservice
                port:
                  number: 3550
            path: /hipstershop.ProductCatalogService
            pathType: Prefix
          - backend:
              service:
                name: currencyservice
                port:
                  number: 7000
            path: /hipstershop.CurrencyService
            pathType: Prefix


grpcurl -plaintext -proto ./demo.proto $INGRESS_IP:80 hipstershop.CurrencyService/GetSupportedCurrencies
grpcurl -plaintext -proto ./demo.proto $INGRESS_IP:80 hipstershop.ProductCatalogService/ListProducts
```

## 🔑 SNI-based Ingress rules
Due to their frontal position in the cluster, Kubernetes Ingress are commonly used as TLS termination proxies.

In this challenge, we will use self-signed certificates for simplicity. In production, it is recommended to either inject known CA credentials, or generate dynamic certificates using e.g. Cert Manager.

## 🔑 Create TLS certificate and private key
For demonstration purposes we will use a TLS certificate signed by a made-up, self-signed certificate authority (CA). One easy way to do this is with mkcert.

Create a certificate that will validate bookinfo.cilium.rocks and hipstershop.cilium.rocks, as these are the host names used in this ingress example:

```sh
mkcert '*.cilium.rocks'
Mkcert created a key (_wildcard.cilium.rocks-key.pem) and a certificate (_wildcard.cilium.rocks.pem) that we will use for the ingress service.

Create a Kubernetes secret with this key and certificate:

kubectl create secret tls demo-cert \
  --key=_wildcard.cilium.rocks-key.pem \
  --cert=_wildcard.cilium.rocks.pem
```

🚪 Deploy the ingress
Delete the ingresses created in the last two challenges:

kubectl delete ingress basic-ingress
kubectl delete ingress grpc-ingress
The ingress configuration for this demo provides the same routing as those demos but with the addition of TLS termination:

the /hipstershop.CurrencyService prefix will be routed to the currency gRPC service deployed in the gRPC challenge
the /details prefix will be routed to the details HTTP service deployed in the HTTP challenge
the / prefix will be routed to the productpage HTTP service deployed in the HTTP challenge
These three services will be secured via TLS and accessible on two domain names:

bookinfo.cilium.rocks
hipstershop.cilium.rocks
Inspect the ingress to verify these rules:

yq tls-ingress.yaml

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-ingress
  namespace: default
spec:
  ingressClassName: cilium
  rules:
    - host: hipstershop.cilium.rocks
      http:
        paths:
          - backend:
              service:
                name: productcatalogservice
                port:
                  number: 3550
            path: /hipstershop.ProductCatalogService
            pathType: Prefix
          - backend:
              service:
                name: currencyservice
                port:
                  number: 7000
            path: /hipstershop.CurrencyService
            pathType: Prefix
    - host: bookinfo.cilium.rocks
      http:
        paths:
          - backend:
              service:
                name: details
                port:
                  number: 9080
            path: /details
            pathType: Prefix
          - backend:
              service:
                name: productpage
                port:
                  number: 9080
            path: /
            pathType: Prefix
  tls:
    - hosts:
        - bookinfo.cilium.rocks
        - hipstershop.cilium.rocks
      secretName: demo-cert
```

Then deploy the ingress to the cluster:

kubectl apply -f tls-ingress.yaml
This creates a LoadBalancer service, which after around 30 seconds or so should be populated with an external IP address.

Verify that the Ingress has an load balancer IP address assigned:

kubectl get ingress tls-ingress
Then assign this IP to the INGRESS_IP variable so we can make use of it:

INGRESS_IP=$(kubectl get ingress tls-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $INGRESS_IP

🗺️ Edit /etc/hosts
In this ingress configuration, the host names hipstershop.cilium.rocks and bookinfo.cilium.rocks are specified in the path routing rules.

Since we do not have DNS entries for these names, we will modify the /etc/hosts file on the host to manually associate these names to the known ingress IP we retrieved. Execute the following in the >_ Terminal tab:

```sh
cat << EOF >> /etc/hosts
${INGRESS_IP} bookinfo.cilium.rocks
${INGRESS_IP} hipstershop.cilium.rocks
EOF
Requests to these names will now be directed to the ingress.
```

🌐 Make requests
Install the Mkcert CA into your system so cURL can trust it:

```sh
mkcert -install
Now let's make a request to the ingress:

curl -s https://bookinfo.cilium.rocks/details/1 | jq
The data should be properly retrieved, using HTTPS (and thus, the TLS handshake was properly achieved).

Similarly you can test a gRPC request:

grpcurl -proto ./demo.proto hipstershop.cilium.rocks:443 hipstershop.ProductCatalogService/ListProducts | jq
```

We can see that our single Ingress resource now allows to access both applications, in a secure manner over HTTPS, using a valid TLS certificate.

Let's validate our learning by taking a short quiz and a lab!