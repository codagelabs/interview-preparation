# CKAD Quick Cheat Sheet

## Most Essential Commands for CKAD Exam

### 1. Pod Management
```bash
# Create pod
kubectl run nginx --image=nginx:latest

# Create pod with specific namespace
kubectl run nginx --image=nginx:latest -n my-namespace

# Get pods
kubectl get pods
kubectl get pods -o wide
kubectl get pods -o yaml

# Describe pod
kubectl describe pod <pod-name>

# Delete pod
kubectl delete pod <pod-name>

# Execute in pod
kubectl exec -it <pod-name> -- /bin/bash

# View logs
kubectl logs <pod-name>
kubectl logs -f <pod-name>
```

### 2. Deployment Management
```bash
# Create deployment
kubectl create deployment nginx --image=nginx:latest

# Scale deployment
kubectl scale deployment nginx --replicas=3

# Update image
kubectl set image deployment/nginx nginx=nginx:1.19

# Rollout status
kubectl rollout status deployment/nginx
kubectl rollout history deployment/nginx
kubectl rollout undo deployment/nginx
```

### 3. Service Management
```bash
# Expose deployment as service
kubectl expose deployment nginx --port=80 --target-port=80

# Expose as NodePort
kubectl expose deployment nginx --port=80 --target-port=80 --type=NodePort

# Expose as LoadBalancer
kubectl expose deployment nginx --port=80 --target-port=80 --type=LoadBalancer
```

### 4. ConfigMap and Secrets
```bash
# Create ConfigMap
kubectl create configmap my-config --from-literal=key1=value1

# Create Secret
kubectl create secret generic my-secret --from-literal=username=admin --from-literal=password=secret

# Get ConfigMaps/Secrets
kubectl get configmaps
kubectl get secrets
```

### 5. Namespace Operations
```bash
# Create namespace
kubectl create namespace my-namespace

# Switch to namespace
kubectl config set-context --current --namespace=my-namespace

# Get resources in namespace
kubectl get all -n my-namespace
```

### 6. Resource Management
```bash
# Get all resources
kubectl get all

# Get specific resources
kubectl get pods,services,deployments

# Get resources with labels
kubectl get pods -l app=nginx

# Get resources in all namespaces
kubectl get all --all-namespaces
```

### 7. Debugging Commands
```bash
# Get events
kubectl get events

# Get resource usage
kubectl top pods

# Explain resource
kubectl explain pod

# Port forward
kubectl port-forward <pod-name> 8080:80
```

### 8. Context and Configuration
```bash
# Get current context
kubectl config current-context

# List contexts
kubectl config get-contexts

# Switch context
kubectl config use-context <context-name>
```

## Common Exam Scenarios

### Scenario 1: Create a Multi-container Pod
```bash
# Create base pod
kubectl run nginx --image=nginx:latest --port=80

# Add sidecar container
kubectl patch pod nginx -p '{"spec":{"containers":[{"name":"sidecar","image":"busybox","command":["sleep","3600"]}]}}'
```

### Scenario 2: Create Pod with Volume
```bash
# Create pod with emptyDir volume
kubectl run nginx --image=nginx:latest --overrides='{"spec":{"volumes":[{"name":"shared-data","emptyDir":{}}],"containers":[{"name":"nginx","image":"nginx:latest","volumeMounts":[{"name":"shared-data","mountPath":"/shared"}]}]}}'
```

### Scenario 3: Create Pod with Environment Variables
```bash
# Create pod with env vars
kubectl run nginx --image=nginx:latest --env=DB_HOST=mysql --env=DB_PORT=3306
```

### Scenario 4: Create Pod with Resource Limits
```bash
# Create pod with resource limits
kubectl run nginx --image=nginx:latest --requests=cpu=100m,memory=128Mi --limits=cpu=200m,memory=256Mi
```

## Quick Tips for Exam

1. **Use `kubectl run`** for quick pod creation
2. **Use `kubectl create deployment`** for deployments
3. **Use `kubectl expose`** for services
4. **Use `kubectl scale`** for scaling
5. **Use `kubectl set image`** for image updates
6. **Use `kubectl rollout`** for deployment management
7. **Use `kubectl exec`** for debugging
8. **Use `kubectl logs`** for troubleshooting
9. **Use `kubectl describe`** for detailed info
10. **Use `kubectl get`** with different output formats

## Common Flags to Remember
- `-n` or `--namespace`: Specify namespace
- `-o` or `--output`: Output format (yaml, json, wide)
- `-l` or `--selector`: Label selector
- `-f` or `--follow`: Follow logs
- `-it`: Interactive terminal
- `--all-namespaces`: All namespaces
- `--field-selector`: Field selector
- `--sort-by`: Sort output 