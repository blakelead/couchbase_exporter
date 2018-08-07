Couchbase Exporter
==================

[![Build Status](https://travis-ci.org/blakelead/couchbase_exporter.svg?branch=master)](https://travis-ci.org/blakelead/couchbase_exporter)
[![Coverage Status](https://coveralls.io/repos/github/blakelead/couchbase_exporter/badge.svg?branch=master)](https://coveralls.io/github/blakelead/couchbase_exporter?branch=master)
[![Software License](https://img.shields.io/badge/license-MIT-green.svg)](/LICENSE.txt)

Expose metrics from *Couchbase Community 4.5.1* cluster for consumption by Prometheus.

> WIP: each release will introduce new metrics and probably lots of breaking changes.

Getting Started
---------------

Run from command line:

```bash
$ ./couchbase_exporter [flags]
```

Available flags:

| argument            | description                                | default               |
|---------------------|--------------------------------------------|-----------------------|
| -web.listen-address | The address to listen on for HTTP requests | :9191                 |
| -web.telemetry-path | Path under which to expose metrics         | /metrics              |
| -db.url             | The address of Couchbase cluster           | http://localhost:8091 |
| -db.user            | The administrator username                 | admin                 |
| -db.pwd             | The administrator password                 | password              |
| -log.level          | Log level: info, debug, warn, error, fatal | info                  |
| -log.format         | Log format: text, json                     | text                  |
| -help               | Command line help                          |                       |

Metrics
-------

| name                                      | description                                     |
|-------------------------------------------|-------------------------------------------------|
| cb_up                                     | State of last cluster scrape                    |
| cb_node_status                            | Status of couchbase node                        |
| cb_node_cluster_membership                | Status of node cluster membership               |
| cb_node_cpu_utilization_rate              | CPU utilization rate                            |
| cb_node_ram_usage_bytes                   | RAM used per node in bytes                      |
| cb_cluster_ram_total_bytes                | Total RAM in the cluster                        |
| cb_cluster_ram_used_bytes                 | Used RAM in the cluster                         |
| cb_cluster_ram_used_by_data_bytes         | Used RAM by data in the cluster                 |
| cb_cluster_ram_quota_total_bytes          | Total quota RAM in the cluster                  |
| cb_cluster_ram_quota_total_per_node_bytes | Total quota RAM per node in the cluster         |
| cb_cluster_ram_quota_used_bytes           | Used quota RAM in the cluster                   |
| cb_cluster_ram_quota_used_per_node_bytes  | Used quota RAM per node in the cluster          |
| cb_cluster_disk_total_bytes               | Total disk in the cluster                       |
| cb_cluster_disk_quota_total_bytes         | Disk quota in the cluster                       |
| cb_cluster_disk_used_bytes                | Used disk in the cluster                        |
| cb_cluster_disk_used_by_data_bytes        | Disk used by data in the cluster                |
| cb_cluster_disk_free_bytes                | Free disk in the cluster                        |
| cb_cluster_fts_ram_quota_bytes            | RAM quota for Full text search bucket           |
| cb_cluster_index_ram_quota_bytes          | RAM quota for Index bucket                      |
| cb_cluster_data_ram_quota_bytes           | RAM quota for Data bucket                       |
| cb_cluster_rebalance_status               | Occurrence of rebalancing in the cluster        |
| cb_cluster_max_bucket_count               | Maximum number of buckets                       |
| cb_cluster_failover_node_count            | Number of failovers since cluster is up         |
| cb_cluster_rebalance_success_count        | Number of rebalance success since cluster is up |
| cb_cluster_rebalance_start_count          | Number of rebalance start since cluster is up   |
| cb_cluster_rebalance_fail_count           | Number of rebalance failure since cluster is up |

### Node metrics

|                   name                    |                   description                   |
| ----------------------------------------- | ----------------------------------------------- |
| node_ram_total_bytes                      | Node total RAM                                  |
| node_ram_usage_bytes                      | Node used RAM                                   |
| node_ram_used_by_data_bytes               | Node RAM used by data                           |
| node_ram_quota_total_bytes                | Node RAM quota total                            |
| node_ram_quota_used_bytes                 | Node RAM quota used                             |
| node_disk_total_bytes                     | Node Total disk                                 |
| node_disk_quota_total_bytes               | Node disk quota total                           |
| node_disk_used_bytes                      | Node used disk                                  |
| node_disk_quota_used_bytes                | Node disk quota total                           |
| node_disk_free_bytes                      | Node free disk                                  |
| node_cpu_utilization_rate                 | CPU utilization rate                            |
| node_swap_total_bytes                     | Node total swap                                 |
| node_swap_used_bytes                      | Node used swap                                  |
| node_stats_cmd_get                        | Node stats: cmd_get                             |
| node_stats_couch_docs_actual_disk_size    | Node stats: couch_docs_actual_disk_size         |
| node_stats_couch_docs_data_size           | Node stats: couch_docs_data_size                |
| node_stats_couch_spatial_data_size        | Node stats: couch_spatial_data_size             |
| node_stats_couch_spatial_disk_size        | Node stats: couch_spatial_disk_size             |
| node_stats_couch_views_actual_disk_size   | Node stats: couch_views_actual_disk_siz         |
| node_stats_couch_views_data_size          | Node stats: couch_views_data_size               |
| node_stats_curr_items                     | Node stats: curr_items                          |
| node_stats_curr_items_tot                 | Node stats: curr_items_tot                      |
| node_stats_ep_bg_fetched                  | Node stats: ep_bg_fetched                       |
| node_stats_get_hits                       | Node stats: get_hits                            |
| node_stats_mem_used                       | Node stats: mem_used                            |
| node_stats_ops                            | Node stats: ops                                 |
| node_stats_vb_replica_curr_items          | Node stats: vb_replica_curr_items               |
| node_uptime_seconds                       | Node uptime                                     |
| node_cluster_membership                   | Status of node cluster membership               |
| node_status                               | Status of couchbase node                        |
| node_fts_ram_quota_bytes                  | Node quota for Full text search bucket          |
| node_index_ram_quota_bytes                | Node quota for Index bucket                     |
| node_data_ram_quota_bytes                 | Node quota for Data bucket                      |

Docker
------

Get the latest image from Docker Hub:

```bash
$ docker pull blakelead/couchbase-exporter:latest
```

Systemd
-------

You can use `exporter.service` to execute **couchbase_exporter** with systemd.

```bash
$ sudo mv exporter.service /etc/systemd/system/couchbase-exporter.service
$ sudo systemctl enable couchbase-exporter.service
$ sudo systemctl start couchbase-exporter.service
```

Author Information
------------------

Adel Abdelhak