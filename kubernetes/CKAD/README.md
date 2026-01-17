# CKAD (Certified Kubernetes Application Developer) Resources

This directory contains comprehensive resources for preparing for the CKAD exam, focusing on imperative commands and practical scenarios.

## 📚 Resources Included

### 1. [Imperative Commands Guide](./imperative_commands.md)
Complete reference of all essential `kubectl` imperative commands organized by topic:
- Pod Management
- Deployment Management  
- Service Management
- ConfigMap and Secret Management
- Namespace Management
- Resource Quotas and Limits
- Network Policies
- Debugging and Troubleshooting
- Multi-container Pods
- Init Containers
- Probes
- Volumes and Storage
- Security Context
- Service Accounts
- RBAC

### 2. [Quick Cheat Sheet](./quick_cheat_sheet.md)
Essential commands for quick reference during the exam:
- Most commonly used commands
- Common exam scenarios
- Quick tips and flags to remember
- Muscle memory commands

### 3. [Practical Examples](./practical_examples.md)
Real-world scenarios with step-by-step solutions:
- Multi-container pods with sidecar pattern
- Pods with ConfigMaps and Secrets
- Deployments with rolling updates
- Services of different types
- Pods with resource limits and probes
- Network policies and security contexts
- Resource quotas and RBAC
- Debugging techniques

## 🎯 CKAD Exam Focus Areas

The CKAD exam tests your ability to:
1. **Design and Build** - Create and configure Kubernetes resources
2. **Deploy** - Deploy applications using different strategies
3. **Expose and Configure** - Create services and configure networking
4. **Troubleshoot** - Debug and troubleshoot applications
5. **Security** - Configure security contexts and RBAC

## 🚀 How to Use These Resources

### For Beginners
1. Start with the [Quick Cheat Sheet](./quick_cheat_sheet.md)
2. Practice the basic commands in each section
3. Move to [Practical Examples](./practical_examples.md) for real scenarios
4. Use [Imperative Commands Guide](./imperative_commands.md) as reference

### For Intermediate Users
1. Focus on [Practical Examples](./practical_examples.md)
2. Practice complex scenarios with multi-container pods
3. Master debugging and troubleshooting techniques
4. Understand security and RBAC concepts

### For Advanced Users
1. Practice all scenarios in [Practical Examples](./practical_examples.md)
2. Master the complete [Imperative Commands Guide](./imperative_commands.md)
3. Focus on time management during practice
4. Simulate exam conditions

## 📋 Essential Commands to Master

### Quick Pod Creation
```bash
kubectl run nginx --image=nginx:latest
kubectl run nginx --image=nginx:latest -n my-namespace
kubectl run nginx --image=nginx:latest --port=80
```

### Deployment Management
```bash
kubectl create deployment nginx --image=nginx:latest
kubectl scale deployment nginx --replicas=3
kubectl set image deployment/nginx nginx=nginx:1.19
kubectl rollout status deployment/nginx
```

### Service Creation
```bash
kubectl expose deployment nginx --port=80 --target-port=80
kubectl expose deployment nginx --port=80 --target-port=80 --type=NodePort
```

### Configuration Management
```bash
kubectl create configmap my-config --from-literal=key1=value1
kubectl create secret generic my-secret --from-literal=username=admin
```

### Debugging
```bash
kubectl get pods -o wide
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl exec -it <pod-name> -- /bin/bash
```

## ⏱️ Exam Time Management Tips

1. **Use imperative commands** - They're faster than YAML for simple tasks
2. **Master `kubectl run`** - Quick pod creation
3. **Use `kubectl expose`** - Fast service creation
4. **Practice `kubectl patch`** - Modify existing resources
5. **Know your debugging commands** - `kubectl describe`, `kubectl logs`, `kubectl exec`

## 🎯 Key Exam Topics Covered

### Core Concepts (25%)
- Multi-container Pod Design
- Pod Design
- Observability
- Services & Networking

### Configuration (25%)
- ConfigMaps
- Security Contexts
- Resource Requirements
- Secrets

### Multi-container Pods (10%)
- Sidecar Containers
- Init Containers
- Ambassador Containers
- Adapter Containers

### Pod Design (20%)
- Labels and Annotations
- Deployments
- Jobs and CronJobs
- Health Checks

### State Persistence (10%)
- Volumes
- Persistent Volumes
- Persistent Volume Claims

### Services & Networking (10%)
- Services
- Network Policies
- Ingress
- Service Mesh

## 🔧 Practice Environment Setup

### Local Kubernetes Cluster
```bash
# Install minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start cluster
minikube start

# Enable addons
minikube addons enable ingress
minikube addons enable metrics-server
```

### Alternative: Docker Desktop
- Enable Kubernetes in Docker Desktop settings
- Use the built-in Kubernetes cluster

### Alternative: Kind
```bash
# Install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Create cluster
kind create cluster
```

## 📝 Exam Preparation Checklist

- [ ] Master all imperative commands in the cheat sheet
- [ ] Practice all scenarios in practical examples
- [ ] Understand multi-container pod patterns
- [ ] Know how to create and manage ConfigMaps and Secrets
- [ ] Practice debugging techniques
- [ ] Understand service types and networking
- [ ] Know RBAC and security contexts
- [ ] Practice time management
- [ ] Take mock exams under time pressure

## 🎓 Additional Resources

- [Official CKAD Curriculum](https://github.com/cncf/curriculum)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [CKAD Practice Tests](https://github.com/arush-sal/cka-practice-environment)

## 💡 Pro Tips

1. **Practice daily** - Even 30 minutes of practice helps
2. **Use imperative commands** - They're faster for the exam
3. **Know your debugging tools** - `kubectl describe` and `kubectl logs` are your friends
4. **Time management** - Don't spend too long on any single question
5. **Read carefully** - Understand what the question is asking for
6. **Verify your work** - Use `kubectl get` and `kubectl describe` to check your work

Good luck with your CKAD exam! 🚀 