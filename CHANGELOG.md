Changelog
=========

0.1.0
-----

First working version.

The exporter exposes the following metrics:

- **cb_up**: State of last cluster scrape
- **cb_node_status**: Status of couchbase node
- **cb_node_cluster_membership**: Status of node cluster membership
- **cb_node_cpu_utilization_rate**: CPU utilization rate
- **cb_node_ram_usage_bytes**: RAM used per node in bytes
- **cb_cluster_ram_total_bytes**: Total RAM in the cluster