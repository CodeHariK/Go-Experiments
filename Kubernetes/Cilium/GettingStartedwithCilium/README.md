# [Getting Started with Cilium](https://isovalent.com/labs/cilium-getting-started/)

👮🏻 Ensuring Security on Kubernetes
In a Cloud Native/Kubernetes environment, how do you

enforce policies?
troubleshoot the network?
stay secure with minimal effort?
After all, in the highly dynamic and complex world of microservices, IP addresses and ports are no longer relevant.

A new approach is needed — meet Cilium!

Cilium provides Connectivity, Observability and Security capabilities in a Cloud Native World, and is based on eBPF.

![cilim](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/c602a2a595a6cc09ad58a6211dd93d9b/assets/cilium_overview.png)

Thanks to eBPF, a revolutionary new kernel extensibility mechanism of Linux, we have the opportunity to rethink the Linux networking and security stack for the age of microservices.

## 🪪 Identities, Protocol parsing & Observability
From inception, Cilium was designed for large-scale, highly-dynamic containerized environments.

Cilium:

natively understands container identities
parses API protocols like HTTP, gRPC, and Kafka
and provides visibility and security that is both simpler and more powerful than traditional approaches

![cilium](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/c8aefb5236f2b02b7c3d775ae88d8539/assets/identity_store.png)

![cilium](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/a782fe149f90c8a88716a69170e039f1/assets/cilium-arch.png)

## 🏛 The Kind Cluster
Let's have a look at this lab's environment.

We are running a Kind Kubernetes cluster, and on top of that Cilium.

While the Kind cluster is finishing to start, let's have a look at its configuration:

cat /etc/kind/${KIND_CONFIG}.yaml

🖥 Nodes
In the nodes section, you can see that the cluster consists of three nodes:

1 control-plane node running the Kubernetes control plane and etcd
2 worker nodes to deploy the applications

🔀 Networking
In the networking section of the configuration file, the default CNI has been disabled so the cluster won't have any Pod network when it starts. Instead, Cilium is being deployed to the cluster to provide this functionality.

To see if the Kind cluster is ready, verify that the cluster is properly running by listing its nodes:

kubectl get nodes
You should see the three nodes appear, all marked as NotReady. This is normal, since the CNI is disabled, and we will install Cilium in the next step. If you don't see all nodes, the workers nodes might still be joining the cluster. Relaunch the command until you can see all three nodes listed.

Now that we have a Kind cluster, let's install Cilium on it!

```sh
root@server:~# cat /etc/kind/${KIND_CONFIG}.yaml
---
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  # localhost.run proxy
  - containerPort: 32042
    hostPort: 32042
  # Hubble relay
  - containerPort: 31234
    hostPort: 31234
  # Hubble UI
  - containerPort: 31235
    hostPort: 31235
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
- role: worker
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
- role: worker
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
networking:
  disableDefaultCNI: true

root@server:~# kubectl get nodes
NAME                 STATUS     ROLES           AGE     VERSION
kind-control-plane   NotReady   control-plane   8m53s   v1.27.3
kind-worker          NotReady   <none>          8m32s   v1.27.3
kind-worker2         NotReady   <none>          8m32s   v1.27.3
```

## 🖥️ The Cilium CLI
The cilium CLI tool can install and update Cilium on a cluster, as well as activate features —such as Hubble and Cluster Mesh.

```sh
cilium install
cilium status --wait
```

## 🛸 Star Wars Demo
To learn how to use and enforce policies with Cilium, we have prepared a demo example.

In the following Star Wars-inspired example, there are three microservice applications: deathstar, tiefighter, and xwing.

🌐 The deathstar service
The deathstar runs an HTTP webservice on port 80, which is exposed as a Kubernetes Service to load-balance requests to deathstar across two pod replicas.

The deathstar service provides landing services to the empire’s spaceships so that they can request a landing port.

👮🏻 Allowing ship access
The tiefighter pod represents a landing-request client service on a typical empire ship and xwing represents a similar service on an alliance ship.

With this setup, we can test different security policies for access control to deathstar landing services.

![cilium](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/5c8d0309fb62c5a9f574bcdf3ffd11e0/assets/star_wars.png)

Let's deploy a simple empire demo application. It is made of several microservices, each identified by Kubernetes labels:

the Death Star: org=empire, class=deathstar
the Imperial TIE fighter: org=empire, class=tiefighter
the Rebel X-Wing: org=alliance, class=xwing

kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/http-sw-app.yaml

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: deathstar
  labels:
    app.kubernetes.io/name: deathstar
