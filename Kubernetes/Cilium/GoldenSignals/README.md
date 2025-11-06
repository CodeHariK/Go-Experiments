# [Golden Signals with Hubble and Grafana](https://isovalent.com/labs/hubble-grafana-golden-signals/)

* Grafana Dashboard https://github.com/cilium/cilium/blob/main/install/kubernetes/cilium/files/hubble/dashboards/hubble-l7-http-metrics-by-workload.json

⏱️ Latency
Most often represented as response time in milliseconds (ms) at the application layer.

Application response time is affected by latency across all of the core system resources including network, storage, processor (CPU), and memory.

Latency at the application layer also needs to be correlated to latency and resource usage that may be happening internally within the application processes, between pods/services, across the network/mesh etc.

![latency](https://play.instruqt.com/assets/tracks/a9amhryys7kc/99a4e233a51466b3a18c193114ab2ded/assets/latency.png)

🚦 Throughput
Sometimes referred to as traffic, throughput is the volume and types of requests that are being sent and received by services and applications from within and from outside a Kubernetes environment.

Throughput metrics include examples like web requests, API calls, and is described as the demand commonly represented as the number of requests per second.

It should be measured across all layers to identify requests to and from services, and also which I/O is going further down to the node.

![throughput](https://play.instruqt.com/assets/tracks/a9amhryys7kc/6fd91afe7ba6f455b79539b5d9d9b67e/assets/throughput.png)

⚠️ Errors
The number of requests (traffic) which are failing, often represented either in absolute numbers or as the percentage of requests with errors versus the total number of requests.

There may be errors that happen due to application issues, possible misconfiguration, and some errors happen as defined by policy.

Policy-driven error may indicate accidental misconfiguration or potentially a malicious process.

💥 Saturation
The overall utilization of resources including CPU (capacity, quotas, throttling), memory (capacity, allocation), storage (capacity, allocation, and I/O throughput), and network.

Some resources saturate linearly (e.g. storage capacity) while others (memory, CPU, and network) fluctuate much more with the ephemeral nature of containerized applications.

Network saturation is a great example of the complexity of monitoring Kubernetes because there is node networking, service-to-service network throughput, and once a service mesh is in place, there are more paths, and potentially more bottlenecks that can be saturated.

![saturation](https://play.instruqt.com/assets/tracks/a9amhryys7kc/2a6cc8be0809382939d089556ce17ad3/assets/golden-signals-blog-post-ressources-consumption.png)

## ⎈ Golden Signals and Kubernetes Observability
Observability is about being able to ask questions of the system rather than just piling up monitoring data to attempt to correlate it.

What are the questions that we can ask the system in order to understand the current state? For example,

What is causing slowness for application X?
Which services are affecting errors in my front-end application?
What is causing congestion on my cluster nodes?
Why is my application unable to contact my messaging service in both clouds?

🕵️‍♀️ Why is Observability in Kubernetes a Multi-Dimensional Challenge?
Our 4 golden signals for observing Kubernetes are especially interesting (and challenging) because each has its own measurement of health. An aggregate combination of golden signals defines the overall system health. Both visually and mathematically, this is a multi-dimensional challenge.

How often do SRE and container Ops teams get asked “what’s going on with Application X?” without having a specific application monitoring trigger or alert? Or even when certain alerts are appearing but aren’t the singular reason an application is negatively affected.

There could be a combination of latency, utilization, and errors that have to be correlated to the root cause.

```yaml
operator:
  # only 1 replica needed on a single node setup
  replicas: 1
  prometheus:
    enabled: true
    serviceMonitor:
      enabled: true

hubble:
  relay:
    # enable relay in 02
    # enabled: true
    service:
      type: NodePort
    prometheus:
      enabled: true
      serviceMonitor:
        enabled: true

  metrics:
    serviceMonitor:
      enabled: true
    enableOpenMetrics: true
    enabled:
      - dns
      - drop
      - tcp
      - icmp
      - "flow:sourceContext=workload-name|reserved-identity;destinationContext=workload-name|reserved-identity"
      - "kafka:labelsContext=source_namespace,source_workload,destination_namespace,destination_workload,traffic_direction;sourceContext=workload-name|reserved-identity;destinationContext=workload-name|reserved-identity"
      - "httpV2:exemplars=true;labelsContext=source_ip,source_namespace,source_workload,destination_ip,destination_namespace,destination_workload,traffic_direction;sourceContext=workload-name|reserved-identity;destinationContext=workload-name|reserved-identity"
    dashboards:
      enabled: true
      namespace: monitoring
      annotations:
        grafana_folder: "Hubble"

  ui:
    # enable UI in 02
    # enabled: true
    service:
      type: NodePort

prometheus:
  enabled: true
  serviceMonitor:
    enabled: true
```

In particular, we've enabled the httpv2 Hubble metrics with:

httpV2:exemplars=true;labelsContext=source_ip,source_namespace,source_workload,destination_ip,destination_namespace,destination_workload,traffic_direction;sourceContext=workload-name|reserved-identity;destinationContext=workload-name|reserved-identity
This is a rather long line. Let's split the options to understand what they do:

examplars=true will let us display OpenTelemetry trace points from application traces as an overlay on the Grafana graphs
labelsContext is set to add extra labels to metrics including source/destination IPs, source/destination namespaces, source/destination workloads, as well as traffic direction (ingress or egress)
sourceContext sets how the source label is built, in this case using the workload name when possible, or a reserved identity (e.g. world) otherwise
destinationContext does the same for destinations
