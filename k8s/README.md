# Kubernetes Deployment (Minikube)

## Prerequisites

- [minikube](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Docker (for building images)

## Setup

### 1. Start Minikube

```bash
minikube start
```

### 2. Elasticsearch (optional but recommended)

Elasticsearch needs `vm.max_map_count` for cluster formation:

```bash
minikube ssh "sudo sysctl -w vm.max_map_count=262144"
```

### 3. Enable Ingress (for GraphQL access)

```bash
minikube addons enable ingress
```

### 4. Build, Load, and Deploy

```bash
make k8s-deploy
```

Or step by step:

```bash
make k8s-build    # Build all images
make k8s-load     # Load images into minikube
make k8s-apply    # Apply k8s manifests
```

## Access GraphQL

After pods are ready:

```bash
minikube tunnel   # Run in separate terminal for LoadBalancer, or use:
kubectl port-forward svc/graphql 8000:80
```

Then open http://localhost:8000/playground

With ingress enabled:

```bash
minikube addons enable ingress
# Add minikube ip to /etc/hosts or use:
curl $(minikube ip)
```

## Useful Commands

```bash
kubectl get pods -w
kubectl logs deployment/account
kubectl describe pod <pod-name>
make k8s-delete   # Remove all k8s resources
```
