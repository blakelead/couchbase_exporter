Couchbase Exporter
==================

Expose metrics from *Couchbase Community 4.5.1* cluster for consumption by Prometheus.

> WIP: each release will introduce new metrics and probably lots of breaking changes.

Getting Started
---------------

Run from terminal:

```bash
./couchbase_exporter [flags]
```

Available flags:

| argument            | description                                | default               |
|---------------------|--------------------------------------------|-----------------------|
| -web.listen-address | The address to listen on for HTTP requests | :9191                 |
| -web.telemetry-path | Path under which to expose metrics         | /metrics              |
| -db.url             | The address of Couchbase cluster           | http://localhost:8091 |
| -db.user            | The administrator username                 | admin                 |
| -db.pwd             | The administrator password                 | password              |
| -help               | Command line help                          |                       |

Metrics
-------

| name                         | description                                |
|------------------------------|--------------------------------------------|
| cb_up                        | State of last cluster scrape               |
| cb_node_status               | Status of couchbase node                   |
| cb_node_cluster_membership   | Status of node cluster membership          |
| cb_node_cpu_utilization_rate | CPU utilization rate                       |
| cb_node_ram_usage_bytes      | RAM used per node in bytes                 |
| cb_cluster_ram_total_bytes   | Total RAM in the cluster                   |

Author Information
------------------

Blake Lead