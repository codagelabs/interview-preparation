# CKAD Practical Examples

## Common Exam Scenarios and Solutions

### 1. Multi-container Pod with Sidecar Pattern

**Scenario**: Create a pod with nginx and a sidecar container for logging.

```bash
# Create base pod
kubectl run nginx-sidecar --image=nginx:latest --port=80

# Add sidecar container using patch
kubectl patch pod nginx-sidecar -p '{
  "spec": {
    "containers": [
      {
        "name": "nginx",
        "image": "nginx:latest",
        "ports": [{"containerPort": 80}]
      },
      {
        "name": "sidecar",
        "image": "busybox",
        "command": ["sh", "-c", "while true; do echo $(date) >> /shared/logs.txt; sleep 10; done"],
        "volumeMounts": [{"name": "shared-volume", "mountPath": "/shared"}]
      }
    ],
    "volumes": [{"name": "shared-volume", "emptyDir": {}}]
  }
}'
```

### 2. Pod with ConfigMap and Environment Variables

**Scenario**: Create a pod that uses ConfigMap for configuration.

```bash
# Create ConfigMap
kubectl create configmap app-config --from-literal=DB_HOST=mysql --from-literal=DB_PORT=3306 --from-literal=APP_ENV=production

# Create pod with ConfigMap
kubectl run app-pod --image=myapp:latest --overrides='{
  "spec": {
    "containers": [{
      "name": "app",
      "image": "myapp:latest",
      "envFrom": [{"configMapRef": {"name": "app-config"}}]
    }]
  }
}'
```

### 3. Pod with Secrets

**Scenario**: Create a pod that uses secrets for sensitive data.

```bash
# Create secret
kubectl create secret generic db-secret --from-literal=username=admin --from-literal=password=secret123

# Create pod with secret
kubectl run db-app --image=myapp:latest --overrides='{
  "spec": {
    "containers": [{
      "name": "app",
      "image": "myapp:latest",
      "env": [
        {"name": "DB_USER", "valueFrom": {"secretKeyRef": {"name": "db-secret", "key": "username"}}},
        {"name": "DB_PASS", "valueFrom": {"secretKeyRef": {"name": "db-secret", "key": "password"}}}
      ]
    }]
  }
}'
```

### 4. Deployment with Rolling Update Strategy

**Scenario**: Create a deployment with specific update strategy.

```bash
# Create deployment with rolling update
kubectl create deployment nginx-deployment --image=nginx:1.18 --replicas=3

# Update strategy using patch
kubectl patch deployment nginx-deployment -p '{
  "spec": {
    "strategy": {
      "type": "RollingUpdate",
      "rollingUpdate": {
        "maxSurge": 1,
        "maxUnavailable": 0
      }
    }
  }
}'
```

### 5. Service with Different Types

**Scenario**: Create services of different types for the same deployment.

```bash
# Create deployment
kubectl create deployment web-app --image=nginx:latest --replicas=3

# Create ClusterIP service (internal)
kubectl expose deployment web-app --name=web-app-internal --port=80 --target-port=80

# Create NodePort service (external access)
kubectl expose deployment web-app --name=web-app-external --port=80 --target-port=80 --type=NodePort

# Create LoadBalancer service (cloud load balancer)
kubectl expose deployment web-app --name=web-app-lb --port=80 --target-port=80 --type=LoadBalancer
```

### 6. Pod with Resource Limits and Requests

**Scenario**: Create a pod with specific resource requirements.

```bash
# Create pod with resource limits
kubectl run resource-pod --image=nginx:latest --requests=cpu=100m,memory=128Mi --limits=cpu=200m,memory=256Mi
```

### 7. Pod with Health Checks (Probes)

**Scenario**: Create a pod with readiness and liveness probes.

```bash
# Create pod with probes
kubectl run healthy-pod --image=nginx:latest --port=80

# Add probes using patch
kubectl patch pod healthy-pod -p '{
  "spec": {
    "containers": [{
      "name": "nginx",
      "image": "nginx:latest",
      "ports": [{"containerPort": 80}],
      "readinessProbe": {
        "httpGet": {"path": "/", "port": 80},
        "initialDelaySeconds": 5,
        "periodSeconds": 10
      },
      "livenessProbe": {
        "httpGet": {"path": "/", "port": 80},
        "initialDelaySeconds": 15,
        "periodSeconds": 20
      }
    }]
  }
}'
```

### 8. Pod with Volume Mounts

**Scenario**: Create a pod with persistent storage.

