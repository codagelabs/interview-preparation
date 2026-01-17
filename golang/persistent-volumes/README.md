# Kubernetes Persistent Volumes

This directory contains configurations and examples for Kubernetes Persistent Volumes (PVs) and Persistent Volume Claims (PVCs).

## Overview

Persistent Volumes (PVs) are a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. They are a resource in the cluster just like a node is a cluster resource.

## Directory Structure

```
persistent-volumes/
├── README.md                 # This file
├── static/                   # Manually provisioned PV configurations
├── dynamic/                  # Dynamically provisioned PV configurations
└── examples/                 # Example configurations and use cases
```

## Key Concepts

### Persistent Volume (PV)
- A storage resource in the cluster
- Can be provisioned statically or dynamically
- Has a specific storage capacity
- Has access modes (ReadWriteOnce, ReadOnlyMany, ReadWriteMany)
- Has a reclaim policy (Retain, Recycle, Delete)

### Persistent Volume Claim (PVC)
- A request for storage by a user
- Similar to a pod consuming node resources
- Can request specific size and access modes

### Storage Class
- Describes the "classes" of storage offered
- Different classes might map to quality-of-service levels
- Can be used for dynamic provisioning

## PV vs PVC: Key Differences

| Aspect | Persistent Volume (PV) | Persistent Volume Claim (PVC) |
|--------|------------------------|-------------------------------|
| **Definition** | Cluster resource representing actual storage | Request for storage by a user |
| **Creation** | Created by cluster administrator | Created by application developer/user |
| **Purpose** | Represents actual storage in the cluster | Requests storage from available PVs |
| **Lifecycle** | Independent of pods | Tied to pod lifecycle |
| **Binding** | Can exist without being bound | Must be bound to a PV to be used |
| **Access Control** | Cluster-wide resource | Namespace-scoped resource |
| **Storage Details** | Contains actual storage implementation details | Specifies storage requirements |
| **Example** | Physical storage like NFS server, EBS volume | Request for 5Gi of storage with RWX access |

### Key Points:
1. **PV is the Resource, PVC is the Request**
   - PV is the actual storage resource in the cluster
   - PVC is a request for storage with specific requirements

2. **Binding Process**
   - PVCs are bound to PVs based on matching criteria
   - One PV can only be bound to one PVC
   - One PVC can only be bound to one PV

3. **Access Control**
   - PVs are cluster-wide resources
   - PVCs are namespace-scoped
   - PVCs provide namespace isolation for storage

4. **Lifecycle Management**
   - PVs can exist independently of PVCs
   - PVCs are typically created and deleted with applications
   - PVs can be retained, recycled, or deleted based on reclaim policy

## Static vs Dynamic PV Provisioning

### Static Provisioning
| Aspect | Description |
|--------|-------------|
| **Definition** | Administrator manually creates PVs in advance |
| **Process** | 1. Admin creates PVs<br>2. User creates PVC<br>3. System binds PVC to matching PV |
| **When to Use** | - Predictable storage needs<br>- Specific storage requirements<br>- Pre-existing storage infrastructure |
| **Advantages** | - More control over storage<br>- Can use existing storage systems<br>- Predictable capacity planning |
| **Disadvantages** | - Manual management required<br>- May lead to underutilization<br>- Scaling requires manual intervention |
| **Example** | Pre-provisioned NFS shares or EBS volumes |

### Dynamic Provisioning
| Aspect | Description |
|--------|-------------|
| **Definition** | PVs are automatically created when PVCs are created |
| **Process** | 1. Admin creates StorageClass<br>2. User creates PVC<br>3. System automatically provisions PV |
| **When to Use** | - Unpredictable storage needs<br>- Cloud environments<br>- Need for automation |
| **Advantages** | - Automatic PV creation<br>- Better resource utilization<br>- Easier scaling<br>- Reduced administrative overhead |
| **Disadvantages** | - Less control over specific PVs<br>- May incur unexpected costs<br>- Requires compatible storage backend |
| **Example** | Cloud storage like AWS EBS or Azure Disk |

### Key Differences
1. **Provisioning Time**
   - Static: PVs exist before PVCs
   - Dynamic: PVs created on-demand with PVCs

2. **Management Overhead**
   - Static: Higher administrative overhead
   - Dynamic: Automated, lower overhead

3. **Flexibility**
   - Static: More control over specific PVs
   - Dynamic: More flexible for scaling

4. **Use Cases**
   - Static: Traditional storage, specific requirements
   - Dynamic: Cloud-native applications, auto-scaling

## Complete PV Provisioning Process with Examples

### Static Provisioning Process

1. **Create NFS Server (Example)**
```bash
# On NFS Server
sudo apt-get update
sudo apt-get install nfs-kernel-server
sudo mkdir -p /mnt/nfs_share
sudo chown nobody:nogroup /mnt/nfs_share
sudo chmod 777 /mnt/nfs_share
echo "/mnt/nfs_share *(rw,sync,no_subtree_check)" | sudo tee -a /etc/exports
sudo exportfs -a
sudo systemctl restart nfs-kernel-server
```

2. **Create PV (Example)**
```yaml
# static-pv.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: static-nfs-pv
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteMany
  persistentVolumeReclaimPolicy: Retain
  storageClassName: manual
  nfs:
    server: 192.168.1.100  # NFS server IP
    path: "/mnt/nfs_share"
```

3. **Create PVC (Example)**
```yaml
# static-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: static-nfs-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 5Gi
  storageClassName: manual
```

