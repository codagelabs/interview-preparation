# Pods
Pods are the smallest deployable units of computing that you can create and manage in Kubernetes.

A Pod is a group of one or more containers, with shared storage and network resources, and a specification for how to run the containers.

## Usage of Pods in a Kubernetes Cluster

Pods in a Kubernetes cluster are used in two main ways:

1. **Single Container Pod**: This is the most common use case where a Pod runs a single container. This is similar to running a single Docker container, but with the added benefits of Kubernetes orchestration, such as automated scaling, self-healing, and easy updates.

2. **Multi-Container Pod**: In some cases, a Pod may run multiple containers that need to work closely together. These containers share the same network namespace, which allows them to communicate with each other using `localhost`. They can also share storage volumes, making it easy to manage data that needs to be accessed by multiple containers. This pattern is useful for sidecar containers, where one container enhances or extends the functionality of the main application container. 

## Pod Lifecycle
Pods follow a defined lifecycle, starting in the Pending phase, moving through Running if at least one of its primary containers starts OK, and then through either the Succeeded or Failed phases depending on whether any container in the Pod terminated in failure.

### Pod lifetime
- assigning a Pod to a specific node is called binding
- process of selecting which node to use is called scheduling. 
- Once a Pod has been scheduled and is bound to a node, Kubernetes tries to run that Pod on the node.
- The Pod runs on that node until it stops, or until the Pod is terminated; 

### Pod phases
A Pod's status field is a PodStatus object, which has a phase field. The possible values for the phase field are:

1. **Pending**: The Pod has been accepted by the Kubernetes system, but one or more of the container images has not been created. This includes time before being scheduled as well as time spent downloading images over the network.

2. **Running**: The Pod has been bound to a node, and all of the containers have been created. At least one container is still running, or is in the process of starting or restarting.

3. **Succeeded**: All containers in the Pod have terminated in success, and will not be restarted.

4. **Failed**: All containers in the Pod have terminated, and at least one container has terminated in failure. That is, the container either exited with non-zero status or was terminated by the system.

5. **Unknown**: For some reason the state of the Pod could not be obtained, typically due to an error in communicating with the host of the Pod.

### Container states
Each container in a Pod has its own state, which can be:
- **Waiting**: The container is waiting to start
- **Running**: The container is running
- **Terminated**: The container has completed execution

### Pod conditions
A Pod has a PodStatus, which has an array of PodConditions through which the Pod has or has not passed:
- **PodScheduled**: the Pod has been scheduled to a node
- **ContainersReady**: all containers in the Pod are ready
- **Initialized**: all init containers have completed successfully
- **Ready**: the Pod is able to serve requests and should be added to the load balancing pools of all matching Services

### Relationship Between Container States and Pod Phases

The phase of a Pod is determined by the collective states of its containers. Here's how container states influence Pod phases:

1. **Pending Phase**:
   - All containers are in the **Waiting** state
   - This occurs during image pulling, container creation, or scheduling

2. **Running Phase**:
   - At least one container is in the **Running** state
   - Other containers may be in **Waiting** or **Terminated** states
   - The Pod is actively executing its workload

3. **Succeeded Phase**:
   - All containers are in the **Terminated** state
   - All containers exited with a status code of 0 (success)
   - No containers are configured to restart

4. **Failed Phase**:
   - All containers are in the **Terminated** state
   - At least one container exited with a non-zero status code
   - All restart attempts have been exhausted

5. **Unknown Phase**:
   - Container states cannot be determined
   - Usually indicates a communication problem with the node

#### Key Points:
- A Pod's phase is an aggregate of its containers' states
- The phase represents the overall state of the Pod
- Container states provide more granular information about individual containers
- Pod conditions (like Ready) are additional status indicators that complement the phase information
- The phase is particularly useful for monitoring and automation purposes

## Example: Basic Pod YAML

Below is an example of a basic Pod YAML configuration:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  labels:
    app: example
spec:
  containers:
  - name: example-container
    image: nginx:latest
    ports:
    - containerPort: 80
```

### Description
- **apiVersion**: Specifies the version of the Kubernetes API to use. Here, `v1` is used for core resources like Pods.
- **kind**: Defines the type of Kubernetes resource. In this case, it's a `Pod`.
- **metadata**: Contains data that helps uniquely identify the Pod, such as its `name` and `labels`.
  - **name**: The name of the Pod, which must be unique within a namespace.
  - **labels**: Key-value pairs used to organize and select resources.
- **spec**: Describes the desired state of the Pod.
  - **containers**: A list of containers that will run in the Pod.
    - **name**: The name of the container.
    - **image**: The container image to run, specified as `name:tag`.
    - **ports**: A list of ports to expose from the container. Here, port 80 is exposed for HTTP traffic.

This YAML file defines a simple Pod running a single Nginx container, which listens on port 80. It's a basic example to demonstrate how to define a Pod in Kubernetes.

example