```bash
# Create pod with volume
kubectl run volume-pod --image=nginx:latest --overrides='{
  "spec": {
    "containers": [{
      "name": "nginx",
      "image": "nginx:latest",
      "volumeMounts": [{"name": "data-volume", "mountPath": "/data"}]
    }],
    "volumes": [{"name": "data-volume", "emptyDir": {}}]
  }
}'
```

### 9. Pod with Security Context

**Scenario**: Create a pod with specific security settings.

```bash
# Create pod with security context
kubectl run secure-pod --image=nginx:latest --overrides='{
  "spec": {
    "securityContext": {
      "runAsUser": 1000,
      "runAsGroup": 3000,
      "fsGroup": 2000
    },
    "containers": [{
      "name": "nginx",
      "image": "nginx:latest",
      "securityContext": {
        "allowPrivilegeEscalation": false,
        "readOnlyRootFilesystem": true
      }
    }]
  }
}'
```

### 10. Network Policy

**Scenario**: Create a network policy to control traffic.

```bash
# Create network policy (requires YAML file)
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-nginx
spec:
  podSelector:
    matchLabels:
      app: nginx
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 80
EOF
```

### 11. Resource Quotas and Limits

**Scenario**: Create resource quotas for a namespace.

```bash
# Create namespace
kubectl create namespace quota-demo

# Create resource quota
kubectl create quota compute-quota --hard=cpu=4,memory=8Gi,pods=10 -n quota-demo

# Create limit range
kubectl create limitrange memory-limit --min=memory=50Mi --max=memory=1Gi --default=memory=100Mi -n quota-demo
```

### 12. Service Account and RBAC

**Scenario**: Create service account with specific permissions.

```bash
# Create service account
kubectl create serviceaccount app-sa

# Create role
kubectl create role pod-reader --verb=get,list,watch --resource=pods

# Create role binding
kubectl create rolebinding read-pods --role=pod-reader --serviceaccount=default:app-sa
```

### 13. Init Containers

**Scenario**: Create a pod with init containers for setup.

```bash
# Create pod with init container (requires YAML)
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
  initContainers:
  - name: init-setup
    image: busybox
    command: ['sh', '-c', 'echo "Setup complete" > /shared/init.txt']
    volumeMounts:
    - name: shared-volume
      mountPath: /shared
  volumes:
  - name: shared-volume
    emptyDir: {}
EOF
```

### 14. Debugging Commands

**Scenario**: Common debugging scenarios.

```bash
# Check pod status
kubectl get pods -o wide

# View pod events
kubectl describe pod <pod-name>

# Check pod logs
kubectl logs <pod-name>
kubectl logs -f <pod-name>

# Execute commands in pod
kubectl exec -it <pod-name> -- /bin/bash

# Port forward for debugging
kubectl port-forward <pod-name> 8080:80

# Get resource usage
kubectl top pods
kubectl top nodes

# Check events
kubectl get events --sort-by=.metadata.creationTimestamp
```

### 15. Advanced Pod Configuration

**Scenario**: Complex pod with multiple features.

```bash
# Create complex pod with multiple features
kubectl run complex-pod --image=nginx:latest --overrides='{
  "spec": {
    "containers": [{
      "name": "nginx",
      "image": "nginx:latest",
      "ports": [{"containerPort": 80}],
      "env": [
        {"name": "ENV", "value": "production"},
        {"name": "VERSION", "value": "1.0"}
      ],
      "resources": {
        "requests": {"cpu": "100m", "memory": "128Mi"},
        "limits": {"cpu": "200m", "memory": "256Mi"}
      },
      "volumeMounts": [{"name": "config-volume", "mountPath": "/etc/config"}]
    }],
    "volumes": [{"name": "config-volume", "emptyDir": {}}]
  }
}'
```

## Quick Reference for Common Tasks

### Create Resources
```bash
# Pod
kubectl run <name> --image=<image>

# Deployment
kubectl create deployment <name> --image=<image>

# Service
kubectl expose deployment <name> --port=<port>

# ConfigMap
kubectl create configmap <name> --from-literal=<key>=<value>

# Secret
kubectl create secret generic <name> --from-literal=<key>=<value>

# Namespace
kubectl create namespace <name>
```

### Manage Resources
```bash
# Scale
kubectl scale deployment <name> --replicas=<number>

# Update image
kubectl set image deployment/<name> <container>=<image>

# Rollout
kubectl rollout status deployment/<name>
kubectl rollout undo deployment/<name>
```

### Debug Resources
```bash
# Get info
kubectl get <resource>
kubectl describe <resource> <name>

# Logs
kubectl logs <pod-name>

# Execute
kubectl exec -it <pod-name> -- /bin/bash

# Port forward
kubectl port-forward <pod-name> <local-port>:<pod-port>
```

Remember: Practice these scenarios repeatedly to build muscle memory for the exam! 