# Service

A Kubernetes Service is an abstraction that defines a logical set of Pods and a policy by which to access them. Services enable loose coupling between dependent Pods.

## Key Features of a Kubernetes Service

- **Stable Network Identity**: Services provide a stable IP address and DNS name for a set of Pods. This allows other Pods to reliably communicate with the Service, even if the underlying Pods change.

- **Load Balancing**: Services can distribute traffic across multiple Pods, ensuring that the load is balanced and no single Pod is overwhelmed.

- **Service Discovery**: Kubernetes provides built-in service discovery mechanisms, allowing Pods to find and communicate with Services easily.

## Types of Services

Kubernetes supports several types of Services, each suited for different use cases:

1. **ClusterIP** (default):
    - his default Service type assigns an IP address from a pool of IP addresses that your cluster has reserved for that purpose.
    - Exposes the Service on a cluster-internal IP.
    - Makes the Service only reachable from within the cluster.
    - Suitable for internal communication between Pods.
    - If you define a Service that has the .spec.clusterIP set to "None" then Kubernetes does not assign an IP address.

   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: clusterip-service
   spec:
     selector:
       app: example
     ports:
       - protocol: TCP
         port: 80
         targetPort: 9376
     type: ClusterIP
   ```

2. **NodePort**:

   - Exposes the Service on each Node's IP at a static port.
   - Makes the Service accessible from outside the cluster using `<NodeIP>:<NodePort>`.
   - Useful for simple, external access to Services.


   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: nodeport-service
   spec:
     selector:
       app: example
     ports:
       - protocol: TCP
         port: 80
         targetPort: 9376
         nodePort: 30007
     type: NodePort
   ```

3. **LoadBalancer**:
   - Externally exposes the Service using a cloud provider's load balancer.
   - Automatically provisions a load balancer for the Service.
   - Ideal for production environments where external access is required.

   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: loadbalancer-service
   spec:
     selector:
       app: example
     ports:
       - protocol: TCP
         port: 80
         targetPort: 9376
     type: LoadBalancer
   ```

4. **ExternalName**:
   - Maps the Service to the contents of the `externalName` field (e.g., `foo.bar.example.com`).
   - Returns a CNAME record with the name.
   - Useful for integrating with external services not managed by Kubernetes.

   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: externalname-service
   spec:
     type: ExternalName
     externalName: foo.bar.example.com
   ```

## Example: Basic Service YAML

Below is an example of a basic Service YAML configuration:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: example-service
spec:
  selector:
    app: example
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376
  type: ClusterIP
```

### Description
- **apiVersion**: Specifies the version of the Kubernetes API to use. Here, `v1` is used for core resources like Services.
- **kind**: Defines the type of Kubernetes resource. In this case, it's a `Service`.
- **metadata**: Contains data that helps uniquely identify the Service, such as its `name`.
- **spec**: Describes the desired state of the Service.
  - **selector**: Defines how the Service finds the Pods it routes traffic to, using labels.
  - **ports**: Specifies the ports that the Service should expose and the target ports on the Pods.
  - **type**: Defines the type of Service, which determines how the Service is exposed.

This YAML file defines a simple Service that routes traffic to Pods with the label `app: example`, exposing port 80 and targeting port 9376 on the Pods.