spec:
  type: ClusterIP
  ports:
  - port: 80
  selector:
    org: empire
    class: deathstar
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deathstar
  labels:
    app.kubernetes.io/name: deathstar
spec:
  replicas: 2
  selector:
    matchLabels:
      org: empire
      class: deathstar
  template:
    metadata:
      labels:
        org: empire
        class: deathstar
        app.kubernetes.io/name: deathstar
    spec:
      containers:
      - name: deathstar
        image: docker.io/cilium/starwars
---
apiVersion: v1
kind: Pod
metadata:
  name: tiefighter
  labels:
    org: empire
    class: tiefighter
    app.kubernetes.io/name: tiefighter
spec:
  containers:
  - name: spaceship
    image: docker.io/cilium/json-mock
---
apiVersion: v1
kind: Pod
metadata:
  name: xwing
  labels:
    app.kubernetes.io/name: xwing
    org: alliance
    class: xwing
spec:
  containers:
  - name: spaceship
    image: docker.io/cilium/json-mock
```

Each pod will also be represented in Cilium as an Endpoint. To retrieve a list of all endpoints managed by Cilium, the Cilium Endpoint (or cep) resource can be used:

```sh
root@server:~# kubectl get pods,svc
NAME                             READY   STATUS    RESTARTS   AGE
pod/deathstar-7848d6c4d5-j6vnm   1/1     Running   0          42s
pod/deathstar-7848d6c4d5-pj8l2   1/1     Running   0          42s
pod/tiefighter                   1/1     Running   0          42s
pod/xwing                        1/1     Running   0          42s

NAME                 TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/deathstar    ClusterIP   10.96.81.109   <none>        80/TCP    42s
service/kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP   20m


root@server:~# kubectl get cep --all-namespaces
NAMESPACE            NAME                                      ENDPOINT ID   IDENTITY ID   INGRESS ENFORCEMENT   EGRESS ENFORCEMENT   VISIBILITY POLICY   ENDPOINT STATE   IPV4           IPV6
default              deathstar-7848d6c4d5-j6vnm                1746          49585         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.49    
default              deathstar-7848d6c4d5-pj8l2                1052          49585         <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.239   
default              tiefighter                                3166          14860         <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.4     
default              xwing                                     854           3900          <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.118   
kube-system          coredns-5d78c9869d-5f8pz                  180           32465         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.152   
kube-system          coredns-5d78c9869d-xxg5d                  1596          32465         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.26    
local-path-storage   local-path-provisioner-6bc4bddd6b-hdhjk   2746          14054         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.234   


root@server:~# kubectl api-resources | grep cilium
ciliumcidrgroups                   ccg                                 cilium.io/v2alpha1                     false        CiliumCIDRGroup
ciliumclusterwidenetworkpolicies   ccnp                                cilium.io/v2                           false        CiliumClusterwideNetworkPolicy
ciliumendpoints                    cep,ciliumep                        cilium.io/v2                           true         CiliumEndpoint
ciliumexternalworkloads            cew                                 cilium.io/v2                           false        CiliumExternalWorkload
ciliumidentities                   ciliumid                            cilium.io/v2                           false        CiliumIdentity
ciliuml2announcementpolicies       l2announcement                      cilium.io/v2alpha1                     false        CiliumL2AnnouncementPolicy
ciliumloadbalancerippools          ippools,ippool,lbippool,lbippools   cilium.io/v2alpha1                     false        CiliumLoadBalancerIPPool
ciliumnetworkpolicies              cnp,ciliumnp                        cilium.io/v2                           true         CiliumNetworkPolicy
ciliumnodeconfigs                                                      cilium.io/v2alpha1                     true         CiliumNodeConfig
ciliumnodes                        cn,ciliumn                          cilium.io/v2                           false        CiliumNode
ciliumpodippools                   cpip                                cilium.io/v2alpha1                     false        CiliumPodIPPool
```

## 🌐 Death Star access
From the perspective of the deathstar service, only the ships with label org=empire are allowed to connect and request landing!

🛡️ No rules enforced
But we have no rules enforced, so what will happen if not only tiefighter but also xwing request landing?

Let's find out!

🔍 Check Current Access
To simulate our connectivity tests, we will be executing simple API calls using curl.

Let's test if we can land our TIE fighter on the Death Star by running the following command:

```sh
kubectl exec tiefighter -- \
  curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

The command above lets us get a shell on the tiefighter pod and run a HTTP POST request to the deathstar Service to request landing.

