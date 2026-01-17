# CKAD Imperative Commands Guide

This guide contains all the essential imperative kubectl commands that you need to know for the CKAD exam.

## Table of Contents
1. [Pod Management](#pod-management)
2. [Deployment Management](#deployment-management)
3. [Service Management](#service-management)
4. [ConfigMap and Secret Management](#configmap-and-secret-management)
5. [Namespace Management](#namespace-management)
6. [Resource Quotas and Limits](#resource-quotas-and-limits)
7. [Network Policies](#network-policies)
8. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
9. [Multi-container Pods](#multi-container-pods)
10. [Init Containers](#init-containers)
11. [Probes](#probes)
12. [Volumes and Storage](#volumes-and-storage)
13. [Security Context](#security-context)
14. [Service Accounts](#service-accounts)
15. [RBAC](#rbac)

---

## Pod Management

### Create Pods
```bash
# Create a simple pod
kubectl run nginx-pod --image=nginx:latest

# Create a pod with specific namespace
kubectl run nginx-pod --image=nginx:latest -n my-namespace

# Create a pod with labels
kubectl run nginx-pod --image=nginx:latest --labels=app=nginx,env=prod

# Create a pod with resource limits
kubectl run nginx-pod --image=nginx:latest --requests=cpu=100m,memory=128Mi --limits=cpu=200m,memory=256Mi

# Create a pod with environment variables
kubectl run nginx-pod --image=nginx:latest --env=DB_HOST=mysql --env=DB_PORT=3306

# Create a pod with command override
kubectl run nginx-pod --image=nginx:latest --command -- nginx -g 'daemon off;'

# Create a pod with port exposure
kubectl run nginx-pod --image=nginx:latest --port=80
```

### Pod Operations
```bash
# Get pod information
kubectl get pods
kubectl get pods -o wide
kubectl get pods -o yaml
kubectl get pods -o json

# Describe pod
kubectl describe pod <pod-name>

# Delete pod
kubectl delete pod <pod-name>

# Edit pod
kubectl edit pod <pod-name>

# Execute commands in pod
kubectl exec <pod-name> -- ls /
kubectl exec -it <pod-name> -- /bin/bash

# Copy files to/from pod
kubectl cp <pod-name>:/path/to/file ./local-file
kubectl cp ./local-file <pod-name>:/path/to/file

# View pod logs
kubectl logs <pod-name>
kubectl logs -f <pod-name>  # Follow logs
kubectl logs --previous <pod-name>  # Previous container logs
```

---

## Deployment Management

### Create Deployments
```bash
# Create a deployment
kubectl create deployment nginx-deployment --image=nginx:latest

# Create deployment with replicas
kubectl create deployment nginx-deployment --image=nginx:latest --replicas=3

# Create deployment with labels
kubectl create deployment nginx-deployment --image=nginx:latest --labels=app=nginx,env=prod

# Create deployment with port
kubectl create deployment nginx-deployment --image=nginx:latest --port=80
```

### Deployment Operations
```bash
# Scale deployment
kubectl scale deployment nginx-deployment --replicas=5

# Update deployment image
kubectl set image deployment/nginx-deployment nginx=nginx:1.19

# Rollout operations
kubectl rollout status deployment/nginx-deployment
kubectl rollout history deployment/nginx-deployment
kubectl rollout undo deployment/nginx-deployment
kubectl rollout undo deployment/nginx-deployment --to-revision=2
kubectl rollout pause deployment/nginx-deployment
kubectl rollout resume deployment/nginx-deployment

# Delete deployment
kubectl delete deployment nginx-deployment
```

---

## Service Management

### Create Services
```bash
# Create ClusterIP service
kubectl expose deployment nginx-deployment --port=80 --target-port=80

# Create NodePort service
kubectl expose deployment nginx-deployment --port=80 --target-port=80 --type=NodePort

# Create LoadBalancer service
kubectl expose deployment nginx-deployment --port=80 --target-port=80 --type=LoadBalancer

# Create service with specific name
kubectl expose deployment nginx-deployment --name=nginx-service --port=80 --target-port=80

# Create service with labels
kubectl expose deployment nginx-deployment --port=80 --target-port=80 --labels=app=nginx
```

### Service Operations
```bash
# Get services
kubectl get services
kubectl get svc

# Describe service
kubectl describe service <service-name>

# Delete service
kubectl delete service <service-name>
```

---

## ConfigMap and Secret Management

### Create ConfigMaps
```bash
# Create ConfigMap from literal values
kubectl create configmap my-config --from-literal=key1=value1 --from-literal=key2=value2

# Create ConfigMap from file
kubectl create configmap my-config --from-file=config.properties

# Create ConfigMap from directory
kubectl create configmap my-config --from-file=./config/

# Create ConfigMap from env file
kubectl create configmap my-config --from-env-file=config.env
```

### Create Secrets
```bash
# Create secret from literal values
kubectl create secret generic my-secret --from-literal=username=admin --from-literal=password=secret

# Create secret from file
kubectl create secret generic my-secret --from-file=username.txt --from-file=password.txt

# Create TLS secret
kubectl create secret tls my-tls-secret --cert=tls.crt --key=tls.key

# Create docker-registry secret
kubectl create secret docker-registry my-registry-secret --docker-server=my-registry.com --docker-username=user --docker-password=pass
```

### ConfigMap and Secret Operations
```bash
# Get ConfigMaps and Secrets
kubectl get configmaps
kubectl get secrets

# Describe ConfigMap/Secret
kubectl describe configmap <configmap-name>
kubectl describe secret <secret-name>

# Delete ConfigMap/Secret
kubectl delete configmap <configmap-name>
kubectl delete secret <secret-name>
```

---

## Namespace Management

### Namespace Operations
```bash
# Create namespace
kubectl create namespace my-namespace

# Get namespaces
kubectl get namespaces
kubectl get ns

# Switch context to namespace
kubectl config set-context --current --namespace=my-namespace

# Delete namespace
kubectl delete namespace my-namespace

# Get resources in specific namespace
kubectl get all -n my-namespace
```

---

## Resource Quotas and Limits

### Create Resource Quotas
```bash
# Create resource quota
kubectl create quota my-quota --hard=cpu=4,memory=8Gi,pods=10

# Create quota with scopes
kubectl create quota my-quota --hard=cpu=4,memory=8Gi --scopes=BestEffort
```

### Create Limit Ranges
```bash
# Create limit range
kubectl create limitrange my-limits --min=cpu=100m,memory=128Mi --max=cpu=2,memory=2Gi --default=cpu=500m,memory=512Mi
```

---

## Network Policies

### Create Network Policies
```bash
# Create network policy (requires YAML file)
kubectl apply -f network-policy.yaml

# Get network policies
kubectl get networkpolicies
kubectl get netpol
```

---

## Debugging and Troubleshooting

### Debugging Commands
```bash
# Get events
kubectl get events
kubectl get events --sort-by=.metadata.creationTimestamp

# Get resource usage
kubectl top pods
kubectl top nodes

# Get API resources
kubectl api-resources
kubectl api-versions

# Explain resource
kubectl explain pod
kubectl explain pod.spec.containers

# Get resource in different formats
kubectl get pod <pod-name> -o yaml
kubectl get pod <pod-name> -o json
kubectl get pod <pod-name> -o jsonpath='{.spec.containers[0].image}'
kubectl get pod <pod-name> -o custom-columns=NAME:.metadata.name,IMAGE:.spec.containers[0].image
```

---

## Multi-container Pods

### Create Multi-container Pods
```bash
# Create pod with sidecar container
kubectl run nginx-pod --image=nginx:latest --port=80
kubectl patch pod nginx-pod -p '{"spec":{"containers":[{"name":"sidecar","image":"busybox","command":["sleep","3600"]}]}}'
```

---

## Init Containers

### Create Pods with Init Containers
```bash
# Create pod with init container (requires YAML)
kubectl apply -f init-container-pod.yaml
```

---

## Probes

### Create Pods with Probes
```bash
# Create pod with readiness probe
kubectl run nginx-pod --image=nginx:latest --port=80
kubectl patch pod nginx-pod -p '{"spec":{"containers":[{"name":"nginx","readinessProbe":{"httpGet":{"path":"/","port":80},"initialDelaySeconds":5,"periodSeconds":10}}]}}'
```

---

## Volumes and Storage

### Create Pods with Volumes
```bash
# Create pod with emptyDir volume
kubectl run nginx-pod --image=nginx:latest --overrides='{"spec":{"volumes":[{"name":"shared-data","emptyDir":{}}],"containers":[{"name":"nginx","image":"nginx:latest","volumeMounts":[{"name":"shared-data","mountPath":"/shared"}]}]}}'

# Create pod with configMap volume
kubectl run nginx-pod --image=nginx:latest --overrides='{"spec":{"volumes":[{"name":"config-volume","configMap":{"name":"my-config"}}],"containers":[{"name":"nginx","image":"nginx:latest","volumeMounts":[{"name":"config-volume","mountPath":"/etc/config"}]}]}}'
```

---

## Security Context

### Create Pods with Security Context
```bash
# Create pod with security context
kubectl run nginx-pod --image=nginx:latest --overrides='{"spec":{"securityContext":{"runAsUser":1000,"runAsGroup":3000},"containers":[{"name":"nginx","image":"nginx:latest","securityContext":{"allowPrivilegeEscalation":false}}]}}'
```

---

## Service Accounts

### Service Account Operations
```bash
# Create service account
kubectl create serviceaccount my-serviceaccount

# Get service accounts
kubectl get serviceaccounts
kubectl get sa

# Delete service account
kubectl delete serviceaccount my-serviceaccount
```

---

## RBAC

### RBAC Operations
```bash
# Create role
kubectl create role pod-reader --verb=get,list,watch --resource=pods

# Create cluster role
kubectl create clusterrole pod-reader --verb=get,list,watch --resource=pods

# Create role binding
kubectl create rolebinding read-pods --role=pod-reader --user=jane

# Create cluster role binding
kubectl create clusterrolebinding read-pods-global --clusterrole=pod-reader --user=jane

# Get RBAC resources
kubectl get roles
kubectl get clusterroles
kubectl get rolebindings
kubectl get clusterrolebindings
```

---

## Quick Reference Commands

### Essential Commands for CKAD
```bash
# Get all resources
kubectl get all

# Get resources in all namespaces
kubectl get all --all-namespaces

# Get specific resource
kubectl get pods,services,deployments

# Watch resources
kubectl get pods -w

# Port forward
kubectl port-forward <pod-name> 8080:80

# Get resource in specific namespace
kubectl get pods -n <namespace>

# Label resources
kubectl label pod <pod-name> app=nginx

# Annotate resources
kubectl annotate pod <pod-name> description="nginx web server"

# Get resource with labels
kubectl get pods -l app=nginx

# Get resource with field selectors
kubectl get pods --field-selector status.phase=Running
```

### Context and Configuration
```bash
# Get current context
kubectl config current-context

# List contexts
kubectl config get-contexts

# Switch context
kubectl config use-context <context-name>

# Set namespace for current context
kubectl config set-context --current --namespace=<namespace>
```

---

## Exam Tips

1. **Use `kubectl run`** for quick pod creation
2. **Use `kubectl create`** for more complex resources
3. **Use `kubectl expose`** for service creation
4. **Use `kubectl scale`** for scaling deployments
5. **Use `kubectl set`** for updating images
6. **Use `kubectl rollout`** for deployment management
7. **Use `kubectl exec`** for debugging
8. **Use `kubectl logs`** for troubleshooting
9. **Use `kubectl describe`** for detailed information
10. **Use `kubectl get`** with different output formats

Remember: The CKAD exam focuses on practical skills, so practice these commands extensively! 