4. **Apply Configurations**
```bash
# Create PV
kubectl apply -f static-pv.yaml

# Create PVC
kubectl apply -f static-pvc.yaml

# Verify binding
kubectl get pv
kubectl get pvc
```

### Dynamic Provisioning Process

1. **Create StorageClass (Example for AWS EBS)**
```yaml
# storage-class.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: aws-gp2
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
```

2. **Create PVC (Example)**
```yaml
# dynamic-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: dynamic-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: aws-gp2
```

3. **Apply Configurations**
```bash
# Create StorageClass
kubectl apply -f storage-class.yaml

# Create PVC (PV will be automatically created)
kubectl apply -f dynamic-pvc.yaml

# Verify PV creation and binding
kubectl get pv
kubectl get pvc
```

### Using PVs in Pods

1. **Static PV Example Pod**
```yaml
# static-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: static-pod
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: static-volume
      mountPath: /usr/share/nginx/html
  volumes:
  - name: static-volume
    persistentVolumeClaim:
      claimName: static-nfs-pvc
```

2. **Dynamic PV Example Pod**
```yaml
# dynamic-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dynamic-pod
spec:
  containers:
  - name: nginx
    image: nginx
    volumeMounts:
    - name: dynamic-volume
      mountPath: /usr/share/nginx/html
  volumes:
  - name: dynamic-volume
    persistentVolumeClaim:
      claimName: dynamic-pvc
```

### Verification Commands

```bash
# Check PVs
kubectl get pv

# Check PVCs
kubectl get pvc

# Check StorageClasses
kubectl get storageclass

# Check Pods and their volumes
kubectl get pods -o wide
kubectl describe pod <pod-name>

# Check volume mounts in pod
kubectl exec -it <pod-name> -- df -h
```

### Cleanup Process

1. **Static Provisioning Cleanup**
```bash
# Delete pod
kubectl delete -f static-pod.yaml

# Delete PVC
kubectl delete -f static-pvc.yaml

# Delete PV
kubectl delete -f static-pv.yaml
```

2. **Dynamic Provisioning Cleanup**
```bash
# Delete pod
kubectl delete -f dynamic-pod.yaml

# Delete PVC (PV will be automatically deleted)
kubectl delete -f dynamic-pvc.yaml

# Delete StorageClass
kubectl delete -f storage-class.yaml
```

## Types of Persistent Volumes

### Local Storage Types
1. **hostPath**
   - Uses a directory on the host node
   - Good for development and testing
   - Not suitable for production
   - Example: Local development environments

2. **local**
   - Similar to hostPath but with better scheduling
   - Supports dynamic provisioning
   - Good for single-node clusters
   - Example: Local SSD storage

### Network Storage Types
1. **NFS**
   - Network File System
   - Supports ReadWriteMany
   - Good for shared storage
   - Example: Shared file systems

2. **iSCSI**
   - Internet Small Computer System Interface
   - Block storage over IP
   - Good for database storage
   - Example: Database volumes

3. **CephFS**
   - Distributed file system
   - Supports ReadWriteMany
   - Good for large-scale storage
   - Example: Large file storage

4. **GlusterFS**
   - Distributed file system
   - Supports ReadWriteMany
   - Good for shared storage
   - Example: Shared application data

### Cloud Provider Storage Types
1. **AWS EBS**
   - Elastic Block Store
   - Block storage for AWS
   - Good for AWS deployments
   - Example: AWS-based applications

2. **GCP Persistent Disk**
   - Block storage for GCP
   - Good for GCP deployments
   - Example: GCP-based applications

3. **Azure Disk**
   - Block storage for Azure
   - Good for Azure deployments
   - Example: Azure-based applications

### Container Storage Interface (CSI) Volumes
1. **CSI Drivers**
   - Standard interface for storage plugins
   - Supports various storage backends
   - Good for custom storage solutions
   - Example: Custom storage providers

2. **Common CSI Drivers**
   - Longhorn
   - Rook
   - Portworx
   - OpenEBS

## Common Use Cases

1. **Database Storage**: Persistent storage for database applications
2. **File Sharing**: Shared file systems across multiple pods
3. **Log Storage**: Storing application logs persistently
4. **Data Backup**: Storing backup data

## Best Practices

1. **Access Modes**:
   - ReadWriteOnce (RWO): Can be mounted as read-write by a single node
   - ReadOnlyMany (ROX): Can be mounted read-only by many nodes
   - ReadWriteMany (RWX): Can be mounted as read-write by many nodes

2. **Reclaim Policies**:
   - Retain: Manual reclamation
   - Recycle: Basic scrub (rm -rf /thevolume/*)
   - Delete: Associated storage asset deleted

3. **Storage Classes**:
   - Define different classes of storage
   - Enable dynamic provisioning
   - Specify parameters for the provisioner

## Configuration Examples

### Static Provisioning
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: example-pv
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: standard
  hostPath:
    path: "/mnt/data"
```

### Dynamic Provisioning
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: example-pvc
spec:
  accessModes:
    - ReadWriteOnce 
  resources:
    requests:
      storage: 5Gi
  storageClassName: standard
```

## Troubleshooting

Common issues and solutions:
1. **PV not binding to PVC**: Check storage class, access modes, and capacity
2. **Permission issues**: Verify security context and access modes
3. **Storage not available**: Check underlying storage system status

## Additional Resources

- [Kubernetes Persistent Volumes Documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [Storage Classes Documentation](https://kubernetes.io/docs/concepts/storage/storage-classes/)
- [Volume Types Documentation](https://kubernetes.io/docs/concepts/storage/volumes/) 