The command should work —as the TIE fighter and the Death Star are on the same side of the galactic wars (i.e. the bad guys).

Now test if you can land your X-wing (i.e. the good guys) with:

```sh
kubectl exec xwing -- \
  curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

So far, it seems access is allowed! This is good for the rebel alliance —unfettered access to the DeathStar— but this should not be allowed, right?

There is a security policy missing!

## 🪪 Identities and Cloud Native

IP addresses are no longer relevant for Cloud Native workloads. Security policies need something else.

Cilium provides this: Cilium uses the labels assigned to pods to define security policies.

## ✍️ Writing Network Policies

We’ll start with a basic policy restricting deathstar landing requests to only the ships that have the label org=empire.

This blocks any ships without the org=empire label to even connect to the deathstar service.

This is a simple policy that filters only on network layer 3 (IP protocol) and network layer 4 (TCP protocol), so it is often referred to as a L3/L4 network security policy.

![policy](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/6ff094083acf6bf5d3e4e12a511628af/assets/star_wars_l3l4.png)

```yaml
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rule1"
spec:
  description: "L3-L4 policy to restrict deathstar access to empire ships only"
  endpointSelector:
    matchLabels:
      org: empire
      class: deathstar
  ingress:
  - fromEndpoints:
    - matchLabels:
        org: empire
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
```

🛡️ Enforcing the Network Policy
Once you are done visualizing the policy in the editor, change back to the >_ Terminal tab. There we can apply a preconfigured network policy with the values discussed above to our demo system:

```sh
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/sw_l3_l4_policy.yaml
Now let's try to land the empire tiefighter again (HTTP POST from tiefighter to deathstar on the /v1/request-landing path):

kubectl exec tiefighter -- \
  curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
This still works, which is expected.

In comparison, if you try to request landing from the xwing pod, you will see that the request will eventually time out:

kubectl exec xwing -- \
  curl -s -XPOST deathstar.default.svc.cluster.local/v1/request-landing
```

Kill the request with Ctrl+C once you realize that it hangs.

We have successfully blocked access to the deathstar from an X-Wing ship. Let's now see how we could make this policy a bit more fine-grained using L7 rules.

## 🛡️ Tighter Rules

So far it was sufficient to either give tiefighter/xwing full access to deathstar’s API or no access at all. But are you absolutely sure that you can trust all thousands of tiefighter pilots of the entire empire?

We must provide the strongest security (i.e., enforce least-privilege isolation) between microservices: each service that calls deathstar’s API should be limited to making only the set of HTTP requests it requires for legitimate operation.

🌐 Filtering on HTTP
So we need to filter on a higher level: we need to filter the actual HTTP requests!

![rules](https://play.instruqt.com/assets/tracks/ucdsyxm1sfzh/5d08c712a544067c7d4904e6d59f9346/assets/star_wars_l7.png)

📻 A Radio Contact
The Death Star is now well secured to only allow Imperial vessels to access it.

What if the Rebellion was able to take control of an Imperial Tie Figther though? You —a Rebel officer- have just taken control of a Tie Fighter and are approaching the Death Star.

🛣️ Filtering Paths
Consider that the deathstar service exposes some maintenance APIs which should not be called by random empire ships. To see why those APIs are sensitive, run:

```sh
kubectl exec tiefighter -- \
  curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
```

Yes, there is a Panic: the Death Star just exploded!

As you can see, this leads to rather unwanted results. While this is an illustrative example, unauthorized access such as above can have adverse security repercussions. We need to enforce policies on the HTTP layer, at layer 7, to limit what exact APIs the tiefighter is allowed to call —and which are not.

We need to extend the existing policy with an HTTP rule such as:

```
rules:
  http:
  - method: "POST"
    path: "/v1/request-landing"
```

This will restrict API access to only the /v1/request-landing path and will thus prevent users from accessing the /v1/exhaust-port, which caused a crash as we saw earlier.

```yaml
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rule1"
spec:
  description: "L7 policy to restrict access to specific HTTP call"
  endpointSelector:
    matchLabels:
      org: empire
      class: deathstar
  ingress:
  - fromEndpoints:
    - matchLabels:
        org: empire
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: "POST"
          path: "/v1/request-landing"
```

Run the same test as above, and see the different outcome:

```sh
kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
```

As you can see, with Cilium L7 security policies, we are able to restrict tiefighter's access to only the required API resources on deathstar, thereby implementing a “least privilege” security approach for communication between microservices.

