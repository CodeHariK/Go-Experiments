# [Transparent Encryption with IPSec and WireGuard]()

🔐 Encryption in Cloud Native Applications
Encryption is required for many compliance frameworks.

Kubernetes doesn't natively offer encryption of data in transit.

To offer encryption capabilities, it's often required to implement it directly into your applications or deploy a Service Mesh.

Both options add complexity and operational headaches.

![img](https://cilium.io/045e7a61a5b9c7a8044cc3d0999f0747/wg.gif)

⬢ Encryption in Cilium
Cilium actually provides two options to encrypt traffic between Cilium-managed endpoints: IPsec and WireGuard.

## 🔑 Encrypt Key Management
One of the common challenges with cryptography is the management of keys. Users have to take into consideration aspects such as generation, rotation and distribution of keys.

We'll look at all these aspects in this lab and see the differences between using IPsec and WireGuard as they both have pros and cons. The way it is addressed in Cilium is elegant - the IPsec configuration and associated key are stored as a Kubernetes secret. All secrets are automatically shared across all nodes and therefore all endpoints are aware of the keys.

## 🔑 Generating the Key
First, let's create a Kubernetes secret for the IPsec configuration to be stored.

The format for such IPsec Configuration and key is the following: key-id encryption-algorithms PSK-in-hex-format key-size.

Let's start by generating a random pre-shared key (PSK). We're going to create a random string of 20 characters (using dd with /dev/urandom as a source), then encode it as a hexdump with the xxd command.

Run the following command:

```sh
PSK=($(dd if=/dev/urandom count=20 bs=1 2> /dev/null | xxd -p -c 64))
echo $PSK
```

The $PSK variable now contains our hexdumped PSK.

In order to configure IPsec, you will need to pass this PSK along with a key ID (we'll choose 3 here), and a specification of the algorithm to be used with IPsec (we'll use GCM-128-AES, so we'll specify rfc4106(gcm(aes))). We'll specify the block size accordingly to 128.

As a result, the Kubernetes secret will contain the value 3 rfc4106(gcm(aes)) $PSK 128.

Create a Kubernetes secret called cilium-ipsec-keys, and use this newly created PSK:

```sh
kubectl create -n kube-system secret generic cilium-ipsec-keys \
    --from-literal=keys="3 rfc4106(gcm(aes)) $PSK 128"
```

This command might look confusing at first, but essentially a Kubernetes secret is a key-value pair, with the key being the name of the file to be mounted as a volume in the cilium-agent Pods while the value is the IPsec configuration in the format described earlier.

Decoding the secret created earlier is simple:

```sh
SECRET="$(kubectl get secrets cilium-ipsec-keys -o jsonpath='{.data}' -n kube-system | jq -r ".keys")"
echo $SECRET | base64 --decode
```

Your secret should be similar to this:

3 rfc4106(gcm(aes)) da630c6acdbef2757ab7f5215b8b1811420e3f61 128
This maps to the following Cilium IPsec configuration :

key-id (an identifier of the key): arbitrarily set to 3
encryption-algorithms: AES-GCM GCM
PSK: da630c6acdbef2757ab7f5215b8b1811420e3f61
key-size: 128
Now that the IPSec configuration has been generated, let's install Cilium and IPsec.


⬢ The Cilium CLI
The cilium CLI tool will be used to install and check the status of Cilium in the cluster.

Let's start by installing Cilium on the Kind cluster, with IPsec enabled.

```sh
cilium install --version 1.14.5 \
  --set encryption.enabled=true \
  --set encryption.type=ipsec
```

Wait for the installation to finish:

cilium status --wait

The Cilium status should be OK.

Cilium is now functional on our cluster.

Let's verify that IPsec was enabled by checking that the enable-ipsec key is set to true.

cilium config view | grep enable-ipsec

In the next task, we will verify that traffic has been encrypted, and learn how we can rotate keys.

IPsec encryption was easy to install but we need to verify that traffic has been encrypted.

We will be using the tcpdump packet capture tool for this purpose.

⚙️ Day 2 Operations
Additionally, there will come a point where users will want to rotate keys.

Periodically and automatically rotating keys is a recommended security practice. Cilium currently uses 32-bit keys that can become exhausted depending on the amount of traffic in the cluster. This makes key rotation even more critical.

Some industry standards, such as Payment Card Industry Data Security Standard (PCI DSS), require the regular rotation of keys.

We will see how this can be achieved.

🚀 Deploy a demo app
The endor.yaml manifest will deploy a Star Wars-inspired demo application which consists of:

an endor Namespace, containing
a deathstar Deployment with 1 replicas.
a Kubernetes Service to access the Death Star pods
a tiefighter Deployment with 1 replica
an xwing Deployment with 1 replica
Deploy it:

kubectl apply -f endor.yaml
Verify that the components are properly deployed (execute the command until all pods are running):

kubectl get -f endor.yaml

⏺️ Capture IPsec traffic with tcpdump
Now that applications are deployed in the cluster, let's verify the traffic between the components is encrypted and encapsulated in IPsec tunnels.

First, let's run a shell in one of the Cilium agents:

kubectl -n kube-system exec -ti ds/cilium -- bash
Let's then install the packet analyzer tcpdump to inspect some of the traffic (you may not want to run these in production environments 😅).

apt-get update
apt-get -y install tcpdump
Let's now run tcpdump. We are filtering based on traffic on the cilium_vxlan interface.

When using Kind, Cilium is deployed by default in vxlan tunnel mode - meaning we set VXLAN tunnels between our nodes.

In Cilium's IPsec implementation, we use Encapsulating Security Payload (ESP) as the protocol to provide confidentiality and integrity.

Let's now run tcpdump and filter based on this protocol to show IPsec traffic:

tcpdump -n -i cilium_vxlan esp
Just wait a few seconds and you should see output similar to this:

08:57:55.720756 IP 10.244.3.112 > 10.244.0.74: ESP(spi=0x00000003,seq=0xc5), length 80
08:57:55.721022 IP 10.244.0.74 > 10.244.3.112: ESP(spi=0x00000003,seq=0xc5), length 192
08:57:55.721350 IP 10.244.3.112 > 10.244.0.74: ESP(spi=0x00000003,seq=0xc6), length 164
08:57:55.721530 IP 10.244.0.74 > 10.244.3.112: ESP(spi=0x00000003,seq=0xc6), length 88
08:57:56.978819 IP 10.244.1.155 > 10.244.3.112: ESP(spi=0x00000003,seq=0xc4), length 192
08:57:56.979216 IP 10.244.3.112 > 10.244.1.155: ESP(spi=0x00000003,seq=0xc5), length 164
08:57:56.979390 IP 10.244.1.155 > 10.244.3.112: ESP(spi=0x00000003,seq=0xc5), length 88
08:57:56.981314 IP 10.244.1.155 > 10.244.3.112: ESP(spi=0x00000003,seq=0xc6), length 80
08:57:56.981440 IP 10.244.3.112 > 10.244.1.155: ESP(spi=0x00000003,seq=0xc6), length 80
In the example above, there are three IPs (10.244.1.155, 10.244.3.112, 10.244.0.74); yours are likely to be different). These are the IP addresses of Cilium agents and what we are seeing in the logs is a mesh of IPsec tunnels established between our agents. Notice all these tunnels were automatically provisioned by Cilium.

Every 15 seconds or so, you should see some new traffic, corresponding to the heartbeats between the Cilium agents.

Exit the tcpdump stream with Ctrl+c.


🔑 Key Rotation
As we have seen earlier, the Cilium IPsec configuration and associated key are stored as a Kubernetes secret.

To rotate the key, you will therefore need to patch the previously created cilium-ipsec-keys Kubernetes secret, with kubectl patch secret. During the transition, the new and old keys will be used.

Let's try this now.

Exit the Cilium agent shell (with a prompt similar to root@kind-worker2:/home/cilium#):

exit
You should be back to the green root@server:~# prompt.

Now, let's extract and print some of the variables from our existing secret.

read KEYID ALGO PSK KEYSIZE < <(kubectl get secret -n kube-system cilium-ipsec-keys -o go-template='{{.data.keys | base64decode}}{{printf "\n"}}')
echo $KEYID
echo $PSK
When you run echo $KEYID, it should return 3. We could have guessed this, since we used 3 as the key ID when we initially generated the Kubernetes secret.

Notice the value of the existing PSK by running echo $PSK.

Let's rotate the key. We'll increment the Key ID by 1 and generate a new PSK. We'll use the same key size and encryption algorithm.

NEW_PSK=($(dd if=/dev/urandom count=20 bs=1 2> /dev/null | xxd -p -c 64))
echo $NEW_PSK
patch='{"stringData":{"keys":"'$((($KEYID+1)))' rfc4106(gcm(aes)) '$NEW_PSK' 128"}}'
kubectl patch secret -n kube-system cilium-ipsec-keys -p="${patch}" -v=1
You should see this response: secret/cilium-ipsec-keys patched.

Check the IPsec configuration again:

read NEWKEYID ALGO NEWPSK KEYSIZE < <(kubectl get secret -n kube-system cilium-ipsec-keys -o go-template='{{.data.keys | base64decode}}{{printf "\n"}}')
echo $NEWKEYID
echo $NEWPSK
You can see that the key ID was incremented to 4 and that the PSK value has changed. This example illustrates simple key management with IPsec with Cilium. Production use would probably be more sophisticated.


👍 Well done!
You can now see how you can easily encrypt traffic between pods using Cilium's IPsec implementation.

In the next task, we'll see how we can achieve similar outcomes, using a different technology: WireGuard.

## 🚦 WireGuard
As we saw in the previous task, IPsec encryption provided a great method to achieve confidentiality and integrity.

In addition to IPsec support, Cilium 1.10 also introduced an alternative technology to provide pod-to-pod encryption: WireGuard.

🐉 Wireguard
WireGuard, as described on its official website, is "an extremely simple yet fast and modern VPN that utilizes state-of-the-art cryptography".

Compared to IPsec, "it aims to be faster, simpler, leaner, and more useful, while avoiding the massive headache."

🐉 WireGuard on Cilium
One of the appeals of WireGuard is that it is very opinionated: it leverages very robust cryptography and does not let the user choose ciphers and protocols, like we did for IPsec. It is also very simple to use.

From a Cilium user perspective, the experience is very similar to the IPsec deployment, albeit operationally even simpler. Indeed, the encryption key pair for each node is automatically generated by Cilium and key rotation is performed transparently by the WireGuard kernel module.

Let's get started.


⬢ Installing Cilium with WireGuard
Again, we are using the cilium CLI tool to install Cilium, with WireGuard this time.

Before we start though, we should check that the kernel we are using has support for WireGuard:

uname -ar
WireGuard was integrated into the Linux kernel from 5.6, so our kernel is recent enough to support it. Note that WireGuard was backported to some older Kernels, such as the currently 5.4-based Ubuntu 20.04 LTS.

Cilium was automatically uninstalled before this challenge, so we can go ahead and install Cilium again, this time with WireGuard:

cilium install --version 1.14.5 \
  --set encryption.enabled=true \
  --set encryption.type=wireguard
The installation usually takes a minute or so.

Let's verify the Cilium status.

cilium status --wait
Cilium is now functional on our cluster.


☑️ Validate the setup
You might have noticed that, unlike with IPsec, we didn't have to manually create an encryption key.

One advantage of WireGuard over IPsec is the fact that each node automatically creates its own encryption key-pair and distributes its public key via the network.cilium.io/wg-pub-key annotation in the Kubernetes CiliumNode custom resource object.

Each node's public key is then used by other nodes to decrypt and encrypt traffic from and to Cilium-managed endpoints running on that node.

You can verify this by checking the annotation on the Cilium node kind-worker2, which contains its public key:

kubectl get ciliumnode kind-worker2 \
  -o jsonpath='{.metadata.annotations.network\.cilium\.io/wg-pub-key}'
Let's now run a shell in one of the Cilium agents on the kind-worker2 node.

First, let's get the name of the Cilium agent.

CILIUM_POD=$(kubectl -n kube-system get po -l k8s-app=cilium --field-selector spec.nodeName=kind-worker2 -o name)
echo $CILIUM_POD
Let's now run a shell on the agent.

kubectl -n kube-system exec -ti $CILIUM_POD -- bash
The prompt should be root@kind-worker2:/home/cilium#.

Let's verify that WireGuard was installed:

cilium status | grep Encryption
You should see an output like this one:

Encryption:                            Wireguard       [NodeEncryption: Disabled, cilium_wg0 (Pubkey: qCzNE4dZv6MsMgdk0xFlT8q72c3ZIArvtyFDNlly4gA=, Port: 51871, Peers: 3)]
Let's explain this briefly, going backwards from the last entry:

We have 3 peers: the agent running on each cluster node has established a secure WireGuard tunnel between itself and all other known nodes in the cluster. The WireGuard tunnel interface is named cilium_wg0.
The WireGuard tunnel endpoints are exposed on UDP port 51871.
Notice the public key's value is the same one you previously saw in the node's annotation.
NodeEncryption (the ability to encrypt the traffic between Kubernetes nodes) is disabled. We will enable it on the next task.

🔒 Validate Traffic Encryption
Let's now install the packet analyzer tcpdump to inspect some of the traffic (it may already be on the agent, from the previous task).

apt-get update
apt-get -y install tcpdump
Let's now run tcpdump. Instead of capturing traffic on the VXLAN tunnel interface, we are going to capture traffic on the WireGuard tunnel interface itself, cilium_wg0.

tcpdump -n -i cilium_wg0
Note there should be no output as we've not deployed any Pods yet.

Go to the >_ Terminal 2 tab and deploy a couple of Pods:

kubectl apply -f pod1.yaml -f pod2.yaml -o yaml
We will use these two Pods to run some pings between them and verify that traffic is being encrypted and sent through the WireGuard tunnel.

View the manifests for the two Pods, and notice that we are pinning the pods to different nodes (nodeName: kind-worker and nodeName: kind-worker2) for the purpose of the demo (it's not necessarily a good practice in production).

Verify that both pods are running (launch the command until they are):

kubectl get -f pod1.yaml -f pod2.yaml
Let's get the IP address from our Pod on kind-worker2.

POD2=$(kubectl get pod pod-worker2 --template '{{.status.podIP}}')
echo $POD2
Let's now run a simple ping from the Pod on the kind-worker node:

kubectl exec -ti pod-worker -- ping $POD2
Head back to the >_ Terminal 1 tab. You should see output such as:

18:05:12.321179 IP 10.244.2.8 > 10.244.1.163: ICMP echo request, id 24849, seq 1, length 64
18:05:12.321332 IP 10.244.1.163 > 10.244.2.8: ICMP echo reply, id 24849, seq 1, length 64
18:05:13.321980 IP 10.244.2.8 > 10.244.1.163: ICMP echo request, id 24849, seq 2, length 64
18:05:13.322070 IP 10.244.1.163 > 10.244.2.8: ICMP echo reply, id 24849, seq 2, length 64
Traffic between pods on different nodes has been sent across the WireGuard tunnels and is therefore encrypted.

That's how simple Transparent Encryption is, using WireGuard with Cilium.