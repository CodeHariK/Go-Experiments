# [Cilium Envoy L7 Proxy](https://isovalent.com/labs/cilium-envoy-l7-proxy/)

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

```sh
kubectl get cnp rule1 -o yaml | yq .spec

kubectl exec tiefighter -- \
  curl -s -X PUT deathstar.default.svc.cluster.local/v1/exhaust-port
  
hubble observe \
  --from-pod default/tiefighter \
  --to-pod default/deathstar
```

## 🛠️ Embedded Envoy
In >_ Terminal 1, make the Tie Fighter request landing in a loop:

while [ 1 ]; do
  kubectl exec tiefighter -- \
    curl -s --max-time 1 -X POST deathstar.default.svc.cluster.local/v1/request-landing
  sleep 1
done
In >_ Terminal 2, enable Cilium debug:

cilium config set debug true
This will make the Cilium pods restart on the nodes.

Wait for Cilium to be ready again:

cilium status --wait
List the pods in the kube-system namespace:

kubectl -n kube-system get pods
There are no pods specific to Envoy.

Inspect a Cilium pod on a node where the Death Star is deployed:

NODE=$(kubectl get po -l class=deathstar -o jsonpath='{.items[].spec.nodeName}')
echo $NODE
CILIUM=$(kubectl -n kube-system get po -l k8s-app=cilium --field-selector spec.nodeName=$NODE -o name)
echo $CILIUM
kubectl -n kube-system exec $CILIUM -c cilium-agent -- \
  ps axu | grep envoy
You can see that a process called cilium-envoy-starter is running in order to enforce the L7 network policy.

Go back to >_ Terminal 1 and observe that some requests have failed. This is because when the Cilium pods restarted, it also restarted the Envoy proxy on the node since it is embedded in the pod.


🧩 Separate DaemonSet
In >_ Terminal 2, upgrade Cilium with the envoy.enabled=true Helm option to use a separate DaemonSet for Envoy:

helm repo add cilium https://helm.cilium.io
helm -n kube-system upgrade cilium cilium/cilium \
  --reuse-values \
  --set envoy.enabled=true
Wait for Cilium to be ready:

cilium status --wait
The status should indicate "Envoy DaemonSet: OK".

There might still be a few packets lost in >_ Terminal 1 since the Cilium pods restarted again.

In >_ Terminal 2, check a Cilium pod again:

CILIUM=$(kubectl -n kube-system get po -l k8s-app=cilium --field-selector spec.nodeName=$NODE -o name)
kubectl -n kube-system exec $CILIUM -c cilium-agent -- \
  ps axu | grep envoy
There's no Envoy process running in the Cilium pod anymore!

Instead, there's now a dedicated DaemonSet which deployed one Envoy pod per node:

kubectl -n kube-system get po \
  -l k8s-app=cilium-envoy
This allows Envoy workers to be managed separately from the Cilium agents. In particular, it means operators can now tune resources for both Cilium and Envoy separately, as the pod privileges are also better tuned.

Besides that, the Envoy logs are now accessible in the Envoy pods directly instead of the Cilium pods:

kubectl -n kube-system logs daemonsets/cilium-envoy

