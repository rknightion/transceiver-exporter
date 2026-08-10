---
title: Deployment
description: Running transceiver-exporter as a Kubernetes DaemonSet or a systemd unit, and the Prometheus scrape config for either.
---

# Deployment

`transceiver-exporter` reads NIC state directly from the host it runs on, so it is a
**one-instance-per-host** tool: it belongs on every host whose transceivers you want visibility
into, not behind a load balancer or as a horizontally-scaled service. The two natural shapes are a
Kubernetes DaemonSet (one pod per node) and a systemd unit (one process per bare-metal or VM
host). Neither pattern is shipped in this repository — both are derived from the same
`--network host` + `CAP_NET_ADMIN` requirement covered in [Permissions](permissions.md); adapt
them to your environment.

## Kubernetes: DaemonSet

A DaemonSet needs the pod on the **host network namespace** (the same reason `docker run` needs
`--network host`) and the capability grant from [Permissions](permissions.md):

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: transceiver-exporter
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: transceiver-exporter
  template:
    metadata:
      labels:
        app: transceiver-exporter
    spec:
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      containers:
        - name: transceiver-exporter
          image: ghcr.io/rknightion/transceiver-exporter:latest
          args:
            - -web.listen-address=[::]:9458
          ports:
            - containerPort: 9458
              hostPort: 9458
          securityContext:
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["ALL"]
              add: ["NET_ADMIN"]
```

`hostNetwork: true` is what gives the pod visibility into the node's physical NICs, exactly like
`--network host` does for `docker run`. Dropping all capabilities and adding back only
`NET_ADMIN` mirrors the container posture in [Installation](installation.md) — the pod still runs
as root (the image sets no non-root user; see [Permissions](permissions.md)), bounded by the
capability set, not the UID.

### Scrape config

With `hostNetwork: true` the pod is reachable on the node's own address, so a node-based discovery
role is the natural fit:

```yaml
scrape_configs:
  - job_name: transceiver-exporter
    kubernetes_sd_configs:
      - role: pod
        selectors:
          - role: pod
            label: "app=transceiver-exporter"
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_host_ip]
        target_label: __address__
        replacement: "${1}:9458"
      - source_labels: [__meta_kubernetes_pod_node_name]
        target_label: node
```

A plain node-address static config works just as well if you already enumerate nodes elsewhere:

```yaml
scrape_configs:
  - job_name: transceiver-exporter
    static_configs:
      - targets: ["node-a:9458", "node-b:9458"]
```

## systemd

For bare-metal or VM hosts, run the released binary under systemd with the ambient-capability
grant covered in [Permissions](permissions.md#systemd-the-minimum-grant-without-running-as-root):

```ini
[Unit]
Description=transceiver-exporter
After=network.target

[Service]
ExecStart=/usr/local/bin/transceiver-exporter -web.listen-address=[::]:9458
User=transceiver-exporter
AmbientCapabilities=CAP_NET_ADMIN
CapabilityBoundingSet=CAP_NET_ADMIN
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo useradd --system --no-create-home transceiver-exporter
sudo cp transceiver-exporter /usr/local/bin/
sudo systemctl daemon-reload
sudo systemctl enable --now transceiver-exporter
```

### Scrape config

```yaml
scrape_configs:
  - job_name: transceiver-exporter
    static_configs:
      - targets: ["<host>:9458"]
```

## Fleet-wide filtering

On a switch or dense multi-NIC host, `-collector.interface-features.enable=false` and an
`-exclude.interfaces-regex` for management/loopback-style virtual interfaces keep cardinality down
across a fleet of these — see [Configuration](configuration.md).
