# Codes

* https://kubernetes.io/docs/reference/kubectl/quick-reference/

vim ~/.kube/config

curl --insecure ~~api-server-command~~ [Forbidden]

kubectl get nodes -v6 [Get api server url, kubeconfig file]
kubectl get nodes -v9 [Too verbose]
kubectl get nodes -o wide
kubectl get nodes -o json
kubectl get nodes node1 -o yaml
kubectl describe nodes node1 -o yaml
