# Networking

* Install crictl https://github.com/kubernetes-sigs/cri-tools
* https://kubernetes.io/docs/concepts/cluster-administration/addons/

```sh
kind create cluster --name kubenet
kind get clusters
kubectl get nodes
kubectl cluster-info --context kind-kubenet

docker ps
CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS         PORTS                       NAMES
bcf7cb9bcb8d   kindest/node:v1.29.2   "/usr/local/bin/entr…"   6 minutes ago   Up 6 minutes   127.0.0.1:49791->6443/tcp   kubenet-control-plane

docker exec -it bcf7cb9bcb8d crictl images

kind delete cluster --name kubenet
```

```sh
# https://docs.cilium.io/en/stable/gettingstarted/k8s-install-default/
curl -LO https://raw.githubusercontent.com/cilium/cilium/1.15.1/Documentation/installation/kind-config.yaml

kind create cluster --config=kind-config.yaml --name kubenet

cilium install --version 1.15.1

cilium status --wait

kubectl --namespace kube-system get ds

kubectl --namespace kube-system get pods --selector k8s-app=cilium --sort-by='.status.containerStatuses[0].restartCount'

kubectl get pods -A

cilium config view

kubectl run web --image=httpd
kubectl run client -it --image=busybox

kubectl get pod -o wide
NAME     READY   STATUS    RESTARTS      AGE   IP             NODE             NOMINATED NODE   READINESS GATES
client   1/1     Running   1 (73s ago)   85s   10.244.1.119   kubenet-worker   <none>           <none>
web      1/1     Running   0             7s    10.244.1.129   kubenet-worker   <none>           <none>

kubectl exec -it client -- sh
/ # ping 10.244.1.129

kubectl delete pod web
kubectl delete pod client

kubectl create deployment website --replicas=3 --image=httpd

kubectl get po --show-labels

kuebctl apply -f service.yaml

kubectl get service

kubectl get all

kubectl get ns

kubectl cluster-info dump | grep -m 1 service-cluster-ip-range

kubectl get endpoints website

sudo iptables -L -vn -t nat | grep '10.111.148.30'
sudo iptables -L -vn -t nat | grep -A4 'Chain KUBE-SVC-RYQJBQ5TR32XWAUN'

kubectl logs website-5746f499f-qzjjq
kubectl logs -l app=website

kubectl get all -n kube-system

kubectl get sts,po

kubectl attach -it client

kubectl expose deployment website --port=80

/ # nslookup website
/ # nslookup -q=srv website

/ #  wget -qO - website
<html><body><h1>It works!</h1></body></html>
/ # wget -qO - website.kube-system
wget: bad address 'website.kube-system'
/ # wget -qO - website.default
<html><body><h1>It works!</h1></body></html>
/ # wget -qO - website.default.svc
<html><body><h1>It works!</h1></body></html>
/ # wget -qO -  website.default.svc.cluster.local
<html><body><h1>It works!</h1></body></html>
/ # wget -O -  website.default.svc.cluster.local

/ # cat /etc/resolv.conf

# https://www.getambassador.io/docs/emissary/latest/tutorials/getting-started
kubectl create namespace emissary && \
kubectl apply -f https://app.getambassador.io/yaml/emissary/3.9.1/emissary-crds.yaml && \
kubectl wait --timeout=90s --for=condition=available deployment emissary-apiext -n emissary-system
 
kubectl apply -f https://app.getambassador.io/yaml/emissary/3.9.1/emissary-emissaryns.yaml && \
kubectl -n emissary wait --for condition=available --timeout=90s deploy -lproduct=aes

# Linkerd
curl -sSfL https://run.linkerd.io/install | sh

```
