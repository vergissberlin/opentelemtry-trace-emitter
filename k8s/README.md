# OTel-Trace-Emitter with Kubernetes

This example shows how to deploy the OpenTelemetry Collector as a Kubernetes Deployment.

## Prerequisites

### Prepare - Local cluster with multiple nodes

To test it locally with a realistic Kubernetes cluster, you can use [Kind](https://kind.sigs.k8s.io/).
We can spawn a Kubernetes cluster with multiple nodes using Kind.

1. Create a Kind cluster with the configuration file `kind-cluster.yaml`

```shell
   kind create cluster --name thinkport  --config kind-cluster.yaml
   ```

2. Verify the cluster

```shell
   kubectl cluster-info --context kind-thinkport
   kubectl get nodes
   ```

### Deploy the OpenTelemetry Collector as a Kubernetes DaemonSet

1. Create a namespace for the deployment

```bash
   kubectl create namespace monitoring
   ```

2. Apply the ConfigMap

```shell
   kubectl apply -f collector/otel-configmap.yaml
   ```

3. Apply the DaemonSet

```shell
   kubectl apply -f collector/otel-daemonset.yaml
   ```

4. Verify the DaemonSet

```shell
    kubectl get daemonsets -n monitoring
    kubectl get pods -n monitoring
    ```

5. Delete the DaemonSet
   

```shell
   kubectl delete daemonset otel-collector-daemonset -n monitoring
   ```

## Emitter

This example uses the OpenTelemetry Collector to receive traces and metrics. To send traces and metrics to the OpenTelemetry Collector, you can use the [OpenTelemetry Node.js example](

1. Apply the deployment manifest

```shell
    kubectl delete deployment otel-emitter -n monitoring
    kubectl apply -f emitter/emitter-deployment.yaml -n monitoring
    ```

2. Verify the deployment
    

```shell
    kubectl get deployments -n monitoring
    kubectl get pods -n monitoring
    ```

3. Get the logs
    

```shell
    kubectl logs -l app=otel-emitter -n monitoring
    ```

   
4. Verify the traces and metrics in the OpenTelemetry Collector
    

```shell
    kubectl port-forward svc/otel-collector-daemonset 4317:4317 -n monitoring
    ```

   
5. Delete the deployment
    

```shell
    kubectl delete deployment otel-emitter -n monitoring
    ```

   
