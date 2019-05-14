# Couchbase Exporter

[![Build Status](https://travis-ci.com/blakelead/couchbase_exporter.svg?branch=master)](https://travis-ci.org/blakelead/couchbase_exporter)
[![Coverage Status](https://coveralls.io/repos/github/blakelead/couchbase_exporter/badge.svg?branch=master)](https://coveralls.io/github/blakelead/couchbase_exporter?branch=master)
[![Software License](https://img.shields.io/badge/license-MIT-green.svg)](/LICENSE.txt)

Expose metrics from Couchbase cluster for consumption by Prometheus.

## Getting Started

Run from command line:

```bash
./couchbase_exporter [flags]
```

The exporter supports various configuration ways: command line arguments takes precedence over environment variables that take precedence over configuration file.

Configuration file can be written in json or yaml, must be named `config.json` or `config.yml`, and must be placed in the same directory that the exporter binary. You can find complete examples of configuation files in the sources (directory `examples`).

As for available flags and equivalent environment variables, here is a list:

|    environment variable    |      argument       |              description               |        default        |
| -------------------------- | ------------------- | -------------------------------------- | --------------------- |
| CB_EXPORTER_LISTEN_ADDR    | -web.listen-address | Address to listen on for HTTP requests | :9191                 |
| CB_EXPORTER_TELEMETRY_PATH | -web.telemetry-path | Path under which to expose metrics     | /metrics              |
| CB_EXPORTER_SERVER_TIMEOUT | -web.timeout        | Server read timeout in seconds         | 10s                   |
| CB_EXPORTER_DB_URI         | -db.uri             | Address of Couchbase cluster           | https://127.0.0.1:18091 |
| CB_EXPORTER_DB_URI2        | -db2.uri            | Secondary connection to Couchbase clust| https://127.0.0.1:18091
| CB_EXPORTER_DB_TIMEOUT     | -db.timeout         | Couchbase client timeout in seconds    | 10s                   |
| CB_EXPORTER_DB_USER        | *not allowed*       | Administrator username                 |                       |
| CB_EXPORTER_DB_PASSWORD    | *not allowed*       | Administrator password                 |                       |
| CB_EXPORTER_LOG_LEVEL      | -log.level          | Log level: info,debug,warn,error,fatal | error                 |
| CB_EXPORTER_LOG_FORMAT     | -log.format         | Log format: text, json                 | text                  |
| CB_EXPORTER_SCRAPE_CLUSTER | -scrape.cluster     | If false, wont scrape cluster metrics  | true                  |
| CB_EXPORTER_SCRAPE_NODE    | -scrape.node        | If false, wont scrape node metrics     | true                  |
| CB_EXPORTER_SCRAPE_BUCKET  | -scrape.bucket      | If false, wont scrape bucket metrics   | true                  |
| CB_EXPORTER_SCRAPE_XDCR    | -scrape.xdcr        | If false, wont scrape xdcr metrics     | false                 |
| CB_EXPORTER_TLS_SETTING    | -tls.setting        | If true, will ignore Self Signed Certs | false                 |
|                            | -help               | Command line help                      |                       |

> Important: for security reasons credentials cannot be set with command line arguments.

## Metrics

### Cluster metrics

|                name                |                     description                     |
| ---------------------------------- | --------------------------------------------------- |
| cb_cluster_ram_total_bytes         | Total memory available to the cluster               |
| cb_cluster_ram_used_bytes          | Memory used by the cluster                          |
| cb_cluster_ram_used_by_data_bytes  | Memory used by the data in the cluster              |
| cb_cluster_ram_quota_total_bytes   | Total memory allocated to Couchbase in the cluster  |
| cb_cluster_ram_quota_used_bytes    | Memory quota used by the cluster                    |
| cb_cluster_disk_total_bytes        | Total disk space available to the cluster           |
| cb_cluster_disk_used_bytes         | Disk space used by the cluster                      |
| cb_cluster_disk_quota_total_bytes  | Disk space quota for the cluster                    |
| cb_cluster_disk_used_by_data_bytes | Disk space used by the data in the cluster          |
| cb_cluster_disk_free_bytes         | Free disk space in the cluster                      |
| cb_cluster_fts_ram_quota_bytes     | Memory quota allocated to full text search buckets  |
| cb_cluster_index_ram_quota_bytes   | Memory quota allocated to Index buckets             |
| cb_cluster_data_ram_quota_bytes    | Memory quota allocated to Data buckets              |
| cb_cluster_rebalance_status        | Rebalancing status                                  |
| cb_cluster_max_bucket_count        | Maximum number of buckets allowed                   |
| cb_cluster_failover_node_count     | Number of failovers since cluster is up             |
| cb_cluster_rebalance_success_count | Number of rebalance successes since cluster is up   |
| cb_cluster_rebalance_start_count   | Number of rebalance starts since cluster is up      |
| cb_cluster_rebalance_fail_count    | Number of rebalance fails since cluster is up       |
| cb_cluster_balanced                | Status of cluster balance (in 5.1.1) |

### Node metrics

|                      name                       |                        description                         |
| ----------------------------------------------- | ---------------------------------------------------------- |
| cb_node_service_up                              | Couchbase service healthcheck                              |
| cb_node_ram_total_bytes                         | Total memory available to the node                         |
| cb_node_ram_usage_bytes                         | Memory used by the node                                    |
| cb_node_ram_used_by_data_bytes                  | Memory used by data in the node                            |
| cb_node_ram_quota_total_bytes                   | Memory quota allocated to the node                         |
| cb_node_ram_quota_used_bytes                    | Memory quota used by the node                              |
| cb_node_disk_total_bytes                        | Total disk space available to the node                     |
| cb_node_disk_quota_total_bytes                  | Disk space quota for the node                              |
| cb_node_disk_used_bytes                         | Disk space used by the node                                |
| cb_node_disk_used_by_data_bytes                 | Disk space used by data in the node                        |
| cb_node_disk_free_bytes                         | Free disk space in the node                                |
| cb_node_cpu_utilization_rate                    | CPU utilization rate in percent                            |
| cb_node_swap_total_bytes                        | Total swap space allocated to the node                     |
| cb_node_swap_used_bytes                         | Amount of swap space used by the node                      |
| cb_node_stats_cmd_get                           | Number of get commands                                     |
| cb_node_stats_couch_docs_actual_disk_size       | Disk space used by Couchbase documents                     |
| cb_node_stats_couch_docs_data_size              | Couchbase documents data size in the node                  |
| cb_node_stats_couch_spatial_data_size           | Data size for Couchbase spatial views                      |
| cb_node_stats_couch_spatial_disk_size           | Disk space used by Couchbase spatial views                 |
| cb_node_stats_couch_views_actual_disk_size      | Disk space used by Couchbase views                         |
| cb_node_stats_couch_views_data_size             | Data size for Couchbase views                              |
| cb_node_stats_curr_items                        | Number of current items                                    |
| cb_node_stats_curr_items_tot                    | Total number of items in the node                          |
| cb_node_stats_ep_bg_fetched                     | Number of background disk fetches                          |
| cb_node_stats_get_hits                          | Number of get hits                                         |
| cb_node_stats_mem_used                          | Memory used by the node                                    |
| cb_node_stats_ops                               | Number of operations performed in the node                 |
| cb_node_stats_vb_replica_curr_items             | Number of replicas in current items                        |
| cb_node_stats_vb_active_num_non_resident_number | Number of non-resident items in active vbuckets (in 5.1.1) |
| cb_node_uptime_seconds                          | Node uptime                                                |
| cb_node_cluster_membership                      | Cluster membership                                         |
| cb_node_status                                  | Status of couchbase node                                   |
| cb_node_fts_ram_quota_bytes                     | Memory quota allocated to full text search buckets         |
| cb_node_index_ram_quota_bytes                   | Memory quota allocated to index buckets                    |
| cb_node_data_ram_quota_bytes                    | Memory quota allocated to data buckets                     |

### Bucket metrics

|                        name                        |                                               description                                               |
| -------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| cb_bucket_ram_quota_percent_used                   | Memory used by the bucket in percent                                                                    |
| cb_bucket_ops_per_second                           | Number of operations per second                                                                         |
| cb_bucket_disk_fetches                             | Disk fetches for the bucket                                                                             |
| cb_bucket_item_count                               | Number of items in the bucket                                                                           |
| cb_bucket_disk_used_bytes                          | Disk used by the bucket                                                                                 |
| cb_bucket_data_used_bytes                          | Data loaded in memory                                                                                   |
| cb_bucket_ram_used_bytes                           | Bucket RAM used                                                                                         |
| cb_bucket_couch_total_disk_size                    | Couchbase total disk size                                                                               |
| cb_bucket_couch_docs_fragmentation                 | Couchbase documents fragmentation                                                                       |
| cb_bucket_couch_views_fragmentation                | Couchbase views fragmentation                                                                           |
| cb_bucket_hit_ratio                                | Hit ratio                                                                                               |
| cb_bucket_ep_cache_miss_rate                       | Cache miss rate                                                                                         |
| cb_bucket_ep_resident_items_rate                   | Number of resident items                                                                                |
| cb_bucket_vb_avg_active_queue_age                  | Average age in seconds of active items in the active item queue                                         |
| cb_bucket_vb_avg_replica_queue_age                 | Average age in seconds of replica items in the replica item queue                                       |
| cb_bucket_vb_avg_pending_queue_age                 | Average age in seconds of pending items in the pending item queue                                       |
| cb_bucket_vb_avg_total_queue_age                   | Average age of items in the queue                                                                       |
| cb_bucket_vb_active_resident_items_ratio           | Number of resident items                                                                                |
| cb_bucket_vb_replica_resident_items_ratio          | Number of resident replica items                                                                        |
| cb_bucket_vb_pending_resident_items_ratio          | Number of resident pending items                                                                        |
| cb_bucket_avg_disk_update_time                     | Average disk update time                                                                                |
| cb_bucket_avg_disk_commit_time                     | Average disk commit time                                                                                |
| cb_bucket_avg_bg_wait_time                         | Average background wait time                                                                            |
| cb_bucket_ep_dcp_views_indexes_count               | Number of indexes views DCP connections                                                                 |
| cb_bucket_ep_dcp_views_indexes_items_remaining     | Number of indexes views items remaining to be sent                                                      |
| cb_bucket_ep_dcp_views_indexes_producer_count      | Number of indexes views producers                                                                       |
| cb_bucket_ep_dcp_views_indexes_total_backlog_size  | Number of indexes views items remaining for replication                                                 |
| cb_bucket_ep_dcp_views_indexes_items_sent          | Number of indexes views sent                                                                            |
| cb_bucket_ep_dcp_views_indexes_total_bytes         | Number of bytes per second being sent for indexes views DCP                                             |
| cb_bucket_ep_dcp_views_indexes_backoff             | Number of backoffs for indexes views DCP connections                                                    |
| cb_bucket_bg_wait_count                            | Background wait                                                                                         |
| cb_bucket_bg_wait_total                            | Total background wait                                                                                   |
| cb_bucket_bytes_read                               | Bytes read                                                                                              |
| cb_bucket_bytes_written                            | Bytes written                                                                                           |
| cb_bucket_cas_badval                               | Compare and Swap bad values                                                                             |
| cb_bucket_cas_hits                                 | Compare and Swap hits                                                                                   |
| cb_bucket_cas_misses                               | Compare and Swap misses                                                                                 |
| cb_bucket_cmd_get                                  | Gets from memory                                                                                        |
| cb_bucket_cmd_set                                  | Sets to memory                                                                                          |
| cb_bucket_couch_docs_actual_disk_size              | Total size of documents on disk in bytes                                                                |
| cb_bucket_couch_docs_data_size                     | Documents size in bytes                                                                                 |
| cb_bucket_couch_docs_disk_size                     | Total size of documents in bytes                                                                        |
| cb_bucket_couch_spatial_data_size                  | Size of object data for spatial views                                                                   |
| cb_bucket_couch_spatial_disk_size                  | Amount of disk space occupied by spatial views                                                          |
| cb_bucket_couch_spatial_ops                        | Spatial operations                                                                                      |
| cb_bucket_couch_views_actual_disk_size             | Total size of views on disk in bytes                                                                    |
| cb_bucket_couch_views_data_size                    | Views size in bytes                                                                                     |
| cb_bucket_couch_views_disk_size                    | Total size of views in bytes                                                                            |
| cb_bucket_couch_views_ops                          | View operations                                                                                         |
| cb_bucket_curr_connections                         | Current bucket connections                                                                              |
| cb_bucket_curr_items                               | Number of active items in memory                                                                        |
| cb_bucket_curr_items_tot                           | Total number of items                                                                                   |
| cb_bucket_decr_hits                                | Decrement hits                                                                                          |
| cb_bucket_decr_misses                              | Decrement misses                                                                                        |
| cb_bucket_delete_hits                              | Delete hits                                                                                             |
| cb_bucket_delete_misses                            | Delete misses                                                                                           |
| cb_bucket_disk_commit_count                        | Disk commits                                                                                            |
| cb_bucket_disk_commit_total                        | Total disk commits                                                                                      |
| cb_bucket_disk_update_count                        | Disk updates                                                                                            |
| cb_bucket_disk_update_total                        | Total disk updates                                                                                      |
| cb_bucket_disk_write_queue                         | Disk write queue depth                                                                                  |
| cb_bucket_ep_bg_fetched                            | Disk reads per second                                                                                   |
| cb_bucket_ep_dcp_2i_backoff                        | Number of backoffs for indexes DCP connections                                                          |
| cb_bucket_ep_dcp_2i_count                          | Number of indexes DCP connections                                                                       |
| cb_bucket_ep_dcp_2i_items_remaining                | Number of indexes items remaining to be sent                                                            |
| cb_bucket_ep_dcp_2i_items_sent                     | Number of indexes items sent                                                                            |
| cb_bucket_ep_dcp_2i_producer_count                 | Number of indexes producers                                                                             |
| cb_bucket_ep_dcp_2i_total_backlog_size             | Number of indexes total backlog size                                                                    |
| cb_bucket_ep_dcp_2i_total_bytes                    | Number bytes per second being sent for indexes DCP connections                                          |
| cb_bucket_ep_dcp_fts_backoff                       | Number of backoffs for fts DCP connections                                                              |
| cb_bucket_ep_dcp_fts_count                         | Number of fts DCP connections                                                                           |
| cb_bucket_ep_dcp_fts_items_remaining               | Number of fts items remaining to be sent                                                                |
| cb_bucket_ep_dcp_fts_items_sent                    | Number of fts items sent                                                                                |
| cb_bucket_ep_dcp_fts_producer_count                | Number of fts producers                                                                                 |
| cb_bucket_ep_dcp_fts_total_backlog_size            | Number of fts total backlog size                                                                        |
| cb_bucket_ep_dcp_fts_total_bytes                   | Number bytes per second being sent for fts DCP connections                                              |
| cb_bucket_ep_dcp_other_backoff                     | Number of backoffs for other DCP connections                                                            |
| cb_bucket_ep_dcp_other_count                       | Number of other DCP connections                                                                         |
| cb_bucket_ep_dcp_other_items_remaining             | Number of other items remaining to be sent                                                              |
| cb_bucket_ep_dcp_other_items_sent                  | Number of other items sent                                                                              |
| cb_bucket_ep_dcp_other_producer_count              | Number of other producers                                                                               |
| cb_bucket_ep_dcp_other_total_backlog_size          | Number of other total backlog size                                                                      |
| cb_bucket_ep_dcp_other_total_bytes                 | Number bytes per second being sent for other DCP connections                                            |
| cb_bucket_ep_dcp_replica_backoff                   | Number of backoffs for replica DCP connections                                                          |
| cb_bucket_ep_dcp_replica_count                     | Number of replica DCP connections                                                                       |
| cb_bucket_ep_dcp_replica_items_remaining           | Number of replica items remaining to be sent                                                            |
| cb_bucket_ep_dcp_replica_items_sent                | Number of replica items sent                                                                            |
| cb_bucket_ep_dcp_replica_producer_count            | Number of replica producers                                                                             |
| cb_bucket_ep_dcp_replica_total_backlog_size        | Number of replica total backlog size                                                                    |
| cb_bucket_ep_dcp_replica_total_bytes               | Number bytes per second being sent for replica DCP connections                                          |
| cb_bucket_ep_dcp_views_backoff                     | Number of backoffs for views DCP connections                                                            |
| cb_bucket_ep_dcp_views_count                       | Number of views DCP connections                                                                         |
| cb_bucket_ep_dcp_views_items_remaining             | Number of views items remaining to be sent                                                              |
| cb_bucket_ep_dcp_views_items_sent                  | Number of views items sent                                                                              |
| cb_bucket_ep_dcp_views_producer_count              | Number of views producers                                                                               |
| cb_bucket_ep_dcp_views_total_backlog_size          | Number of views total backlog size                                                                      |
| cb_bucket_ep_dcp_views_total_bytes                 | Number bytes per second being sent for views DCP connections                                            |
| cb_bucket_ep_dcp_xdcr_backoff                      | Number of backoffs for xdcr DCP connections                                                             |
| cb_bucket_ep_dcp_xdcr_count                        | Number of xdcr DCP connections                                                                          |
| cb_bucket_ep_dcp_xdcr_items_remaining              | Number of xdcr items remaining to be sent                                                               |
| cb_bucket_ep_dcp_xdcr_items_sent                   | Number of xdcr items sent                                                                               |
| cb_bucket_ep_dcp_xdcr_producer_count               | Number of xdcr producers                                                                                |
| cb_bucket_ep_dcp_xdcr_total_backlog_size           | Number of xdcr total backlog size                                                                       |
| cb_bucket_ep_dcp_xdcr_total_bytes                  | Number bytes per second being sent for xdcr DCP connections                                             |
| cb_bucket_ep_diskqueue_drain                       | Total Drained items on disk queue                                                                       |
| cb_bucket_ep_diskqueue_fill                        | Total enqueued items on disk queue                                                                      |
| cb_bucket_ep_diskqueue_items                       | Total number of items waiting to be written to disk                                                     |
| cb_bucket_ep_flusher_todo                          | Number of items currently being written                                                                 |
| cb_bucket_ep_item_commit_failed                    | Number of times a transaction failed to commit due to storage errors                                    |
| cb_bucket_ep_kv_size                               | Total amount of user data cached in RAM                                                                 |
| cb_bucket_ep_max_size                              | Maximum amount of memory this bucket can use                                                            |
| cb_bucket_ep_mem_high_wat                          | Memory usage high water mark for auto-evictions                                                         |
| cb_bucket_ep_mem_low_wat                           | Memory usage low water mark for auto-evictions                                                          |
| cb_bucket_ep_meta_data_memory                      | Total amount of item metadata consuming RAM                                                             |
| cb_bucket_ep_num_non_resident                      | Number of non-resident items                                                                            |
| cb_bucket_ep_num_ops_del_meta                      | Number of delete operations per second for this bucket as the target for XDCR                           |
| cb_bucket_ep_num_ops_del_ret_meta                  | Number of delRetMeta operations per second for this bucket as the target for XDCR                       |
| cb_bucket_ep_num_ops_get_meta                      | Number of read operations per second for this bucket as the target for XDCR                             |
| cb_bucket_ep_num_ops_set_meta                      | Number of write operations per second for this bucket as the target for XDCR                            |
| cb_bucket_ep_num_ops_set_ret_meta                  | Number of setRetMeta operations per second for this bucket as the target for XDCR                       |
| cb_bucket_ep_num_value_ejects                      | Number of times item values got ejected from memory to disk                                             |
| cb_bucket_ep_oom_errors                            | Number of times unrecoverable OOMs happened while processing operations                                 |
| cb_bucket_ep_ops_create                            | Create operations                                                                                       |
| cb_bucket_ep_ops_update                            | Update operations                                                                                       |
| cb_bucket_ep_overhead                              | Extra memory used by transient data like persistence queues or checkpoints                              |
| cb_bucket_ep_queue_size                            | Number of items queued for storage                                                                      |
| cb_bucket_ep_tmp_oom_errors                        | Number of times recoverable OOMs happened while processing operations                                   |
| cb_bucket_ep_vb_total                              | Total number of vBuckets for this bucket                                                                |
| cb_bucket_evictions                                | Number of evictions                                                                                     |
| cb_bucket_get_hits                                 | Number of get hits                                                                                      |
| cb_bucket_get_misses                               | Number of get misses                                                                                    |
| cb_bucket_incr_hits                                | Number of increment hits                                                                                |
| cb_bucket_incr_misses                              | Number of increment misses                                                                              |
| cb_bucket_mem_used                                 | Engine's total memory usage (deprecated)                                                                |
| cb_bucket_misses                                   | Total number of misses                                                                                  |
| cb_bucket_ops                                      | Total number of operations                                                                              |
| cb_bucket_vb_active_eject                          | Number of items per second being ejected to disk from active vBuckets                                   |
| cb_bucket_vb_active_itm_memory                     | Amount of active user data cached in RAM                                                                |
| cb_bucket_vb_active_meta_data_memory               | Amount of active item metadata consuming RAM                                                            |
| cb_bucket_vb_active_num                            | Number of active items                                                                                  |
| cb_bucket_vb_active_num_non_resident               | Number of non resident vBuckets in the active state for this bucket                                     |
| cb_bucket_vb_active_ops_create                     | New items per second being inserted into active vBuckets                                                |
| cb_bucket_vb_active_ops_update                     | Number of items updated on active vBucket per second for this bucket                                    |
| cb_bucket_vb_active_queue_age                      | Sum of disk queue item age in milliseconds                                                              |
| cb_bucket_vb_active_queue_drain                    | Total drained items in the queue                                                                        |
| cb_bucket_vb_active_queue_fill                     | Number of active items per second being put on the active item disk queue                               |
| cb_bucket_vb_active_queue_size                     | Number of active items in the queue                                                                     |
| cb_bucket_vb_pending_curr_items                    | Number of items in pending vBuckets                                                                     |
| cb_bucket_vb_pending_eject                         | Number of items per second being ejected to disk from pending vBuckets                                  |
| cb_bucket_vb_pending_itm_memory                    | Amount of pending user data cached in RAM                                                               |
| cb_bucket_vb_pending_meta_data_memory              | Amount of pending item metadata consuming RAM                                                           |
| cb_bucket_vb_pending_num                           | Number of pending items                                                                                 |
| cb_bucket_vb_pending_num_non_resident              | Number of non resident vBuckets in the pending state for this bucket                                    |
| cb_bucket_vb_pending_ops_create                    | Number of pending create operations                                                                     |
| cb_bucket_vb_pending_ops_update                    | Number of items updated on pending vBucket per second for this bucket                                   |
| cb_bucket_vb_pending_queue_age                     | Sum of disk pending queue item age in milliseconds                                                      |
| cb_bucket_vb_pending_queue_drain                   | Total drained pending items in the queue                                                                |
| cb_bucket_vb_pending_queue_fill                    | Total enqueued pending items on disk queue                                                              |
| cb_bucket_vb_pending_queue_size                    | Number of pending items in the queue                                                                    |
| cb_bucket_vb_replica_curr_items                    | Number of in memory items                                                                               |
| cb_bucket_vb_replica_eject                         | Number of items per second being ejected to disk from replica vBuckets                                  |
| cb_bucket_vb_replica_itm_memory                    | Amount of replica user data cached in RAM                                                               |
| cb_bucket_vb_replica_meta_data_memory              | Total metadata memory                                                                                   |
| cb_bucket_vb_replica_num                           | Number of replica vBuckets                                                                              |
| cb_bucket_vb_replica_num_non_resident              | Number of non resident vBuckets in the replica state for this bucket                                    |
| cb_bucket_vb_replica_ops_create                    | Number of replica create operations                                                                     |
| cb_bucket_vb_replica_ops_update                    | Number of items updated on replica vBucket per second for this bucket                                   |
| cb_bucket_vb_replica_queue_age                     | Sum of disk replica queue item age in milliseconds                                                      |
| cb_bucket_vb_replica_queue_drain                   | Total drained replica items in the queue                                                                |
| cb_bucket_vb_replica_queue_fill                    | Total enqueued replica items on disk queue                                                              |
| cb_bucket_vb_replica_queue_size                    | Replica items in disk queue                                                                             |
| cb_bucket_vb_total_queue_age                       | Sum of disk queue item age in milliseconds                                                              |
| cb_bucket_xdc_ops                                  | Number of cross-datacenter replication operations                                                       |
| cb_bucket_cpu_idle_ms                              | CPU idle milliseconds                                                                                   |
| cb_bucket_cpu_local_ms                             | CPU local milliseconds                                                                                  |
| cb_bucket_cpu_utilization_rate                     | CPU utilization percentage                                                                              |
| cb_bucket_hibernated_requests                      | Number of streaming requests now idle                                                                   |
| cb_bucket_hibernated_waked                         | Rate of streaming request wakeups                                                                       |
| cb_bucket_mem_actual_free                          | Actual free memory                                                                                      |
| cb_bucket_mem_actual_used                          | Actual used memory                                                                                      |
| cb_bucket_mem_free                                 | Free memory                                                                                             |
| cb_bucket_mem_total                                | Total memeory                                                                                           |
| cb_bucket_mem_used_sys                             | System memory usage                                                                                     |
| cb_bucket_rest_requests                            | Number of HTTP requests                                                                                 |
| cb_bucket_swap_total                               | Total amount of swap available                                                                          |
| cb_bucket_swap_used                                | Amount of swap used                                                                                     |
| cb_bucket_ep_tap_rebalance_count                   | Number of internal rebalancing TAP queues                                                               |
| cb_bucket_ep_tap_rebalance_qlen                    | Number of items in the rebalance TAP queues                                                             |
| cb_bucket_ep_tap_rebalance_queue_backfillremaining | Number of items in the backfill queues of rebalancing TAP connections                                   |
| cb_bucket_ep_tap_rebalance_queue_backoff           | Number of back-offs received per second while sending data over rebalancing connections                 |
| cb_bucket_ep_tap_rebalance_queue_drain             | Number of items per second being sent over rebalancing TAP connections, i.e. removed from queue         |
| cb_bucket_ep_tap_rebalance_queue_fill              | Number of items per second being sent to queue                                                          |
| cb_bucket_ep_tap_rebalance_queue_itemondisk        | Number of items still on disk to be loaded for rebalancing TAP connections                              |
| cb_bucket_ep_tap_rebalance_total_backlog_size      | Number of remaining items for rebalancing TAP connections                                               |
| cb_bucket_ep_tap_replica_count                     | Number of internal replication TAP queues                                                               |
| cb_bucket_ep_tap_replica_qlen                      | Number of items in the replication TAP queues                                                           |
| cb_bucket_ep_tap_replica_queue_backfillremaining   | Number of items in the backfill queues of replication TAP connections                                   |
| cb_bucket_ep_tap_replica_queue_backoff             | Number of back-offs received per second while sending data over replication connections                 |
| cb_bucket_ep_tap_replica_queue_drain               | Total drained items in the replica queue                                                                |
| cb_bucket_ep_tap_replica_queue_fill                | Number of items per second being sent to queue                                                          |
| cb_bucket_ep_tap_replica_queue_itemondisk          | Number of items still on disk to be loaded for replication TAP connections                              |
| cb_bucket_ep_tap_replica_total_backlog_size        | Number of remaining items for replication TAP connections                                               |
| cb_bucket_ep_tap_total_count                       | Total number of internal TAP queues                                                                     |
| cb_bucket_ep_tap_total_qlen                        | Total number of items in TAP queues                                                                     |
| cb_bucket_ep_tap_total_queue_backfillremaining     | Total number of items in the backfill queues of TAP connections                                         |
| cb_bucket_ep_tap_total_queue_backoff               | Total number of back-offs received per second while sending data over TAP connections                   |
| cb_bucket_ep_tap_total_queue_drain                 | Total drained items in the queue                                                                        |
| cb_bucket_ep_tap_total_queue_fill                  | Total enqueued items in the queue                                                                       |
| cb_bucket_ep_tap_total_queue_itemondisk            | Total number of items still on disk to be loaded for TAP connections                                    |
| cb_bucket_ep_tap_total_total_backlog_size          | Number of remaining items for replication                                                               |
| cb_bucket_ep_tap_user_count                        | Number of internal user TAP queues                                                                      |
| cb_bucket_ep_tap_user_qlen                         | Number of items in user TAP queues                                                                      |
| cb_bucket_ep_tap_user_queue_backfillremaining      | Number of items in the backfill queues of user TAP connections                                          |
| cb_bucket_ep_tap_user_queue_backoff                | Number of back-offs received per second while sending data over user TAP connections                    |
| cb_bucket_ep_tap_user_queue_drain                  | Number of items per second being sent over user TAP connections to this bucket, i.e. removed from queue |
| cb_bucket_ep_tap_user_queue_fill                   | Number of items per second being sent to queue                                                          |
| cb_bucket_ep_tap_user_queue_itemondisk             | Number of items still on disk to be loaded for client TAP connections                                   |
| cb_bucket_ep_tap_user_total_backlog_size           | Number of remaining items for client TAP connections                                                    |
| cb_bucket_avg_active_timestamp_drift               | Average active timestamp drift                                                                          |
| cb_bucket_avg_replica_timestamp_drift              | Average replica timestamp drift                                                                         |
| cb_bucket_ep_active_ahead_exceptions               | Sum total of all active vBuckets drift_ahead_threshold_exceeded counter                                 |
| cb_bucket_ep_active_hlc_drift                      | Total absolute drift for all active vBuckets                                                            |
| cb_bucket_ep_active_hlc_drift_count                | Number of updates applied to ep_active_hlc_drift                                                        |
| cb_bucket_ep_clock_cas_drift_threshold_exceeded    | Ep clock cas drift threshold exceeded                                                                   |
| cb_bucket_ep_replica_ahead_exceptions              | Sum total of all replica vBuckets' drift_ahead_threshold_exceeded counter                               |
| cb_bucket_ep_replica_hlc_drift                     | Total absolute drift for all replica vBuckets                                                           |
| cb_bucket_ep_replica_hlc_drift_count               | Number of updates applied to ep_replica_hlc_drift                                                       |

### XDCR metrics

|              name              |                                                     description                                                     |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------- |
| cb_xdcr_bandwidth_usage        | Bandwidth used during replication, measured in bytes per second                                                     |
| cb_xdcr_changes_left           | Number of updates still pending replication                                                                         |
| cb_xdcr_data_replicated        | Size of data replicated in bytes                                                                                    |
| cb_xdcr_docs_checked           | Number of documents checked for changes                                                                             |
| cb_xdcr_docs_failed_cr_source  | Number of documents that have failed conflict resolution on the source cluster and not replicated to target cluster |
| cb_xdcr_docs_filtered          | Number of documents that have been filtered out and not replicated to target cluster                                |
| cb_xdcr_docs_opt_repd          | Number of docs sent optimistically                                                                                  |
| cb_xdcr_docs_received_from_dcp | Number of documents received from DCP                                                                               |
| cb_xdcr_docs_rep_queue         | Number of documents in replication queue                                                                            |
| cb_xdcr_docs_written           | Number of documents written to the destination cluster via XDCR                                                     |
| cb_xdcr_meta_latency_wt        | Weighted average time for requesting document metadata                                                              |
| cb_xdcr_num_checkpoints        | Number of checkpoints issued in replication queue                                                                   |
| cb_xdcr_num_failedckpts        | Number of checkpoints failed during replication                                                                     |
| cb_xdcr_rate_received_from_dcp | Number of documents received from DCP per second                                                                    |
| cb_xdcr_rate_replicated        | Rate of documents being replicated, measured in documents per second                                                |
| cb_xdcr_size_rep_queue         | Size of replication queue in bytes                                                                                  |
| cb_xdcr_time_committing        | Seconds elapsed during replication                                                                                  |
| cb_xdcr_wtavg_docs_latency     | Weighted average latency for sending replicated changes to destination cluster                                      |
| cb_xdcr_wtavg_meta_latency     | Weighted average time for requesting document metadata                                                              |
| cb_xdcr_percent_completeness   | Percentage of checked items out of all checked and to-be-replicated items                                           |
| cb_xdcr_error_count            | Number of XDCR errors                                                                                               |

## Docker

Get the latest image from Docker Hub:

```bash
docker pull blakelead/couchbase-exporter:latest
```

and use it with environment variables like in the following example:

```bash
docker run --name cbexporter -p 9191:9191 -e CB_EXPORTER_DB_USER=admin -e CB_EXPORTER_DB_PASSWORD=complicatedpassword blakelead/couchbase-exporter:latest
```

## Examples

You can find example files in `examples` directory.

### Prometheus

Some alerting rules that can be used in Prometheus configuration for alert manager notifications.

### Grafana

Minimal dashboard that shows how to use the exporter metrics.

### Systemd

You can adapt and use the service template to run **couchbase_exporter** with systemd.

```bash
sudo mv systemd-exporter.service /etc/systemd/system/couchbase-exporter.service
sudo systemctl enable couchbase-exporter.service
sudo systemctl start couchbase-exporter.service
```

## Author Information

Adel Abdelhak
Additions Don Thomson
