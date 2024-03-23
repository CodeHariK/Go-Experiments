# Echo 035

```sh

🐌 Go-Experiments 🦢 k get nodes | grep control
echo-control-plane   Ready    control-plane   27h   v1.29.2

🐌 Go-Experiments 🦢docker exec -it echo-control-plane bash

    root@echo-control-plane:/# ls /etc/kube
    kubelet     kubernetes/ 

    root@echo-control-plane:/# ls /etc/kubernetes/
    admin.conf               kubelet.conf             pki/                     super-admin.conf         
    controller-manager.conf  manifests/               scheduler.conf           

    root@echo-control-plane:/# cat /etc/kubernetes/admin.conf

    root@echo-control-plane:/# cat /etc/kubernetes/admin.conf | grep server
    server: https://echo-control-plane:6443

helm repo add cilium https://helm.cilium.io/

helm show values cilium/cilium | less

helm upgrade cilium cilium/cilium --version 1.15.2 \
    --namespace kube-system \
    --set kubeProxyReplacement=true \
    --set k8sServiceHost="echo-control-plane" \
    --set k8sServicePort=6443 \
    --set loadBalancer.serviceTopology=true \
    --set ipam.mode="kubernetes"

cilium config view
```

## pods.sh
```sh
#!/bin/bash

# The purpose of this script is to deploy to each node in the cluster 2 pods. 
# Each pod will have an env var that shows it's zone.

function echopod () {
  ZONE=""
  case $1 in
    cilium-worker)
      ZONE=a
      ;;
    cilium-worker2)
      ZONE=a
      ;;
    cilium-worker3)
      ZONE=b
      ;;
    cilium-worker4)
      ZONE=b
      ;;
    cilium-worker5)
      ZONE=c
      ;;
    cilium-worker6)
      ZONE=c
      ;;
  esac
  kubectl run echo${2}-${1} \
     --image overridden  --labels app=echo,pod=echo${2}-${1},node=${1},zone=$ZONE  --overrides \
    '{
      "spec":{
        "hostname": "echo'${2}-${1}'",
	      "subdomain": "test",
        "nodeName": "'$1'",
        "containers":[{
          "name":"echo",
          "image":"inanimate/echo-server",
          "env":[{
            "name":"ZONE",
            "value":"'$ZONE'"
          }]
        }]
      }
    }'
}


for worker in $(kind get nodes --name=cilium | grep worker)
  do 
    for i in {1..2}
      do echopod $worker $i
    done
  done

kubectl create service clusterip echo --tcp 8080 

echo "to mark the service eligible for topology aware hints."
echo "kubectl annotate svc echo service.kubernetes.io/topology-aware-hints=auto"
echo "to remove the annotation"
echo "kubectl annotate svc echo service.kubernetes.io/topology-aware-hints-"
echo "to check it's working"
echo "kubectl get endpointslices -o yaml | grep -i for"
```

