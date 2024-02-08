# Learn Kubernetes — The Easy Way

https://programmingpercy.tech/blog/learn-kubernetes-the-easy-way/

* https://kubernetes.io/docs/reference/using-api/
* https://kubernetes.io/docs/concepts/overview/working-with-objects/
* https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
* https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
* https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
* https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/
* https://kubernetes.io/docs/concepts/services-networking/service/
* https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
* https://kubernetes.io/docs/tasks/administer-cluster/namespaces-walkthrough/
* https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/

* minikube start
* kubectl get nodes
* eval $(minikube -p minikube docker-env)
* docker build -t programmingpercy/hellogopher:5.0 .
* kubectl create -f hellogopher.yml
* kubectl get all
* kubectl get deployment/hellogopher -o yaml
* kubectl expose deployment hellogopher --type=NodePort --port=8080
* minikube service hellogopher
* kubectl delete deployment hellogopher
* kubectl get all --show-labels
* kubectl get pods --show-labels
* kubectl label po/hellogopher-f76b49f9-95v4p author=percy
* kubectl label po/hellogopher-56d8758b6b-2rb4d author=ironman --overwrite
* kubectl label po/hellogopher-56d8758b6b-2rb4d author-
* kubectl get pods --selector app=hellogopher
* kubectl get pods --selector app!=hellogopher
* kubectl delete pods -l app=hellogopher ::: kubectl delete pods --selector app=hellogopher
* kubectl describe pod/hellogopher-df787c4d5-gbv66
* kubectl delete service/hellogopher
* kubectl set image deployment/hellogopher hellogopher=programmingpercy/hellogopher:2.0
* kubectl apply -f hellogopher.yml
* kubectl logs pod/hellogopher-f76b49f9-95v4p
* kubectl exec -it pod/hellogopher-79d5bfdfbd-bnhkf -- /bin/sh
* minikube addons enable dashboard
* minikube addons enable metrics-server
* minikube dashboard
* kubectl config set-context --current --namespace=hellogopher
* kubectl apply -f kubernetes/
* kubectl exec pod/mysql-77bd8d464d-8vd2w -it -- bash
* mysql --user=root --password=$MYSQL_ROOT_PASSWORD
* kubectl config set-context --current --namespace=my-namespace
* minikube service hellogopher -n hellogopher
* kubectl create configmap myConfigMap --from-literal=log_level=debug
* kubectl create configmap myConfigMap --from-env-file=path/to/file
* kubectl get configmaps
* kubectl get configmap/myConfigMap -o yaml
* kubectl delete configmap database-configs
* kubectl get secrets
* kubectl get secrets database-secrets -o yaml
* kubectl patch secret database-secrets --type='json' -p='[{"op" : "replace","path" : "/data/DATABASE_PASSWORD","value" : "test"}]'
* minikube logs
* kubectl get all -n hellogopher
* kubectl describe deployment,service,configmap,secret,pod -n hellogopher
* kubectl rollout status deployment/hellogopher -n hellogopher
* kubectl delete -f '*.yml'
* kubectl create -f kubernetes
* kubectl get service mysql -n hellogopher
