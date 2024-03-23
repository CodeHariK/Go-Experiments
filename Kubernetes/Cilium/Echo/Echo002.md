# Echo 001

```sh
hubble status
hubble status --server localhost:4245
hubble status --server localhost:4245 --last 20
hubble observe
hubble observe --server localhost:4245
hubble observe --server localhost:4245 -f
hubble observe --server localhost:4245 -f --output json | jq .
hubble observe --server localhost:4245 --pod kube-system/coredns-76f75df574-hjl96

k -n=kube-system exec -ti cilium-9zs7k -- cilium status
k -n=kube-system exec -ti cilium-9zs7k -- hubble observe --last 10
k -n=kube-system port-forward svc/hubble-relay --address 0.0.0.0 4245:80

hubble observe --server localhost:4245 -f --port 53
k apply -f dns-visibility.yaml
hubble observe --server localhost:4245 -f --port 53 --output json -t l7
hubble observe --server localhost:4245 -f -t l7

k apply -f abchain.yaml

k -n=kube-system port-forward svc/hubble-ui --address 0.0.0.0 12000:80
```