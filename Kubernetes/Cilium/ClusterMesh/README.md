# [Cluster Mesh](https://isovalent.com/labs/cilium-cluster-mesh/)

In this track, we will set up two Kubernetes clusters and connect them using Cilium Cluster Mesh.

We will then deploy applications and services on both clusters, and see how we can make these applications global and resilient across clusters.

After the services are deployed and made global, we will see how this can be used to improve availability and fault tolerance of multi-cluster applications.

![cilium](https://cilium.io/static/04d2d06e7e32665b74c968a9f7fc0a40/905a7/usecase_ha.webp)


[![Watch the video](https://img.youtube.com/vi/1fsXtqg4Pkw/maxresdefault.jpg)](https://youtu.be/1fsXtqg4Pkw)

We'll have two requirements for these clusters:

disable default CNI so we can easily install Cilium
use disjoint pods and services subnets


1️⃣ Koornacht Cluster
Go to the >_ 1️⃣ Koornacht cluster tab.

Let's have a look at the configuration for the first cluster, which we will be calling Koornacht:

# ⚠️ In the Koornacht tab
yq kind_koornacht.yaml
This cluster will feature one control-plane node and 2 worker nodes, and use 10.1.0.0/16 for the Pod network, and 172.20.1.0/24 for the Services.

Create the Koornacht first cluster with:

---
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
networking:
  disableDefaultCNI: true
  podSubnet: 10.1.0.0/16
  serviceSubnet: 172.20.1.0/24
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
  - role: worker
  - role: worker

# ⚠️ In the Koornacht tab
kind create cluster --name koornacht --config kind_koornacht.yaml
This usually takes about 1 minute.

Verify that all 3 nodes are up:

# ⚠️ In the Koornacht tab
kubectl get nodes
The nodes are marked as NotReady because there is not CNI plugin set up yet.

Then install Cilium on it:

# ⚠️ In the Koornacht tab
cilium install \
  --set cluster.name=koornacht \
  --set cluster.id=1 \
  --set ipam.mode=kubernetes
Let's also enable Hubble for observability, only on the Koornacht cluster:

# ⚠️ In the Koornacht tab
cilium hubble enable --ui
Verify that everything is fine with:

# ⚠️ In the Koornacht tab
cilium status --wait

2️⃣ Tion Cluster
Let's now create a second Kind cluster, which we will call Tion.

Switch to the >_ 2️⃣ Tion tab, and inspect the configuration:

# ⚠️ In the Tion tab
yq kind_tion.yaml
This Tion cluster will also feature one control-plane node and 2 worker nodes, but it will use 10.2.0.0/16 for the Pod network, and 172.20.2.0/24 for the Services.

Create the Tion cluster with:

yq kind_tion.yaml
---
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
networking:
  disableDefaultCNI: true
  podSubnet: 10.2.0.0/16
  serviceSubnet: 172.20.2.0/24
nodes:
  - role: control-plane
  - role: worker
  - role: worker

# ⚠️ In the Tion tab
kind create cluster --name tion --config kind_tion.yaml
Verify that all 3 nodes are up:

# ⚠️ In the Tion tab
kubectl get nodes
Then install Cilium on it:

# ⚠️ In the Tion tab
cilium install \
  --set cluster.name=tion \
  --set cluster.id=2 \
  --set ipam.mode=kubernetes
Verify that everything is fine with:

# ⚠️ In the Tion tab
cilium status --wait
Now that we have two Kind clusters installed with Cilium, let's get them meshed!

🏛️ Cluster Mesh Architecture
When activating Cluster Mesh on Cilium clusters, a new Control Plane is deployed to manage the mesh for this cluster, along with its etcd key-value store.

Agents of other clusters can then access this Cluster Mesh Control Plane in read-only mode, allowing them to access metadata about the cluster, such as Service names and corresponding IPs.

![Architecture](https://play.instruqt.com/assets/tracks/gcdcyyilhrqt/0f3b2c9fcfd6b679ee94210a7c7589d2/assets/architecture.png)

Cilium Cluster Mesh allows to link multiple Kubernetes clusters, provided:

all clusters run Cilium as CNI
all worker nodes have a unique IP address and are able to connect to each other

## Enable Cluster Mesh on both clusters:

# ⚠️ In *both* the Koornacht and Tion tabs
cilium clustermesh enable --service-type NodePort
Note
Several types of connectivity can be set up. We're using NodePort in our case as it's easier and we don't have dynamic load balancers available.

For production clusters, it is strongly recommended to use LoadBalancer instead.

Wait for Cluster Mesh to be ready on both clusters:

# ⚠️ In *both* the Koornacht and Tion tabs
cilium clustermesh status --wait
You can also verify the Cluster Mesh status with cilium status:

# ⚠️ In *both* the Koornacht and Tion tabs
cilium status
You should see a ClusterMesh: OK field.


🤝 Mesh Clusters
Let's now connect the clusters by instructing one cluster to mesh with the second one. This needs to be done in a shell with access to both cluster contexts, so we'll use the >_ 🌐 Global tab for that:

# ⚠️ In the Global tab
cilium clustermesh connect \
  --context kind-koornacht \
  --destination-context kind-tion
Wait for both clusters to be ready:

# ⚠️ In *both* the Koornacht and Tion tabs
cilium clustermesh status --wait
Our two clusters are now meshed together. Let's deploy applications on them!

## 🌌 Deploying an application
We will now deploy a sample application on both Kubernetes clusters.

This application will contain two deployments:

a simple HTTP application called rebel-base, which will return a static JSON document
an x-wing pod which we will use to make requests to the rebel-base service from within the cluster
The only difference between the two deployments will be the ConfigMap resource deployed, which will contain the static JSON document served by rebel-base, and whose content will depend on the cluster.

Are you ready? Let's go!

Let's prepare to deploy on the >_ 1️⃣ Koornacht Cluster.

We will deploy a simple HTTP application returning a JSON, including the name of the cluster:

# ⚠️ In the Koornacht tab
kubectl apply -f deployment.yaml
The ConfigMap for this service contains the JSON reply, with the name of the Cluster hardcoded (-o yaml is added here to show you the content of the resource):

# ⚠️ In the Koornacht tab
kubectl apply -f configmap_koornacht.yaml -o yaml
Check that the pods are running properly (launch until all 4 pods are Running):

# ⚠️ In the Koornacht tab
kubectl get pod
You should see something like:

NAME                          READY   STATUS    RESTARTS   AGE
rebel-base-6985d8f76f-n6qmm   1/1     Running   0          44s
rebel-base-6985d8f76f-rn4ht   1/1     Running   0          44s
x-wing-6d58648f95-2mrpc       1/1     Running   0          40s
x-wing-6d58648f95-nw927       1/1     Running   0          40s
Apply the Service for the application:

# ⚠️ In the Koornacht tab
kubectl apply -f service.yaml
Let's test this service, using the x-wing pod deployed alongside the rebel-base deployment:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl rebel-base; done'
You should see 10 lines of log, all containing:

{"Cluster": "Koornacht", "Planet": "N'Zoth"}
Go the 🔗 🛰️ Hubble UI tab and select the default namespace. You will see the x-wing requests being sent to the rebel-base pods.


2️⃣ Tion Cluster
We will deploy the same application and service on the >_ 2️⃣ Tion cluster, with a small difference: the JSON answer will reply with Tion since we're using a slightly different ConfigMap:

# ⚠️ In the Tion tab
kubectl apply -f deployment.yaml
kubectl apply -f configmap_tion.yaml -o yaml
kubectl apply -f service.yaml
Wait until the pods are ready (run kubectl get po until all pods are Ready) and check this second service:

# ⚠️ In the Tion tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl rebel-base; done'
After the pods start, you should see 10 lines of log, all containing:

{"Cluster": "Tion", "Planet": "Foran Tutha"}
We now have similar applications running on both our clusters. Wouldn't it be great if we could load-balance traffic between them? This is precisely what we'll be doing in the next challenge!

## 🌍 Making Services Global
When two or more clusters are meshed, Cilium allows you to set services as global in one or more clusters, by adding an annotation to them:

service.cilium.io/global: "true"
When this annotation is set, requests to this service will load-balance to all available services with the same name and namespace in all meshed clusters.

💥 Fault Resilience
One obvious usage of global services is fault tolerance.

If a service becomes unavailable in one cluster, traffic can be redirected to the same service in other clusters, ensuring a continuity of service.

![img](https://cilium.io/static/04d2d06e7e32665b74c968a9f7fc0a40/905a7/usecase_ha.webp)

# ⚠️ In *both* the Koornacht and Tion tabs
kubectl annotate service rebel-base service.cilium.io/global="true"
When accessing the service from either cluster, it should now be load-balanced between the two clusters, because it is marked as global:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl rebel-base; done'

## 🌎 Global vs Shared

We've seen how the service.cilium.io/global annotation allows for a cluster to load-balance requests to a service to all meshed clusters with the same annotated service.

What if you want to remove the service of a specific cluster from the global service?

The service.cilium.io/shared annotation can be used for this.

🚫 Disabling Global Service
By default, a service marked as global is considered shared as well, so the value of service.cilium.io/shared is true for all clusters where the service is marked as global.

Setting it to false in a cluster removes that specific service from the global service:

service.cilium.io/shared: "false"

By default, a service marked as global is considered shared as well, so the value of service.cilium.io/shared is true for all clusters where the service is marked as global.

Setting it to false in a cluster removes that specific service from the global service:

## ⏳ Global services & latency
Global services allow to load-balance traffic across multiple clusters.

As we have seen, this is very useful to implement a fallback policy for redundant services.

Most of the time however, it would be useful to limit latency by only using remote services when local ones are not available.

This is the objective of service affinity.

kubectl annotate service rebel-base service.cilium.io/affinity=local

The opposite effect can be obtained by using remote as the annotation value.

## 👮 Securing Cross-Cluster Communication
In the previous examples, we have used a pod (x-wing) as a curl client to access another set of pods (rebel-base) either in a local or a remote cluster.

Cilium Network Policies are Kubernetes resources which allow to restrict access between pods, by labels.

🛡️ Cross-Cluster Network Policies
When using Cilium Cluster Mesh, it is possible to add Cilium Network Policies to filter traffic between clusters.

At the moment, the Koornacht service load-balances to both Koornacht and Tion clusters, with a local affinity.

In the context of a Zero-Trust security policy, we would like to block all traffic, and then allow only what is necessary.

## 🗑 Remove affinity
For this challenge, let's start by removing the local affinity we placed on the >_ 1️⃣ Koornacht service earlier:

# ⚠️ In the Koornacht tab
kubectl annotate service rebel-base service.cilium.io/affinity-
Check that the service balances again to both clusters:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl --max-time 2 rebel-base; done'

❌ Default Deny
By default, all communication is allowed between the pods. In order to implement Network Policies, we thus need to start with a default deny rule, which will disallow communication. We will then add specific rules to add the traffic we want to allow.

Adding a default deny rule is achieved by selecting all pods (using {} as the value for the endpointSelector field) and using empty rules for ingress and egress fields.

However, blocking all egress traffic would prevent nodes from performing DNS requests to Kube DNS, which is something we want to avoid. For this reason, our default deny policy will include an egress rule to allow access to Kube DNS on UDP/53, so all pods are able to resolve service names:

---
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "default-deny"
spec:
  description: "Default Deny"
  endpointSelector: {}
  ingress:
    - {}
  egress:
    - toEndpoints:
        - matchLabels:
            io.kubernetes.pod.namespace: kube-system
            k8s-app: kube-dns
      toPorts:
        - ports:
            - port: "53"
              protocol: UDP
          rules:
            dns:
              - matchPattern: "*"
Copy this Kubernetes manifest and paste it to the default-deny.yaml using the </> Editor tab.

Then apply the manifest to both clusters:

# ⚠️ In *both* the Koornacht and Tion tabs
kubectl apply -f default-deny.yaml
Now test the requests again:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl --max-time 2 rebel-base; done'
As expected from the application of the default deny policy, all requests now time out.


🛰️ Visualizing with Hubble
We installed Hubble, Cilium's observability component, on the Koornacht cluster.

You can use its CLI to visualize packet drops:

# ⚠️ In the Koornacht tab
hubble observe --verdict DROPPED --from-pod default/x-wing
You can see an x-wing pod trying to reach out to rebel-base pods:

Dec  4 14:00:20.568: default/x-wing-64665f7b7b-w658f:34646 (ID:105931) <> default/rebel-base-6cf4b8c8b-6fjn5:80 (ID:151579) policy-verdict:none EGRESS DENIED (TCP Flags: SYN)
Dec  4 14:00:20.568: default/x-wing-64665f7b7b-w658f:34646 (ID:105931) <> default/rebel-base-6cf4b8c8b-6fjn5:80 (ID:151579) Policy denied DROPPED (TCP Flags: SYN)
On each of these lines, the default/x-wing-64665f7b7b-w658f client pod is trying to reach the default/rebel-base-6cf4b8c8b-6fjn5 pod on port TCP/80, sending a SYN TCP flag. These packets are dropped because of the default deny policy, and the client pod never receives a SYN-ACK TCP reply.

You can also verify this in the Hubble UI. Go to the 🔗 🛰️ Hubble UItab and click on the x-wing boxfor the Koornacht cluster.

The logs at the bottom of the screen list the same dropped verdict flows you just listed with the Hubble CLI.

Egress denied in Hubble UI

Note
You might see a more complex diagram with dotted-red lines instead. This is because the flows from before applying the Network Policy are still in the buffer. Wait a few minutes and refresh the Hubble UI interface.


✅ Allowing Cross-Cluster traffic
We want to allow the Koornacht x-wing pods to access the rebel-base pods on both the local and Tion clusters. Since all traffic is now denied by default, we need to add a new Network Policy to allow this specific traffic.

This CiliumNetworkPolicy resource allows the pods with a name=x-wing label located in the koornacht cluster to reach out to any pod with a name=rebel-base label.

---
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "x-wing-to-rebel-base"
spec:
  description: "Allow x-wing in Koornacht to contact rebel-base"
  endpointSelector:
    matchLabels:
      name: x-wing
      io.cilium.k8s.policy.cluster: koornacht
  egress:
  - toEndpoints:
    - matchLabels:
        name: rebel-base
Using the </> Editor tab, save this manifest to x-wing-to-rebel-base.yaml, then apply it in the >_ 1️⃣ Koornacht tab:

# ⚠️ In the Koornacht tab
kubectl apply -f x-wing-to-rebel-base.yaml
Try the request again:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl --max-time 2 rebel-base; done'
The requests are still dropped. Our default deny policy blocks both ingress and egress connections for all pods, but the new policy we've added only allows egress connectivity. We also need to allow ingress connections to reach the rebel-base pods. Let's fix this with a new CiliumNetworkPolicy resource:

---
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rebel-base-from-x-wing"
spec:
  description: "Allow rebel-base to be contacted by Koornacht's x-wing"
  endpointSelector:
    matchLabels:
      name: rebel-base
  ingress:
  - fromEndpoints:
    - matchLabels:
        name: x-wing
        io.cilium.k8s.policy.cluster: koornacht
Using the </> Editor tab, save this manifest to rebel-base-from-x-wing.yaml, then apply it in the >_ 1️⃣ Koornacht tab:

# ⚠️ In the Koornacht tab
kubectl apply -f rebel-base-from-x-wing.yaml
Now test the service again:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl --max-time 2 rebel-base; done'
It works, but only partially, as only the requests to the Koornacht cluster go through:

curl: (28) Connection timed out after 2000 milliseconds
curl: (28) Connection timed out after 2000 milliseconds
curl: (28) Connection timed out after 2000 milliseconds
{"Cluster": "Koornacht", "Planet": "N'Zoth"}
curl: (28) Connection timed out after 2000 milliseconds
curl: (28) Connection timed out after 2000 milliseconds
{"Cluster": "Koornacht", "Planet": "N'Zoth"}
{"Cluster": "Koornacht", "Planet": "N'Zoth"}
{"Cluster": "Koornacht", "Planet": "N'Zoth"}
curl: (28) Connection timed out after 2001 milliseconds
command terminated with exit code 28
This is because we haven't applied any specific policies to the Tion cluster, where the default deny policy was also deployed.

We need to apply the rebel-base-from-x-wing Network Policy to the >_ 2️⃣ Tion cluster to allow the ingress connection:

# ⚠️ In the Tion tab
kubectl apply -f rebel-base-from-x-wing.yaml
Test once more in the >_ 1️⃣ Koornacht tab:

# ⚠️ In the Koornacht tab
kubectl exec -ti deployments/x-wing -- /bin/sh -c 'for i in $(seq 1 10); do curl --max-time 2 rebel-base; done'
The requests all go through, and we have successfully secured our service across clusters!