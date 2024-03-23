# Echo 001

```sh
alias k='kubectl'

curl -LO https://raw.githubusercontent.com/cilium/cilium/1.15.2/Documentation/installation/kind-config.yaml
kind create cluster --config=kind-config.yaml -n echo

cilium install --version 1.15.2

cilium status --wait
cilium status

cilium hubble enable
cilium hubble enable --ui
cilium hubble ui
cilium hubble port-forward
hubble status
hubble status --server localhost:4245
hubble observe
hubble observe --server localhost:4245
hubble observe --server localhost:4245 --pod kube-system/coredns-76f75df574-hjl96

k get po -A -o wide | grep cilium
kube-system          cilium-9zs7k                                 1/1     Running   0          35m   172.20.0.4     echo-worker3      
kube-system          cilium-hfvkj                                 1/1     Running   0          35m   172.20.0.3     echo-worker2      
kube-system          cilium-j64cb                                 1/1     Running   0          35m   172.20.0.5     echo-control-plane
kube-system          cilium-l6xsp                                 1/1     Running   0          35m   172.20.0.2     echo-worker       
kube-system          cilium-operator-6cdc4568cb-nbsgn             1/1     Running   0          35m   172.20.0.4     echo-worker3      

k -n=kube-system  exec -ti cilium-9zs7k -- bash

root@echo-worker3:/home/cilium# cilium endpoint list
ENDPOINT   POLICY (ingress)   POLICY (egress)   IDENTITY   LABELS (source:key[=value])   IPv6   IPv4           STATUS   
           ENFORCEMENT        ENFORCEMENT                                                                      
2425       Disabled           Disabled          1          reserved:host                                       ready   
3373       Disabled           Disabled          4          reserved:health                      10.244.3.243   ready

root@echo-worker3:/home/cilium# cilium bpf endpoint list
IP ADDRESS       LOCAL ENDPOINT INFO
10.244.3.212:0   (localhost)                                                                                      
172.20.0.4:0     (localhost)                                                                                      
10.244.3.243:0   id=3373  sec_id=4     flags=0x0000 ifindex=8   mac=FE:B4:CA:C1:F7:5E nodemac=F2:CD:55:72:41:BD   

root@echo-worker3:/home/cilium# cilium service list
ID   Frontend            Service Type   Backend                           
1    10.96.0.1:443       ClusterIP      1 => 172.20.0.5:6443 (active)     
2    10.96.170.158:443   ClusterIP      1 => 172.20.0.4:4244 (active)     
3    10.96.0.10:53       ClusterIP      1 => 10.244.0.238:53 (active)     
                                        2 => 10.244.0.180:53 (active)     
4    10.96.0.10:9153     ClusterIP      1 => 10.244.0.238:9153 (active)   
                                        2 => 10.244.0.180:9153 (active)   
5    10.96.212.166:80    ClusterIP      1 => 10.244.1.54:4245 (active)    
6    10.96.238.89:80     ClusterIP      1 => 10.244.2.79:8081 (active)

🐌 EchoExamples 🦢 k get svc -A
NAMESPACE     NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                  AGE
default       kubernetes     ClusterIP   10.96.0.1       <none>        443/TCP                  61m
kube-system   hubble-peer    ClusterIP   10.96.170.158   <none>        443/TCP                  59m
kube-system   hubble-relay   ClusterIP   10.96.212.166   <none>        80/TCP                   55m
kube-system   hubble-ui      ClusterIP   10.96.238.89    <none>        80/TCP                   42m
kube-system   kube-dns       ClusterIP   10.96.0.10      <none>        53/UDP,53/TCP,9153/TCP   61m

root@echo-worker3:/home/cilium# cilium bpf lb list
SERVICE ADDRESS     BACKEND ADDRESS (REVNAT_ID) (SLOT)
10.96.0.10:9153     0.0.0.0:0 (4) (0) [ClusterIP, non-routable]                  
                    10.244.0.180:9153 (4) (2)                                    
                    10.244.0.238:9153 (4) (1)                                    
10.96.212.166:80    0.0.0.0:0 (5) (0) [ClusterIP, non-routable]                  
                    10.244.1.54:4245 (5) (1)                                     
10.96.0.1:443       0.0.0.0:0 (1) (0) [ClusterIP, non-routable]                  
                    172.20.0.5:6443 (1) (1)                                      
10.96.0.10:53       0.0.0.0:0 (3) (0) [ClusterIP, non-routable]                  
                    10.244.0.238:53 (3) (1)                                      
                    10.244.0.180:53 (3) (2)                                      
10.96.238.89:80     10.244.2.79:8081 (6) (1)                                     
                    0.0.0.0:0 (6) (0) [ClusterIP, non-routable]                  
10.96.170.158:443   172.20.0.4:4244 (2) (1)                                      
                    0.0.0.0:0 (2) (0) [ClusterIP, InternalLocal, non-routable]  
```