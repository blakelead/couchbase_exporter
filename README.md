Couchbase Exporter
==================

[![Release](https://img.shields.io/badge/release-0.1.0-blue.svg)](https://github.com/blakelead/couchbase_exporter/releases/tag/0.1.0)
[![Build Status](https://travis-ci.com/blakelead/couchbase_exporter.svg?branch=master)](https://travis-ci.org/blakelead/couchbase_exporter)
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

You can either use command-line flags, or environment variables to pass custom parameters. Environment variables take precedence over flags.

Available flags and equivalent environment variable:

|      argument       |    environment variable    |              description               |        default        |
| ------------------- | -------------------------- | -------------------------------------- | --------------------- |
| -web.listen-address | CB_EXPORTER_LISTEN_ADDR    | Address to listen on for HTTP requests | :9191                 |
| -web.telemetry-path | CB_EXPORTER_TELEMETRY_PATH | Path under which to expose metrics     | /metrics              |
| -db.uri             | CB_EXPORTER_DB_URI         | Address of Couchbase cluster           | http://127.0.0.1:8091 |
| -db.user            | CB_EXPORTER_DB_USER        | Administrator username                 | admin                 |
| -db.pwd             | CB_EXPORTER_DB_PASSWORD    | Administrator password                 | password              |
| -log.level          | CB_EXPORTER_LOG_LEVEL      | Log level: info,debug,warn,error,fatal | info                  |
| -log.format         | CB_EXPORTER_LOG_FORMAT     | Log format: text, json                 | text                  |
| -scrape.cluster     | CB_EXPORTER_SCRAPE_CLUSTER        | If false, wont scrape cluster metrics  | true                  |
| -scrape.node        | CB_EXPORTER_SCRAPE_NODE           | If false, wont scrape node metrics     | true                  |
| -scrape.bucket      | CB_EXPORTER_SCRAPE_BUCKET         | If false, wont scrape bucket metrics   | true                  |
| -help               |                            | Command line help                      |                       |

Metrics
-------

### Cluster metrics

|                name                |                    description                     |
| ---------------------------------- | -------------------------------------------------- |
| cb_cluster_ram_total_bytes         | Total memory available to the cluster              |
| cb_cluster_ram_used_bytes          | Memory used by the cluster                         |
| cb_cluster_ram_used_by_data_bytes  | Memory used by the data in the cluster             |
| cb_cluster_ram_quota_total_bytes   | Total memory allocated to Couchbase in the cluster |
| cb_cluster_ram_quota_used_bytes    | Memory quota used by the cluster                   |
| cb_cluster_disk_total_bytes        | Total disk space available to the cluster          |
| cb_cluster_disk_used_bytes         | Disk space used by the cluster                     |
| cb_cluster_disk_quota_total_bytes  | Disk space quota for the cluster                   |
| cb_cluster_disk_used_by_data_bytes | Disk space used by the data in the cluster         |
| cb_cluster_disk_free_bytes         | Free disk space in the cluster                     |
| cb_cluster_fts_ram_quota_bytes     | Memory quota allocated to full text search buckets |
| cb_cluster_index_ram_quota_bytes   | Memory quota allocated to Index buckets            |
| cb_cluster_data_ram_quota_bytes    | Memory quota allocated to Data buckets             |
| cb_cluster_rebalance_status        | Rebalancing status                                 |
| cb_cluster_max_bucket_count        | Maximum number of buckets allowed                  |
| cb_cluster_failover_node_count     | Number of failovers since cluster is up            |
| cb_cluster_rebalance_success_count | Number of rebalance successes since cluster is up  |
| cb_cluster_rebalance_start_count   | Number of rebalance starts since cluster is up     |
| cb_cluster_rebalance_fail_count    | Number of rebalance fails since cluster is up      |

### Node metrics

|                    name                    |                    description                     |
| ------------------------------------------ | -------------------------------------------------- |
| cb_node_service_up                         | Couchbase service healthcheck                      |
| cb_node_ram_total_bytes                    | Total memory available to the node                 |
| cb_node_ram_usage_bytes                    | Memory used by the node                            |
| cb_node_ram_used_by_data_bytes             | Memory used by data in the node                    |
| cb_node_ram_quota_total_bytes              | Memory quota allocated to the node                 |
| cb_node_ram_quota_used_bytes               | Memory quota used by the node                      |
| cb_node_disk_total_bytes                   | Total disk space available to the node             |
| cb_node_disk_quota_total_bytes             | Disk space quota for the node                      |
| cb_node_disk_used_bytes                    | Disk space used by the node                        |
| cb_node_disk_quota_used_bytes              | Disk space quota used by the node                  |
| cb_node_disk_free_bytes                    | Free disk space in the node                        |
| cb_node_cpu_utilization_rate               | CPU utilization rate in percent                    |
| cb_node_swap_total_bytes                   | Total swap space allocated to the node             |
| cb_node_swap_used_bytes                    | Amount of swap space used by the node              |
| cb_node_stats_cmd_get                      | Number of get commands                             |
| cb_node_stats_couch_docs_actual_disk_size  | Disk space used by Couchbase documents             |
| cb_node_stats_couch_docs_data_size         | Couchbase documents data size in the node          |
| cb_node_stats_couch_spatial_data_size      | Data size for Couchbase spatial views              |
| cb_node_stats_couch_spatial_disk_size      | Disk space used by Couchbase spatial views         |
| cb_node_stats_couch_views_actual_disk_size | Disk space used by Couchbase views                 |
| cb_node_stats_couch_views_data_size        | Data size for Couchbase views                      |
| cb_node_stats_curr_items                   | Number of current items                            |
| cb_node_stats_curr_items_tot               | Total number of items in the node                  |
| cb_node_stats_ep_bg_fetched                | Number of background disk fetches                  |
| cb_node_stats_get_hits                     | Number of get hits                                 |
| cb_node_stats_mem_used                     | Memory used by the node                            |
| cb_node_stats_ops                          | Number of operations performed in the node         |
| cb_node_stats_vb_replica_curr_items        | Number of replicas in current items                |
| cb_node_uptime_seconds                     | Node uptime                                        |
| cb_node_cluster_membership                 | Cluster membership                                 |
| cb_node_status                             | Status of couchbase node                           |
| cb_node_fts_ram_quota_bytes                | Memory quota allocated to full text search buckets |
| cb_node_index_ram_quota_bytes              | Memory quota allocated to index buckets            |
| cb_node_data_ram_quota_bytes               | Memory quota allocated to data buckets             |

### Bucket metrics

|               name               |             description              |
| -------------------------------- | ------------------------------------ |
| cb_bucket_proxy_port             | Bucket proxy port                    |
| cb_bucket_replica_index          | Replica index for the bucket         |
| cb_bucket_replica_number         | Number of replicas for the bucket    |
| cb_bucket_threads_number         | Bucket thread number                 |
| cb_bucket_ram_quota_bytes        | Memory used by the bucket            |
| cb_bucket_raw_ram_quota_bytes    | Raw memory used by the bucket        |
| cb_bucket_ram_quota_percent_used | Memory used by the bucket in percent |
| cb_bucket_ops_per_second         | Number of operations per second      |
| cb_bucket_disk_fetches           | Disk fetches for the bucket          |
| cb_bucket_item_count             | Number of items in the bucket        |
| cb_bucket_disk_used_bytes        | Disk used by the bucket              |
| cb_bucket_data_used_bytes        | Data loaded in memory                |
| cb_bucket_ram_used_bytes         | Bucket RAM used                      |

Docker
------

Get the latest image from Docker Hub:

```bash
$ docker pull blakelead/couchbase-exporter:latest
```

and use it with environment variables like in the following example:

```bash
$ docker run --name cbexporter -p 9191:9191 -e CB_EXPORTER_DB_USER=admin -e CB_EXPORTER_DB_PASSWORD=complicatedpassword blakelead/couchbase-exporter:latest
```

Systemd
-------

You can use `exporter.service` to execute **couchbase_exporter** with systemd.

```bash
$ sudo mv exporter.service /etc/systemd/system/couchbase-exporter.service
$ sudo systemctl enable couchbase-exporter.service
$ sudo systemctl start couchbase-exporter.service
```

Todo
-----

- Unit tests
- XDCR metrics
- Couchbase 5 compatibility
- Advanced bucket metrics
- Cleaner code

Author Information
------------------

Adel Abdelhak