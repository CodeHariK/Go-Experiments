# Github Action, ArgoCD, Kind kubernetes local cluster

```sh

# 1. Create two new github project
mkdir gitaction
git init
gh repo create gitaction --public --source=. --remote=upstream
git remote add origin https://github.com/CodeHariK/gitaction.git
git branch -M main
git push -u origin main

mkdir gitops
git init
gh repo create gitops --public --source=. --remote=upstream
git remote add origin https://github.com/CodeHariK/gitops.git
git branch -M main
git push -u origin main

ssh-keygen -t ed25519 -C "codeharik@gmail.com"

# 2. Create new kind kubernetes cluster
kind create cluster -n gitaction

brew install argocd

kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'

argocd admin initial-password -n argocd

kubectl port-forward svc/argocd-server -n argocd 1111:443

argocd login localhost:1111

kubectl port-forward svc/astro -n gitops 4321:80
```