Couchbase Exporter
==================

Expose metrics from *Couchbase Community 4.5.1* cluster for consumption by Prometheus.

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

Diclaimer
---------

This project is intented for me to learn Go. It is therefore not suited for real world use.
But if you want to help and make it a valuable project, feel free to contribute!

Author Information
------------------

Blake Lead