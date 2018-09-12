// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// BucketStatsData (/pools/default/buckets/<bucket_name>/stats)
type BucketStatsData struct {
	Op struct {
		Samples struct {
			CouchTotalDiskSize                   []float64 `json:"couch_total_disk_size"`
			CouchDocsFragmentation               []float64 `json:"couch_docs_fragmentation"`
			CouchViewsFragmentation              []float64 `json:"couch_views_fragmentation"`
			HitRatio                             []float64 `json:"hit_ratio"`
			EpCacheMissRate                      []float64 `json:"ep_cache_miss_rate"`
			EpResidentItemsRate                  []float64 `json:"ep_resident_items_rate"`
			VbAvgActiveQueueAge                  []float64 `json:"vb_avg_active_queue_age"`
			VbAvgReplicaQueueAge                 []float64 `json:"vb_avg_replica_queue_age"`
			VbAvgPendingQueueAge                 []float64 `json:"vb_avg_pending_queue_age"`
			VbAvgTotalQueueAge                   []float64 `json:"vb_avg_total_queue_age"`
			VbActiveResidentItemsRatio           []float64 `json:"vb_active_resident_items_ratio"`
			VbReplicaResidentItemsRatio          []float64 `json:"vb_replica_resident_items_ratio"`
			VbPendingResidentItemsRatio          []float64 `json:"vb_pending_resident_items_ratio"`
			AvgDiskUpdateTime                    []float64 `json:"avg_disk_update_time"`
			AvgDiskCommitTime                    []float64 `json:"avg_disk_commit_time"`
			AvgBgWaitTime                        []float64 `json:"avg_bg_wait_time"`
			AvgActiveTimestampDrift              []float64 `json:"avg_active_timestamp_drift"`  // couchbase 5.1.1
			AvgReplicaTimestampDrift             []float64 `json:"avg_replica_timestamp_drift"` // couchbase 5.1.1
			EpDcpViewsIndexesCount               []float64 `json:"ep_dcp_views+indexes_count"`
			EpDcpViewsIndexesItemsRemaining      []float64 `json:"ep_dcp_views+indexes_items_remaining"`
			EpDcpViewsIndexesProducerCount       []float64 `json:"ep_dcp_views+indexes_producer_count"`
			EpDcpViewsIndexesTotalBacklogSize    []float64 `json:"ep_dcp_views+indexes_total_backlog_size"`
			EpDcpViewsIndexesItemsSent           []float64 `json:"ep_dcp_views+indexes_items_sent"`
			EpDcpViewsIndexesTotalBytes          []float64 `json:"ep_dcp_views+indexes_total_bytes"`
			EpDcpViewsIndexesBackoff             []float64 `json:"ep_dcp_views+indexes_backoff"`
			BgWaitCount                          []float64 `json:"bg_wait_count"`
			BgWaitTotal                          []float64 `json:"bg_wait_total"`
			BytesRead                            []float64 `json:"bytes_read"`
			BytesWritten                         []float64 `json:"bytes_written"`
			CasBadval                            []float64 `json:"cas_badval"`
			CasHits                              []float64 `json:"cas_hits"`
			CasMisses                            []float64 `json:"cas_misses"`
			CmdGet                               []float64 `json:"cmd_get"`
			CmdSet                               []float64 `json:"cmd_set"`
			CouchDocsActualDiskSize              []float64 `json:"couch_docs_actual_disk_size"`
			CouchDocsDataSize                    []float64 `json:"couch_docs_data_size"`
			CouchDocsDiskSize                    []float64 `json:"couch_docs_disk_size"`
			CouchSpatialDataSize                 []float64 `json:"couch_spatial_data_size"`
			CouchSpatialDiskSize                 []float64 `json:"couch_spatial_disk_size"`
			CouchSpatialOps                      []float64 `json:"couch_spatial_ops"`
			CouchViewsActualDiskSize             []float64 `json:"couch_views_actual_disk_size"`
			CouchViewsDataSize                   []float64 `json:"couch_views_data_size"`
			CouchViewsDiskSize                   []float64 `json:"couch_views_disk_size"`
			CouchViewsOps                        []float64 `json:"couch_views_ops"`
			CurrConnections                      []float64 `json:"curr_connections"`
			CurrItems                            []float64 `json:"curr_items"`
			CurrItemsTot                         []float64 `json:"curr_items_tot"`
			DecrHits                             []float64 `json:"decr_hits"`
			DecrMisses                           []float64 `json:"decr_misses"`
			DeleteHits                           []float64 `json:"delete_hits"`
			DeleteMisses                         []float64 `json:"delete_misses"`
			DiskCommitCount                      []float64 `json:"disk_commit_count"`
			DiskCommitTotal                      []float64 `json:"disk_commit_total"`
			DiskUpdateCount                      []float64 `json:"disk_update_count"`
			DiskUpdateTotal                      []float64 `json:"disk_update_total"`
			DiskWriteQueue                       []float64 `json:"disk_write_queue"`
			EpActiveAheadExceptions              []float64 `json:"ep_active_ahead_exceptions"` // couchbase 5.1.1
			EpActiveHlcDrift                     []float64 `json:"ep_active_hlc_drift"`        // couchbase 5.1.1
			EpActiveHlcDriftCount                []float64 `json:"ep_active_hlc_drift_count"`  // couchbase 5.1.1
			EpBgFetched                          []float64 `json:"ep_bg_fetched"`
			EpClockCasDriftThresholdExceeded     []float64 `json:"ep_clock_cas_drift_threshold_exceeded"` // couchbase 5.1.1
			EpDcp2IBackoff                       []float64 `json:"ep_dcp_2i_backoff"`
			EpDcp2ICount                         []float64 `json:"ep_dcp_2i_count"`
			EpDcp2IItemsRemaining                []float64 `json:"ep_dcp_2i_items_remaining"`
			EpDcp2IItemsSent                     []float64 `json:"ep_dcp_2i_items_sent"`
			EpDcp2IProducerCount                 []float64 `json:"ep_dcp_2i_producer_count"`
			EpDcp2ITotalBacklogSize              []float64 `json:"ep_dcp_2i_total_backlog_size"`
			EpDcp2ITotalBytes                    []float64 `json:"ep_dcp_2i_total_bytes"`
			EpDcpFtsBackoff                      []float64 `json:"ep_dcp_fts_backoff"`
			EpDcpFtsCount                        []float64 `json:"ep_dcp_fts_count"`
			EpDcpFtsItemsRemaining               []float64 `json:"ep_dcp_fts_items_remaining"`
			EpDcpFtsItemsSent                    []float64 `json:"ep_dcp_fts_items_sent"`
			EpDcpFtsProducerCount                []float64 `json:"ep_dcp_fts_producer_count"`
			EpDcpFtsTotalBacklogSize             []float64 `json:"ep_dcp_fts_total_backlog_size"`
			EpDcpFtsTotalBytes                   []float64 `json:"ep_dcp_fts_total_bytes"`
			EpDcpOtherBackoff                    []float64 `json:"ep_dcp_other_backoff"`
			EpDcpOtherCount                      []float64 `json:"ep_dcp_other_count"`
			EpDcpOtherItemsRemaining             []float64 `json:"ep_dcp_other_items_remaining"`
			EpDcpOtherItemsSent                  []float64 `json:"ep_dcp_other_items_sent"`
			EpDcpOtherProducerCount              []float64 `json:"ep_dcp_other_producer_count"`
			EpDcpOtherTotalBacklogSize           []float64 `json:"ep_dcp_other_total_backlog_size"`
			EpDcpOtherTotalBytes                 []float64 `json:"ep_dcp_other_total_bytes"`
			EpDcpReplicaBackoff                  []float64 `json:"ep_dcp_replica_backoff"`
			EpDcpReplicaCount                    []float64 `json:"ep_dcp_replica_count"`
			EpDcpReplicaItemsRemaining           []float64 `json:"ep_dcp_replica_items_remaining"`
			EpDcpReplicaItemsSent                []float64 `json:"ep_dcp_replica_items_sent"`
			EpDcpReplicaProducerCount            []float64 `json:"ep_dcp_replica_producer_count"`
			EpDcpReplicaTotalBacklogSize         []float64 `json:"ep_dcp_replica_total_backlog_size"`
			EpDcpReplicaTotalBytes               []float64 `json:"ep_dcp_replica_total_bytes"`
			EpDcpViewsBackoff                    []float64 `json:"ep_dcp_views_backoff"`
			EpDcpViewsCount                      []float64 `json:"ep_dcp_views_count"`
			EpDcpViewsItemsRemaining             []float64 `json:"ep_dcp_views_items_remaining"`
			EpDcpViewsItemsSent                  []float64 `json:"ep_dcp_views_items_sent"`
			EpDcpViewsProducerCount              []float64 `json:"ep_dcp_views_producer_count"`
			EpDcpViewsTotalBacklogSize           []float64 `json:"ep_dcp_views_total_backlog_size"`
			EpDcpViewsTotalBytes                 []float64 `json:"ep_dcp_views_total_bytes"`
			EpDcpXdcrBackoff                     []float64 `json:"ep_dcp_xdcr_backoff"`
			EpDcpXdcrCount                       []float64 `json:"ep_dcp_xdcr_count"`
			EpDcpXdcrItemsRemaining              []float64 `json:"ep_dcp_xdcr_items_remaining"`
			EpDcpXdcrItemsSent                   []float64 `json:"ep_dcp_xdcr_items_sent"`
			EpDcpXdcrProducerCount               []float64 `json:"ep_dcp_xdcr_producer_count"`
			EpDcpXdcrTotalBacklogSize            []float64 `json:"ep_dcp_xdcr_total_backlog_size"`
			EpDcpXdcrTotalBytes                  []float64 `json:"ep_dcp_xdcr_total_bytes"`
			EpDiskqueueDrain                     []float64 `json:"ep_diskqueue_drain"`
			EpDiskqueueFill                      []float64 `json:"ep_diskqueue_fill"`
			EpDiskqueueItems                     []float64 `json:"ep_diskqueue_items"`
			EpFlusherTodo                        []float64 `json:"ep_flusher_todo"`
			EpItemCommitFailed                   []float64 `json:"ep_item_commit_failed"`
			EpKvSize                             []float64 `json:"ep_kv_size"`
			EpMaxSize                            []float64 `json:"ep_max_size"`
			EpMemHighWat                         []float64 `json:"ep_mem_high_wat"`
			EpMemLowWat                          []float64 `json:"ep_mem_low_wat"`
			EpMetaDataMemory                     []float64 `json:"ep_meta_data_memory"`
			EpNumNonResident                     []float64 `json:"ep_num_non_resident"`
			EpNumOpsDelMeta                      []float64 `json:"ep_num_ops_del_meta"`
			EpNumOpsDelRetMeta                   []float64 `json:"ep_num_ops_del_ret_meta"`
			EpNumOpsGetMeta                      []float64 `json:"ep_num_ops_get_meta"`
			EpNumOpsSetMeta                      []float64 `json:"ep_num_ops_set_meta"`
			EpNumOpsSetRetMeta                   []float64 `json:"ep_num_ops_set_ret_meta"`
			EpNumValueEjects                     []float64 `json:"ep_num_value_ejects"`
			EpOomErrors                          []float64 `json:"ep_oom_errors"`
			EpOpsCreate                          []float64 `json:"ep_ops_create"`
			EpOpsUpdate                          []float64 `json:"ep_ops_update"`
			EpOverhead                           []float64 `json:"ep_overhead"`
			EpQueueSize                          []float64 `json:"ep_queue_size"`
			EpReplicaAheadExceptions             []float64 `json:"ep_replica_ahead_exceptions"`              // couchbase 5.1.1
			EpReplicaHlcDrift                    []float64 `json:"ep_replica_hlc_drift"`                     // couchbase 5.1.1
			EpReplicaHlcDriftCount               []float64 `json:"ep_replica_hlc_drift_count"`               // couchbase 5.1.1
			EpTapRebalanceCount                  []float64 `json:"ep_tap_rebalance_count"`                   // couchbase 4.5.1
			EpTapRebalanceQlen                   []float64 `json:"ep_tap_rebalance_qlen"`                    // couchbase 4.5.1
			EpTapRebalanceQueueBackfillremaining []float64 `json:"ep_tap_rebalance_queue_backfillremaining"` // couchbase 4.5.1
			EpTapRebalanceQueueBackoff           []float64 `json:"ep_tap_rebalance_queue_backoff"`           // couchbase 4.5.1
			EpTapRebalanceQueueDrain             []float64 `json:"ep_tap_rebalance_queue_drain"`             // couchbase 4.5.1
			EpTapRebalanceQueueFill              []float64 `json:"ep_tap_rebalance_queue_fill"`              // couchbase 4.5.1
			EpTapRebalanceQueueItemondisk        []float64 `json:"ep_tap_rebalance_queue_itemondisk"`        // couchbase 4.5.1
			EpTapRebalanceTotalBacklogSize       []float64 `json:"ep_tap_rebalance_total_backlog_size"`      // couchbase 4.5.1
			EpTapReplicaCount                    []float64 `json:"ep_tap_replica_count"`                     // couchbase 4.5.1
			EpTapReplicaQlen                     []float64 `json:"ep_tap_replica_qlen"`                      // couchbase 4.5.1
			EpTapReplicaQueueBackfillremaining   []float64 `json:"ep_tap_replica_queue_backfillremaining"`   // couchbase 4.5.1
			EpTapReplicaQueueBackoff             []float64 `json:"ep_tap_replica_queue_backoff"`             // couchbase 4.5.1
			EpTapReplicaQueueDrain               []float64 `json:"ep_tap_replica_queue_drain"`               // couchbase 4.5.1
			EpTapReplicaQueueFill                []float64 `json:"ep_tap_replica_queue_fill"`                // couchbase 4.5.1
			EpTapReplicaQueueItemondisk          []float64 `json:"ep_tap_replica_queue_itemondisk"`          // couchbase 4.5.1
			EpTapReplicaTotalBacklogSize         []float64 `json:"ep_tap_replica_total_backlog_size"`        // couchbase 4.5.1
			EpTapTotalCount                      []float64 `json:"ep_tap_total_count"`                       // couchbase 4.5.1
			EpTapTotalQlen                       []float64 `json:"ep_tap_total_qlen"`                        // couchbase 4.5.1
			EpTapTotalQueueBackfillremaining     []float64 `json:"ep_tap_total_queue_backfillremaining"`     // couchbase 4.5.1
			EpTapTotalQueueBackoff               []float64 `json:"ep_tap_total_queue_backoff"`               // couchbase 4.5.1
			EpTapTotalQueueDrain                 []float64 `json:"ep_tap_total_queue_drain"`                 // couchbase 4.5.1
			EpTapTotalQueueFill                  []float64 `json:"ep_tap_total_queue_fill"`                  // couchbase 4.5.1
			EpTapTotalQueueItemondisk            []float64 `json:"ep_tap_total_queue_itemondisk"`            // couchbase 4.5.1
			EpTapTotalTotalBacklogSize           []float64 `json:"ep_tap_total_total_backlog_size"`          // couchbase 4.5.1
			EpTapUserCount                       []float64 `json:"ep_tap_user_count"`                        // couchbase 4.5.1
			EpTapUserQlen                        []float64 `json:"ep_tap_user_qlen"`                         // couchbase 4.5.1
			EpTapUserQueueBackfillremaining      []float64 `json:"ep_tap_user_queue_backfillremaining"`      // couchbase 4.5.1
			EpTapUserQueueBackoff                []float64 `json:"ep_tap_user_queue_backoff"`                // couchbase 4.5.1
			EpTapUserQueueDrain                  []float64 `json:"ep_tap_user_queue_drain"`                  // couchbase 4.5.1
			EpTapUserQueueFill                   []float64 `json:"ep_tap_user_queue_fill"`                   // couchbase 4.5.1
			EpTapUserQueueItemondisk             []float64 `json:"ep_tap_user_queue_itemondisk"`             // couchbase 4.5.1
			EpTapUserTotalBacklogSize            []float64 `json:"ep_tap_user_total_backlog_size"`           // couchbase 4.5.1
			EpTmpOomErrors                       []float64 `json:"ep_tmp_oom_errors"`
			EpVbTotal                            []float64 `json:"ep_vb_total"`
			Evictions                            []float64 `json:"evictions"`
			GetHits                              []float64 `json:"get_hits"`
			GetMisses                            []float64 `json:"get_misses"`
			IncrHits                             []float64 `json:"incr_hits"`
			IncrMisses                           []float64 `json:"incr_misses"`
			MemUsed                              []float64 `json:"mem_used"`
			Misses                               []float64 `json:"misses"`
			Ops                                  []float64 `json:"ops"`
			VbActiveEject                        []float64 `json:"vb_active_eject"`
			VbActiveItmMemory                    []float64 `json:"vb_active_itm_memory"`
			VbActiveMetaDataMemory               []float64 `json:"vb_active_meta_data_memory"`
			VbActiveNum                          []float64 `json:"vb_active_num"`
			VbActiveNumNonResident               []float64 `json:"vb_active_num_non_resident"`
			VbActiveOpsCreate                    []float64 `json:"vb_active_ops_create"`
			VbActiveOpsUpdate                    []float64 `json:"vb_active_ops_update"`
			VbActiveQueueAge                     []float64 `json:"vb_active_queue_age"`
			VbActiveQueueDrain                   []float64 `json:"vb_active_queue_drain"`
			VbActiveQueueFill                    []float64 `json:"vb_active_queue_fill"`
			VbActiveQueueSize                    []float64 `json:"vb_active_queue_size"`
			VbPendingCurrItems                   []float64 `json:"vb_pending_curr_items"`
			VbPendingEject                       []float64 `json:"vb_pending_eject"`
			VbPendingItmMemory                   []float64 `json:"vb_pending_itm_memory"`
			VbPendingMetaDataMemory              []float64 `json:"vb_pending_meta_data_memory"`
			VbPendingNum                         []float64 `json:"vb_pending_num"`
			VbPendingNumNonResident              []float64 `json:"vb_pending_num_non_resident"`
			VbPendingOpsCreate                   []float64 `json:"vb_pending_ops_create"`
			VbPendingOpsUpdate                   []float64 `json:"vb_pending_ops_update"`
			VbPendingQueueAge                    []float64 `json:"vb_pending_queue_age"`
			VbPendingQueueDrain                  []float64 `json:"vb_pending_queue_drain"`
			VbPendingQueueFill                   []float64 `json:"vb_pending_queue_fill"`
			VbPendingQueueSize                   []float64 `json:"vb_pending_queue_size"`
			VbReplicaCurrItems                   []float64 `json:"vb_replica_curr_items"`
			VbReplicaEject                       []float64 `json:"vb_replica_eject"`
			VbReplicaItmMemory                   []float64 `json:"vb_replica_itm_memory"`
			VbReplicaMetaDataMemory              []float64 `json:"vb_replica_meta_data_memory"`
			VbReplicaNum                         []float64 `json:"vb_replica_num"`
			VbReplicaNumNonResident              []float64 `json:"vb_replica_num_non_resident"`
			VbReplicaOpsCreate                   []float64 `json:"vb_replica_ops_create"`
			VbReplicaOpsUpdate                   []float64 `json:"vb_replica_ops_update"`
			VbReplicaQueueAge                    []float64 `json:"vb_replica_queue_age"`
			VbReplicaQueueDrain                  []float64 `json:"vb_replica_queue_drain"`
			VbReplicaQueueFill                   []float64 `json:"vb_replica_queue_fill"`
			VbReplicaQueueSize                   []float64 `json:"vb_replica_queue_size"`
			VbTotalQueueAge                      []float64 `json:"vb_total_queue_age"`
			XdcOps                               []float64 `json:"xdc_ops"`
			CPUIdleMs                            []float64 `json:"cpu_idle_ms"`
			CPULocalMs                           []float64 `json:"cpu_local_ms"`
			CPUUtilizationRate                   []float64 `json:"cpu_utilization_rate"`
			HibernatedRequests                   []float64 `json:"hibernated_requests"`
			HibernatedWaked                      []float64 `json:"hibernated_waked"`
			MemActualFree                        []float64 `json:"mem_actual_free"`
			MemActualUsed                        []float64 `json:"mem_actual_used"`
			MemFree                              []float64 `json:"mem_free"`
			MemTotal                             []float64 `json:"mem_total"`
			MemUsedSys                           []float64 `json:"mem_used_sys"`
			RestRequests                         []float64 `json:"rest_requests"`
			SwapTotal                            []float64 `json:"swap_total"`
			SwapUsed                             []float64 `json:"swap_used"`
		} `json:"samples"`
	} `json:"op"`
}

// BucketName (/pools/default/buckets)
type BucketName struct {
	Name string `json:"name"`
}

// BucketStatsExporter describes the exporter object.
type BucketStatsExporter struct {
	context                                    Context
	bucketCouchTotalDiskSize                   *p.GaugeVec
	bucketCouchDocsFragmentation               *p.GaugeVec
	bucketCouchViewsFragmentation              *p.GaugeVec
	bucketHitRatio                             *p.GaugeVec
	bucketEpCacheMissRate                      *p.GaugeVec
	bucketEpResidentItemsRate                  *p.GaugeVec
	bucketVbAvgActiveQueueAge                  *p.GaugeVec
	bucketVbAvgReplicaQueueAge                 *p.GaugeVec
	bucketVbAvgPendingQueueAge                 *p.GaugeVec
	bucketVbAvgTotalQueueAge                   *p.GaugeVec
	bucketVbActiveResidentItemsRatio           *p.GaugeVec
	bucketVbReplicaResidentItemsRatio          *p.GaugeVec
	bucketVbPendingResidentItemsRatio          *p.GaugeVec
	bucketAvgDiskUpdateTime                    *p.GaugeVec
	bucketAvgDiskCommitTime                    *p.GaugeVec
	bucketAvgBgWaitTime                        *p.GaugeVec
	bucketAvgActiveTimestampDrift              *p.GaugeVec
	bucketAvgReplicaTimestampDrift             *p.GaugeVec
	bucketEpDcpViewsIndexesCount               *p.GaugeVec
	bucketEpDcpViewsIndexesItemsRemaining      *p.GaugeVec
	bucketEpDcpViewsIndexesProducerCount       *p.GaugeVec
	bucketEpDcpViewsIndexesTotalBacklogSize    *p.GaugeVec
	bucketEpDcpViewsIndexesItemsSent           *p.GaugeVec
	bucketEpDcpViewsIndexesTotalBytes          *p.GaugeVec
	bucketEpDcpViewsIndexesBackoff             *p.GaugeVec
	bucketBgWaitCount                          *p.GaugeVec
	bucketBgWaitTotal                          *p.GaugeVec
	bucketBytesRead                            *p.GaugeVec
	bucketBytesWritten                         *p.GaugeVec
	bucketCasBadval                            *p.GaugeVec
	bucketCasHits                              *p.GaugeVec
	bucketCasMisses                            *p.GaugeVec
	bucketCmdGet                               *p.GaugeVec
	bucketCmdSet                               *p.GaugeVec
	bucketCouchDocsActualDiskSize              *p.GaugeVec
	bucketCouchDocsDataSize                    *p.GaugeVec
	bucketCouchDocsDiskSize                    *p.GaugeVec
	bucketCouchSpatialDataSize                 *p.GaugeVec
	bucketCouchSpatialDiskSize                 *p.GaugeVec
	bucketCouchSpatialOps                      *p.GaugeVec
	bucketCouchViewsActualDiskSize             *p.GaugeVec
	bucketCouchViewsDataSize                   *p.GaugeVec
	bucketCouchViewsDiskSize                   *p.GaugeVec
	bucketCouchViewsOps                        *p.GaugeVec
	bucketCurrConnections                      *p.GaugeVec
	bucketCurrItems                            *p.GaugeVec
	bucketCurrItemsTot                         *p.GaugeVec
	bucketDecrHits                             *p.GaugeVec
	bucketDecrMisses                           *p.GaugeVec
	bucketDeleteHits                           *p.GaugeVec
	bucketDeleteMisses                         *p.GaugeVec
	bucketDiskCommitCount                      *p.GaugeVec
	bucketDiskCommitTotal                      *p.GaugeVec
	bucketDiskUpdateCount                      *p.GaugeVec
	bucketDiskUpdateTotal                      *p.GaugeVec
	bucketDiskWriteQueue                       *p.GaugeVec
	bucketEpActiveAheadExceptions              *p.GaugeVec
	bucketEpActiveHlcDrift                     *p.GaugeVec
	bucketEpActiveHlcDriftCount                *p.GaugeVec
	bucketEpBgFetched                          *p.GaugeVec
	bucketEpClockCasDriftThresholdExceeded     *p.GaugeVec
	bucketEpDcp2IBackoff                       *p.GaugeVec
	bucketEpDcp2ICount                         *p.GaugeVec
	bucketEpDcp2IItemsRemaining                *p.GaugeVec
	bucketEpDcp2IItemsSent                     *p.GaugeVec
	bucketEpDcp2IProducerCount                 *p.GaugeVec
	bucketEpDcp2ITotalBacklogSize              *p.GaugeVec
	bucketEpDcp2ITotalBytes                    *p.GaugeVec
	bucketEpDcpFtsBackoff                      *p.GaugeVec
	bucketEpDcpFtsCount                        *p.GaugeVec
	bucketEpDcpFtsItemsRemaining               *p.GaugeVec
	bucketEpDcpFtsItemsSent                    *p.GaugeVec
	bucketEpDcpFtsProducerCount                *p.GaugeVec
	bucketEpDcpFtsTotalBacklogSize             *p.GaugeVec
	bucketEpDcpFtsTotalBytes                   *p.GaugeVec
	bucketEpDcpOtherBackoff                    *p.GaugeVec
	bucketEpDcpOtherCount                      *p.GaugeVec
	bucketEpDcpOtherItemsRemaining             *p.GaugeVec
	bucketEpDcpOtherItemsSent                  *p.GaugeVec
	bucketEpDcpOtherProducerCount              *p.GaugeVec
	bucketEpDcpOtherTotalBacklogSize           *p.GaugeVec
	bucketEpDcpOtherTotalBytes                 *p.GaugeVec
	bucketEpDcpReplicaBackoff                  *p.GaugeVec
	bucketEpDcpReplicaCount                    *p.GaugeVec
	bucketEpDcpReplicaItemsRemaining           *p.GaugeVec
	bucketEpDcpReplicaItemsSent                *p.GaugeVec
	bucketEpDcpReplicaProducerCount            *p.GaugeVec
	bucketEpDcpReplicaTotalBacklogSize         *p.GaugeVec
	bucketEpDcpReplicaTotalBytes               *p.GaugeVec
	bucketEpDcpViewsBackoff                    *p.GaugeVec
	bucketEpDcpViewsCount                      *p.GaugeVec
	bucketEpDcpViewsItemsRemaining             *p.GaugeVec
	bucketEpDcpViewsItemsSent                  *p.GaugeVec
	bucketEpDcpViewsProducerCount              *p.GaugeVec
	bucketEpDcpViewsTotalBacklogSize           *p.GaugeVec
	bucketEpDcpViewsTotalBytes                 *p.GaugeVec
	bucketEpDcpXdcrBackoff                     *p.GaugeVec
	bucketEpDcpXdcrCount                       *p.GaugeVec
	bucketEpDcpXdcrItemsRemaining              *p.GaugeVec
	bucketEpDcpXdcrItemsSent                   *p.GaugeVec
	bucketEpDcpXdcrProducerCount               *p.GaugeVec
	bucketEpDcpXdcrTotalBacklogSize            *p.GaugeVec
	bucketEpDcpXdcrTotalBytes                  *p.GaugeVec
	bucketEpDiskqueueDrain                     *p.GaugeVec
	bucketEpDiskqueueFill                      *p.GaugeVec
	bucketEpDiskqueueItems                     *p.GaugeVec
	bucketEpFlusherTodo                        *p.GaugeVec
	bucketEpItemCommitFailed                   *p.GaugeVec
	bucketEpKvSize                             *p.GaugeVec
	bucketEpMaxSize                            *p.GaugeVec
	bucketEpMemHighWat                         *p.GaugeVec
	bucketEpMemLowWat                          *p.GaugeVec
	bucketEpMetaDataMemory                     *p.GaugeVec
	bucketEpNumNonResident                     *p.GaugeVec
	bucketEpNumOpsDelMeta                      *p.GaugeVec
	bucketEpNumOpsDelRetMeta                   *p.GaugeVec
	bucketEpNumOpsGetMeta                      *p.GaugeVec
	bucketEpNumOpsSetMeta                      *p.GaugeVec
	bucketEpNumOpsSetRetMeta                   *p.GaugeVec
	bucketEpNumValueEjects                     *p.GaugeVec
	bucketEpOomErrors                          *p.GaugeVec
	bucketEpOpsCreate                          *p.GaugeVec
	bucketEpOpsUpdate                          *p.GaugeVec
	bucketEpOverhead                           *p.GaugeVec
	bucketEpQueueSize                          *p.GaugeVec
	bucketEpReplicaAheadExceptions             *p.GaugeVec
	bucketEpReplicaHlcDrift                    *p.GaugeVec
	bucketEpReplicaHlcDriftCount               *p.GaugeVec
	bucketEpTapRebalanceCount                  *p.GaugeVec
	bucketEpTapRebalanceQlen                   *p.GaugeVec
	bucketEpTapRebalanceQueueBackfillremaining *p.GaugeVec
	bucketEpTapRebalanceQueueBackoff           *p.GaugeVec
	bucketEpTapRebalanceQueueDrain             *p.GaugeVec
	bucketEpTapRebalanceQueueFill              *p.GaugeVec
	bucketEpTapRebalanceQueueItemondisk        *p.GaugeVec
	bucketEpTapRebalanceTotalBacklogSize       *p.GaugeVec
	bucketEpTapReplicaCount                    *p.GaugeVec
	bucketEpTapReplicaQlen                     *p.GaugeVec
	bucketEpTapReplicaQueueBackfillremaining   *p.GaugeVec
	bucketEpTapReplicaQueueBackoff             *p.GaugeVec
	bucketEpTapReplicaQueueDrain               *p.GaugeVec
	bucketEpTapReplicaQueueFill                *p.GaugeVec
	bucketEpTapReplicaQueueItemondisk          *p.GaugeVec
	bucketEpTapReplicaTotalBacklogSize         *p.GaugeVec
	bucketEpTapTotalCount                      *p.GaugeVec
	bucketEpTapTotalQlen                       *p.GaugeVec
	bucketEpTapTotalQueueBackfillremaining     *p.GaugeVec
	bucketEpTapTotalQueueBackoff               *p.GaugeVec
	bucketEpTapTotalQueueDrain                 *p.GaugeVec
	bucketEpTapTotalQueueFill                  *p.GaugeVec
	bucketEpTapTotalQueueItemondisk            *p.GaugeVec
	bucketEpTapTotalTotalBacklogSize           *p.GaugeVec
	bucketEpTapUserCount                       *p.GaugeVec
	bucketEpTapUserQlen                        *p.GaugeVec
	bucketEpTapUserQueueBackfillremaining      *p.GaugeVec
	bucketEpTapUserQueueBackoff                *p.GaugeVec
	bucketEpTapUserQueueDrain                  *p.GaugeVec
	bucketEpTapUserQueueFill                   *p.GaugeVec
	bucketEpTapUserQueueItemondisk             *p.GaugeVec
	bucketEpTapUserTotalBacklogSize            *p.GaugeVec
	bucketEpTmpOomErrors                       *p.GaugeVec
	bucketEpVbTotal                            *p.GaugeVec
	bucketEvictions                            *p.GaugeVec
	bucketGetHits                              *p.GaugeVec
	bucketGetMisses                            *p.GaugeVec
	bucketIncrHits                             *p.GaugeVec
	bucketIncrMisses                           *p.GaugeVec
	bucketMemUsed                              *p.GaugeVec
	bucketMisses                               *p.GaugeVec
	bucketOps                                  *p.GaugeVec
	bucketVbActiveEject                        *p.GaugeVec
	bucketVbActiveItmMemory                    *p.GaugeVec
	bucketVbActiveMetaDataMemory               *p.GaugeVec
	bucketVbActiveNum                          *p.GaugeVec
	bucketVbActiveNumNonResident               *p.GaugeVec
	bucketVbActiveOpsCreate                    *p.GaugeVec
	bucketVbActiveOpsUpdate                    *p.GaugeVec
	bucketVbActiveQueueAge                     *p.GaugeVec
	bucketVbActiveQueueDrain                   *p.GaugeVec
	bucketVbActiveQueueFill                    *p.GaugeVec
	bucketVbActiveQueueSize                    *p.GaugeVec
	bucketVbPendingCurrItems                   *p.GaugeVec
	bucketVbPendingEject                       *p.GaugeVec
	bucketVbPendingItmMemory                   *p.GaugeVec
	bucketVbPendingMetaDataMemory              *p.GaugeVec
	bucketVbPendingNum                         *p.GaugeVec
	bucketVbPendingNumNonResident              *p.GaugeVec
	bucketVbPendingOpsCreate                   *p.GaugeVec
	bucketVbPendingOpsUpdate                   *p.GaugeVec
	bucketVbPendingQueueAge                    *p.GaugeVec
	bucketVbPendingQueueDrain                  *p.GaugeVec
	bucketVbPendingQueueFill                   *p.GaugeVec
	bucketVbPendingQueueSize                   *p.GaugeVec
	bucketVbReplicaCurrItems                   *p.GaugeVec
	bucketVbReplicaEject                       *p.GaugeVec
	bucketVbReplicaItmMemory                   *p.GaugeVec
	bucketVbReplicaMetaDataMemory              *p.GaugeVec
	bucketVbReplicaNum                         *p.GaugeVec
	bucketVbReplicaNumNonResident              *p.GaugeVec
	bucketVbReplicaOpsCreate                   *p.GaugeVec
	bucketVbReplicaOpsUpdate                   *p.GaugeVec
	bucketVbReplicaQueueAge                    *p.GaugeVec
	bucketVbReplicaQueueDrain                  *p.GaugeVec
	bucketVbReplicaQueueFill                   *p.GaugeVec
	bucketVbReplicaQueueSize                   *p.GaugeVec
	bucketVbTotalQueueAge                      *p.GaugeVec
	bucketXdcOps                               *p.GaugeVec
	bucketCPUIdleMs                            *p.GaugeVec
	bucketCPULocalMs                           *p.GaugeVec
	bucketCPUUtilizationRate                   *p.GaugeVec
	bucketHibernatedRequests                   *p.GaugeVec
	bucketHibernatedWaked                      *p.GaugeVec
	bucketMemActualFree                        *p.GaugeVec
	bucketMemActualUsed                        *p.GaugeVec
	bucketMemFree                              *p.GaugeVec
	bucketMemTotal                             *p.GaugeVec
	bucketMemUsedSys                           *p.GaugeVec
	bucketRestRequests                         *p.GaugeVec
	bucketSwapTotal                            *p.GaugeVec
	bucketSwapUsed                             *p.GaugeVec
}

// NewBucketStatsExporter instantiates the Exporter with the URI and metrics.
func NewBucketStatsExporter(context Context) (*BucketStatsExporter, error) {
	exporter := &BucketStatsExporter{
		context:                                 context,
		bucketCouchTotalDiskSize:                newGaugeVec("bucket_couch_total_disk_size", "Couchbase total disk size", []string{"bucket"}),
		bucketCouchDocsFragmentation:            newGaugeVec("bucket_couch_docs_fragmentation", "Couchbase documents fragmentation", []string{"bucket"}),
		bucketCouchViewsFragmentation:           newGaugeVec("bucket_couch_views_fragmentation", "Couchbase views fragmentation", []string{"bucket"}),
		bucketHitRatio:                          newGaugeVec("bucket_hit_ratio", "Hit ratio", []string{"bucket"}),
		bucketEpCacheMissRate:                   newGaugeVec("bucket_ep_cache_miss_rate", "Cache miss rate", []string{"bucket"}),
		bucketEpResidentItemsRate:               newGaugeVec("bucket_ep_resident_items_rate", "Number of resident items", []string{"bucket"}),
		bucketVbAvgActiveQueueAge:               newGaugeVec("bucket_vb_avg_active_queue_age", "Average age in seconds of active items in the active item queue", []string{"bucket"}),
		bucketVbAvgReplicaQueueAge:              newGaugeVec("bucket_vb_avg_replica_queue_age", "Average age in seconds of replica items in the replica item queue", []string{"bucket"}),
		bucketVbAvgPendingQueueAge:              newGaugeVec("bucket_vb_avg_pending_queue_age", "Average age in seconds of pending items in the pending item queue", []string{"bucket"}),
		bucketVbAvgTotalQueueAge:                newGaugeVec("bucket_vb_avg_total_queue_age", "Average age of items in the queue", []string{"bucket"}),
		bucketVbActiveResidentItemsRatio:        newGaugeVec("bucket_vb_active_resident_items_ratio", "Number of resident items", []string{"bucket"}),
		bucketVbReplicaResidentItemsRatio:       newGaugeVec("bucket_vb_replica_resident_items_ratio", "Number of resident replica items", []string{"bucket"}),
		bucketVbPendingResidentItemsRatio:       newGaugeVec("bucket_vb_pending_resident_items_ratio", "Number of resident pending items", []string{"bucket"}),
		bucketAvgDiskUpdateTime:                 newGaugeVec("bucket_avg_disk_update_time", "Average disk update time", []string{"bucket"}),
		bucketAvgDiskCommitTime:                 newGaugeVec("bucket_avg_disk_commit_time", "Average disk commit time", []string{"bucket"}),
		bucketAvgBgWaitTime:                     newGaugeVec("bucket_avg_bg_wait_time", "Average background wait time", []string{"bucket"}),
		bucketEpDcpViewsIndexesCount:            newGaugeVec("bucket_ep_dcp_views_indexes_count", "Number of indexes views DCP connections", []string{"bucket"}),
		bucketEpDcpViewsIndexesItemsRemaining:   newGaugeVec("bucket_ep_dcp_views_indexes_items_remaining", "Number of indexes views items remaining to be sent", []string{"bucket"}),
		bucketEpDcpViewsIndexesProducerCount:    newGaugeVec("bucket_ep_dcp_views_indexes_producer_count", "Number of indexes views producers", []string{"bucket"}),
		bucketEpDcpViewsIndexesTotalBacklogSize: newGaugeVec("bucket_ep_dcp_views_indexes_total_backlog_size", "Number of indexes views items remaining for replication", []string{"bucket"}),
		bucketEpDcpViewsIndexesItemsSent:        newGaugeVec("bucket_ep_dcp_views_indexes_items_sent", "Number of indexes views sent", []string{"bucket"}),
		bucketEpDcpViewsIndexesTotalBytes:       newGaugeVec("bucket_ep_dcp_views_indexes_total_bytes", "Number of bytes per second being sent for indexes views DCP connections", []string{"bucket"}),
		bucketEpDcpViewsIndexesBackoff:          newGaugeVec("bucket_ep_dcp_views_indexes_backoff", "Number of backoffs for indexes views DCP connections", []string{"bucket"}),
		bucketBgWaitCount:                       newGaugeVec("bucket_bg_wait_count", "Background wait", []string{"bucket"}),
		bucketBgWaitTotal:                       newGaugeVec("bucket_bg_wait_total", "Total background wait", []string{"bucket"}),
		bucketBytesRead:                         newGaugeVec("bucket_bytes_read", "Bytes read", []string{"bucket"}),
		bucketBytesWritten:                      newGaugeVec("bucket_bytes_written", "Bytes written", []string{"bucket"}),
		bucketCasBadval:                         newGaugeVec("bucket_cas_badval", "Compare and Swap bad values", []string{"bucket"}),
		bucketCasHits:                           newGaugeVec("bucket_cas_hits", "Compare and Swap hits", []string{"bucket"}),
		bucketCasMisses:                         newGaugeVec("bucket_cas_misses", "Compare and Swap misses", []string{"bucket"}),
		bucketCmdGet:                            newGaugeVec("bucket_cmd_get", "Gets from memory", []string{"bucket"}),
		bucketCmdSet:                            newGaugeVec("bucket_cmd_set", "Sets to memory", []string{"bucket"}),
		bucketCouchDocsActualDiskSize:           newGaugeVec("bucket_couch_docs_actual_disk_size", "Total size of documents on disk in bytes", []string{"bucket"}),
		bucketCouchDocsDataSize:                 newGaugeVec("bucket_couch_docs_data_size", "Documents size in bytes", []string{"bucket"}),
		bucketCouchDocsDiskSize:                 newGaugeVec("bucket_couch_docs_disk_size", "Total size of documents in bytes", []string{"bucket"}),
		bucketCouchSpatialDataSize:              newGaugeVec("bucket_couch_spatial_data_size", "Size of object data for spatial views", []string{"bucket"}),
		bucketCouchSpatialDiskSize:              newGaugeVec("bucket_couch_spatial_disk_size", "Amount of disk space occupied by spatial views", []string{"bucket"}),
		bucketCouchSpatialOps:                   newGaugeVec("bucket_couch_spatial_ops", "Spatial operations", []string{"bucket"}),
		bucketCouchViewsActualDiskSize:          newGaugeVec("bucket_couch_views_actual_disk_size", "Total size of views on disk in bytes", []string{"bucket"}),
		bucketCouchViewsDataSize:                newGaugeVec("bucket_couch_views_data_size", "Views size in bytes", []string{"bucket"}),
		bucketCouchViewsDiskSize:                newGaugeVec("bucket_couch_views_disk_size", "Total size of views in bytes", []string{"bucket"}),
		bucketCouchViewsOps:                     newGaugeVec("bucket_couch_views_ops", "View operations", []string{"bucket"}),
		bucketCurrConnections:                   newGaugeVec("bucket_curr_connections", "Current bucket connections", []string{"bucket"}),
		bucketCurrItems:                         newGaugeVec("bucket_curr_items", "Number of active items in memory", []string{"bucket"}),
		bucketCurrItemsTot:                      newGaugeVec("bucket_curr_items_tot", "Total number of items", []string{"bucket"}),
		bucketDecrHits:                          newGaugeVec("bucket_decr_hits", "Decrement hits", []string{"bucket"}),
		bucketDecrMisses:                        newGaugeVec("bucket_decr_misses", "Decrement misses", []string{"bucket"}),
		bucketDeleteHits:                        newGaugeVec("bucket_delete_hits", "Delete hits", []string{"bucket"}),
		bucketDeleteMisses:                      newGaugeVec("bucket_delete_misses", "Delete misses", []string{"bucket"}),
		bucketDiskCommitCount:                   newGaugeVec("bucket_disk_commit_count", "Disk commits", []string{"bucket"}),
		bucketDiskCommitTotal:                   newGaugeVec("bucket_disk_commit_total", "Total disk commits", []string{"bucket"}),
		bucketDiskUpdateCount:                   newGaugeVec("bucket_disk_update_count", "Disk updates", []string{"bucket"}),
		bucketDiskUpdateTotal:                   newGaugeVec("bucket_disk_update_total", "Total disk updates", []string{"bucket"}),
		bucketDiskWriteQueue:                    newGaugeVec("bucket_disk_write_queue", "Disk write queue depth", []string{"bucket"}),
		bucketEpBgFetched:                       newGaugeVec("bucket_ep_bg_fetched", "Disk reads per second", []string{"bucket"}),
		bucketEpDcp2IBackoff:                    newGaugeVec("bucket_ep_dcp_2i_backoff", "Number of backoffs for indexes DCP connections", []string{"bucket"}),
		bucketEpDcp2ICount:                      newGaugeVec("bucket_ep_dcp_2i_count", "Number of indexes DCP connections", []string{"bucket"}),
		bucketEpDcp2IItemsRemaining:             newGaugeVec("bucket_ep_dcp_2i_items_remaining", "Number of indexes items remaining to be sent", []string{"bucket"}),
		bucketEpDcp2IItemsSent:                  newGaugeVec("bucket_ep_dcp_2i_items_sent", "Number of indexes items sent", []string{"bucket"}),
		bucketEpDcp2IProducerCount:              newGaugeVec("bucket_ep_dcp_2i_producer_count", "Number of indexes producers", []string{"bucket"}),
		bucketEpDcp2ITotalBacklogSize:           newGaugeVec("bucket_ep_dcp_2i_total_backlog_size", "Number of indexes total backlog size", []string{"bucket"}),
		bucketEpDcp2ITotalBytes:                 newGaugeVec("bucket_ep_dcp_2i_total_bytes", "Number bytes per second being sent for indexes DCP connections", []string{"bucket"}),
		bucketEpDcpFtsBackoff:                   newGaugeVec("bucket_ep_dcp_fts_backoff", "Number of backoffs for fts DCP connections", []string{"bucket"}),
		bucketEpDcpFtsCount:                     newGaugeVec("bucket_ep_dcp_fts_count", "Number of fts DCP connections", []string{"bucket"}),
		bucketEpDcpFtsItemsRemaining:            newGaugeVec("bucket_ep_dcp_fts_items_remaining", "Number of fts items remaining to be sent", []string{"bucket"}),
		bucketEpDcpFtsItemsSent:                 newGaugeVec("bucket_ep_dcp_fts_items_sent", "Number of fts items sent", []string{"bucket"}),
		bucketEpDcpFtsProducerCount:             newGaugeVec("bucket_ep_dcp_fts_producer_count", "Number of fts producers", []string{"bucket"}),
		bucketEpDcpFtsTotalBacklogSize:          newGaugeVec("bucket_ep_dcp_fts_total_backlog_size", "Number of fts total backlog size", []string{"bucket"}),
		bucketEpDcpFtsTotalBytes:                newGaugeVec("bucket_ep_dcp_fts_total_bytes", "Number bytes per second being sent for fts DCP connections", []string{"bucket"}),
		bucketEpDcpOtherBackoff:                 newGaugeVec("bucket_ep_dcp_other_backoff", "Number of backoffs for other DCP connections", []string{"bucket"}),
		bucketEpDcpOtherCount:                   newGaugeVec("bucket_ep_dcp_other_count", "Number of other DCP connections", []string{"bucket"}),
		bucketEpDcpOtherItemsRemaining:          newGaugeVec("bucket_ep_dcp_other_items_remaining", "Number of other items remaining to be sent", []string{"bucket"}),
		bucketEpDcpOtherItemsSent:               newGaugeVec("bucket_ep_dcp_other_items_sent", "Number of other items sent", []string{"bucket"}),
		bucketEpDcpOtherProducerCount:           newGaugeVec("bucket_ep_dcp_other_producer_count", "Number of other producers", []string{"bucket"}),
		bucketEpDcpOtherTotalBacklogSize:        newGaugeVec("bucket_ep_dcp_other_total_backlog_size", "Number of other total backlog size", []string{"bucket"}),
		bucketEpDcpOtherTotalBytes:              newGaugeVec("bucket_ep_dcp_other_total_bytes", "Number bytes per second being sent for other DCP connections", []string{"bucket"}),
		bucketEpDcpReplicaBackoff:               newGaugeVec("bucket_ep_dcp_replica_backoff", "Number of backoffs for replica DCP connections", []string{"bucket"}),
		bucketEpDcpReplicaCount:                 newGaugeVec("bucket_ep_dcp_replica_count", "Number of replica DCP connections", []string{"bucket"}),
		bucketEpDcpReplicaItemsRemaining:        newGaugeVec("bucket_ep_dcp_replica_items_remaining", "Number of replica items remaining to be sent", []string{"bucket"}),
		bucketEpDcpReplicaItemsSent:             newGaugeVec("bucket_ep_dcp_replica_items_sent", "Number of replica items sent", []string{"bucket"}),
		bucketEpDcpReplicaProducerCount:         newGaugeVec("bucket_ep_dcp_replica_producer_count", "Number of replica producers", []string{"bucket"}),
		bucketEpDcpReplicaTotalBacklogSize:      newGaugeVec("bucket_ep_dcp_replica_total_backlog_size", "Number of replica total backlog size", []string{"bucket"}),
		bucketEpDcpReplicaTotalBytes:            newGaugeVec("bucket_ep_dcp_replica_total_bytes", "Number bytes per second being sent for replica DCP connections", []string{"bucket"}),
		bucketEpDcpViewsBackoff:                 newGaugeVec("bucket_ep_dcp_views_backoff", "Number of backoffs for views DCP connections", []string{"bucket"}),
		bucketEpDcpViewsCount:                   newGaugeVec("bucket_ep_dcp_views_count", "Number of views DCP connections", []string{"bucket"}),
		bucketEpDcpViewsItemsRemaining:          newGaugeVec("bucket_ep_dcp_views_items_remaining", "Number of views items remaining to be sent", []string{"bucket"}),
		bucketEpDcpViewsItemsSent:               newGaugeVec("bucket_ep_dcp_views_items_sent", "Number of views items sent", []string{"bucket"}),
		bucketEpDcpViewsProducerCount:           newGaugeVec("bucket_ep_dcp_views_producer_count", "Number of views producers", []string{"bucket"}),
		bucketEpDcpViewsTotalBacklogSize:        newGaugeVec("bucket_ep_dcp_views_total_backlog_size", "Number of views total backlog size", []string{"bucket"}),
		bucketEpDcpViewsTotalBytes:              newGaugeVec("bucket_ep_dcp_views_total_bytes", "Number bytes per second being sent for views DCP connections", []string{"bucket"}),
		bucketEpDcpXdcrBackoff:                  newGaugeVec("bucket_ep_dcp_xdcr_backoff", "Number of backoffs for xdcr DCP connections", []string{"bucket"}),
		bucketEpDcpXdcrCount:                    newGaugeVec("bucket_ep_dcp_xdcr_count", "Number of xdcr DCP connections", []string{"bucket"}),
		bucketEpDcpXdcrItemsRemaining:           newGaugeVec("bucket_ep_dcp_xdcr_items_remaining", "Number of xdcr items remaining to be sent", []string{"bucket"}),
		bucketEpDcpXdcrItemsSent:                newGaugeVec("bucket_ep_dcp_xdcr_items_sent", "Number of xdcr items sent", []string{"bucket"}),
		bucketEpDcpXdcrProducerCount:            newGaugeVec("bucket_ep_dcp_xdcr_producer_count", "Number of xdcr producers", []string{"bucket"}),
		bucketEpDcpXdcrTotalBacklogSize:         newGaugeVec("bucket_ep_dcp_xdcr_total_backlog_size", "Number of xdcr total backlog size", []string{"bucket"}),
		bucketEpDcpXdcrTotalBytes:               newGaugeVec("bucket_ep_dcp_xdcr_total_bytes", "Number bytes per second being sent for xdcr DCP connections", []string{"bucket"}),
		bucketEpDiskqueueDrain:                  newGaugeVec("bucket_ep_diskqueue_drain", "Total Drained items on disk queue", []string{"bucket"}),
		bucketEpDiskqueueFill:                   newGaugeVec("bucket_ep_diskqueue_fill", "Total enqueued items on disk queue", []string{"bucket"}),
		bucketEpDiskqueueItems:                  newGaugeVec("bucket_ep_diskqueue_items", "Total number of items waiting to be written to disk", []string{"bucket"}),
		bucketEpFlusherTodo:                     newGaugeVec("bucket_ep_flusher_todo", "Number of items currently being written", []string{"bucket"}),
		bucketEpItemCommitFailed:                newGaugeVec("bucket_ep_item_commit_failed", "Number of times a transaction failed to commit due to storage errors", []string{"bucket"}),
		bucketEpKvSize:                          newGaugeVec("bucket_ep_kv_size", "Total amount of user data cached in RAM", []string{"bucket"}),
		bucketEpMaxSize:                         newGaugeVec("bucket_ep_max_size", "Maximum amount of memory this bucket can use", []string{"bucket"}),
		bucketEpMemHighWat:                      newGaugeVec("bucket_ep_mem_high_wat", "Memory usage high water mark for auto-evictions", []string{"bucket"}),
		bucketEpMemLowWat:                       newGaugeVec("bucket_ep_mem_low_wat", "Memory usage low water mark for auto-evictions", []string{"bucket"}),
		bucketEpMetaDataMemory:                  newGaugeVec("bucket_ep_meta_data_memory", "Total amount of item metadata consuming RAM", []string{"bucket"}),
		bucketEpNumNonResident:                  newGaugeVec("bucket_ep_num_non_resident", "Number of non-resident items", []string{"bucket"}),
		bucketEpNumOpsDelMeta:                   newGaugeVec("bucket_ep_num_ops_del_meta", "Number of delete operations per second for this bucket as the target for XDCR", []string{"bucket"}),
		bucketEpNumOpsDelRetMeta:                newGaugeVec("bucket_ep_num_ops_del_ret_meta", "Number of delRetMeta operations per second for this bucket as the target for XDCR", []string{"bucket"}),
		bucketEpNumOpsGetMeta:                   newGaugeVec("bucket_ep_num_ops_get_meta", "Number of read operations per second for this bucket as the target for XDCR", []string{"bucket"}),
		bucketEpNumOpsSetMeta:                   newGaugeVec("bucket_ep_num_ops_set_meta", "Number of write operations per second for this bucket as the target for XDCR", []string{"bucket"}),
		bucketEpNumOpsSetRetMeta:                newGaugeVec("bucket_ep_num_ops_set_ret_meta", "Number of setRetMeta operations per second for this bucket as the target for XDCR", []string{"bucket"}),
		bucketEpNumValueEjects:                  newGaugeVec("bucket_ep_num_value_ejects", "Number of times item values got ejected from memory to disk", []string{"bucket"}),
		bucketEpOomErrors:                       newGaugeVec("bucket_ep_oom_errors", "Number of times unrecoverable OOMs happened while processing operations", []string{"bucket"}),
		bucketEpOpsCreate:                       newGaugeVec("bucket_ep_ops_create", "Create operations", []string{"bucket"}),
		bucketEpOpsUpdate:                       newGaugeVec("bucket_ep_ops_update", "Update operations", []string{"bucket"}),
		bucketEpOverhead:                        newGaugeVec("bucket_ep_overhead", "Extra memory used by transient data like persistence queues or checkpoints", []string{"bucket"}),
		bucketEpQueueSize:                       newGaugeVec("bucket_ep_queue_size", "Number of items queued for storage", []string{"bucket"}),
		bucketEpTmpOomErrors:                    newGaugeVec("bucket_ep_tmp_oom_errors", "Number of times recoverable OOMs happened while processing operations", []string{"bucket"}),
		bucketEpVbTotal:                         newGaugeVec("bucket_ep_vb_total", "Total number of vBuckets for this bucket", []string{"bucket"}),
		bucketEvictions:                         newGaugeVec("bucket_evictions", "Number of evictions", []string{"bucket"}),
		bucketGetHits:                           newGaugeVec("bucket_get_hits", "Number of get hits", []string{"bucket"}),
		bucketGetMisses:                         newGaugeVec("bucket_get_misses", "Number of get misses", []string{"bucket"}),
		bucketIncrHits:                          newGaugeVec("bucket_incr_hits", "Number of increment hits", []string{"bucket"}),
		bucketIncrMisses:                        newGaugeVec("bucket_incr_misses", "Number of increment misses", []string{"bucket"}),
		bucketMemUsed:                           newGaugeVec("bucket_mem_used", "Engine's total memory usage (deprecated)", []string{"bucket"}),
		bucketMisses:                            newGaugeVec("bucket_misses", "Total number of misses", []string{"bucket"}),
		bucketOps:                               newGaugeVec("bucket_ops", "Total number of operations", []string{"bucket"}),
		bucketVbActiveEject:                     newGaugeVec("bucket_vb_active_eject", "Number of items per second being ejected to disk from active vBuckets", []string{"bucket"}),
		bucketVbActiveItmMemory:                 newGaugeVec("bucket_vb_active_itm_memory", "Amount of active user data cached in RAM", []string{"bucket"}),
		bucketVbActiveMetaDataMemory:            newGaugeVec("bucket_vb_active_meta_data_memory", "Amount of active item metadata consuming RAM", []string{"bucket"}),
		bucketVbActiveNum:                       newGaugeVec("bucket_vb_active_num", "Number of active items", []string{"bucket"}),
		bucketVbActiveNumNonResident:            newGaugeVec("bucket_vb_active_num_non_resident", "Number of non resident vBuckets in the active state for this bucket", []string{"bucket"}),
		bucketVbActiveOpsCreate:                 newGaugeVec("bucket_vb_active_ops_create", "New items per second being inserted into active vBuckets", []string{"bucket"}),
		bucketVbActiveOpsUpdate:                 newGaugeVec("bucket_vb_active_ops_update", "Number of items updated on active vBucket per second for this bucket", []string{"bucket"}),
		bucketVbActiveQueueAge:                  newGaugeVec("bucket_vb_active_queue_age", "Sum of disk queue item age in milliseconds", []string{"bucket"}),
		bucketVbActiveQueueDrain:                newGaugeVec("bucket_vb_active_queue_drain", "Total drained items in the queue", []string{"bucket"}),
		bucketVbActiveQueueFill:                 newGaugeVec("bucket_vb_active_queue_fill", "Number of active items per second being put on the active item disk queue", []string{"bucket"}),
		bucketVbActiveQueueSize:                 newGaugeVec("bucket_vb_active_queue_size", "Number of active items in the queue", []string{"bucket"}),
		bucketVbPendingCurrItems:                newGaugeVec("bucket_vb_pending_curr_items", "Number of items in pending vBuckets", []string{"bucket"}),
		bucketVbPendingEject:                    newGaugeVec("bucket_vb_pending_eject", "Number of items per second being ejected to disk from pending vBuckets", []string{"bucket"}),
		bucketVbPendingItmMemory:                newGaugeVec("bucket_vb_pending_itm_memory", "Amount of pending user data cached in RAM", []string{"bucket"}),
		bucketVbPendingMetaDataMemory:           newGaugeVec("bucket_vb_pending_meta_data_memory", "Amount of pending item metadata consuming RAM", []string{"bucket"}),
		bucketVbPendingNum:                      newGaugeVec("bucket_vb_pending_num", "Number of pending items", []string{"bucket"}),
		bucketVbPendingNumNonResident:           newGaugeVec("bucket_vb_pending_num_non_resident", "Number of non resident vBuckets in the pending state for this bucket", []string{"bucket"}),
		bucketVbPendingOpsCreate:                newGaugeVec("bucket_vb_pending_ops_create", "Number of pending create operations", []string{"bucket"}),
		bucketVbPendingOpsUpdate:                newGaugeVec("bucket_vb_pending_ops_update", "Number of items updated on pending vBucket per second for this bucket", []string{"bucket"}),
		bucketVbPendingQueueAge:                 newGaugeVec("bucket_vb_pending_queue_age", "Sum of disk pending queue item age in milliseconds", []string{"bucket"}),
		bucketVbPendingQueueDrain:               newGaugeVec("bucket_vb_pending_queue_drain", "Total drained pending items in the queue", []string{"bucket"}),
		bucketVbPendingQueueFill:                newGaugeVec("bucket_vb_pending_queue_fill", "Total enqueued pending items on disk queue", []string{"bucket"}),
		bucketVbPendingQueueSize:                newGaugeVec("bucket_vb_pending_queue_size", "Number of pending items in the queue", []string{"bucket"}),
		bucketVbReplicaCurrItems:                newGaugeVec("bucket_vb_replica_curr_items", "Number of in memory items", []string{"bucket"}),
		bucketVbReplicaEject:                    newGaugeVec("bucket_vb_replica_eject", "Number of items per second being ejected to disk from replica vBuckets", []string{"bucket"}),
		bucketVbReplicaItmMemory:                newGaugeVec("bucket_vb_replica_itm_memory", "Amount of replica user data cached in RAM", []string{"bucket"}),
		bucketVbReplicaMetaDataMemory:           newGaugeVec("bucket_vb_replica_meta_data_memory", "Total metadata memory", []string{"bucket"}),
		bucketVbReplicaNum:                      newGaugeVec("bucket_vb_replica_num", "Number of replica vBuckets", []string{"bucket"}),
		bucketVbReplicaNumNonResident:           newGaugeVec("bucket_vb_replica_num_non_resident", "Number of non resident vBuckets in the replica state for this bucket", []string{"bucket"}),
		bucketVbReplicaOpsCreate:                newGaugeVec("bucket_vb_replica_ops_create", "Number of replica create operations", []string{"bucket"}),
		bucketVbReplicaOpsUpdate:                newGaugeVec("bucket_vb_replica_ops_update", "Number of items updated on replica vBucket per second for this bucket", []string{"bucket"}),
		bucketVbReplicaQueueAge:                 newGaugeVec("bucket_vb_replica_queue_age", "Sum of disk replica queue item age in milliseconds", []string{"bucket"}),
		bucketVbReplicaQueueDrain:               newGaugeVec("bucket_vb_replica_queue_drain", "Total drained replica items in the queue", []string{"bucket"}),
		bucketVbReplicaQueueFill:                newGaugeVec("bucket_vb_replica_queue_fill", "Total enqueued replica items on disk queue", []string{"bucket"}),
		bucketVbReplicaQueueSize:                newGaugeVec("bucket_vb_replica_queue_size", "Replica items in disk queue", []string{"bucket"}),
		bucketVbTotalQueueAge:                   newGaugeVec("bucket_vb_total_queue_age", "Sum of disk queue item age in milliseconds", []string{"bucket"}),
		bucketXdcOps:                            newGaugeVec("bucket_xdc_ops", "Number of cross-datacenter replication operations", []string{"bucket"}),
		bucketCPUIdleMs:                         newGaugeVec("bucket_cpu_idle_ms", "CPU idle milliseconds", []string{"bucket"}),
		bucketCPULocalMs:                        newGaugeVec("bucket_cpu_local_ms", "CPU local milliseconds", []string{"bucket"}),
		bucketCPUUtilizationRate:                newGaugeVec("bucket_cpu_utilization_rate", "CPU utilization percentage", []string{"bucket"}),
		bucketHibernatedRequests:                newGaugeVec("bucket_hibernated_requests", "Number of streaming requests now idle", []string{"bucket"}),
		bucketHibernatedWaked:                   newGaugeVec("bucket_hibernated_waked", "Rate of streaming request wakeups", []string{"bucket"}),
		bucketMemActualFree:                     newGaugeVec("bucket_mem_actual_free", "Actual free memory", []string{"bucket"}),
		bucketMemActualUsed:                     newGaugeVec("bucket_mem_actual_used", "Actual used memory", []string{"bucket"}),
		bucketMemFree:                           newGaugeVec("bucket_mem_free", "Free memory", []string{"bucket"}),
		bucketMemTotal:                          newGaugeVec("bucket_mem_total", "Total memeory", []string{"bucket"}),
		bucketMemUsedSys:                        newGaugeVec("bucket_mem_used_sys", "System memory usage", []string{"bucket"}),
		bucketRestRequests:                      newGaugeVec("bucket_rest_requests", "Number of HTTP requests", []string{"bucket"}),
		bucketSwapTotal:                         newGaugeVec("bucket_swap_total", "Total amount of swap available", []string{"bucket"}),
		bucketSwapUsed:                          newGaugeVec("bucket_swap_used", "Amount of swap used", []string{"bucket"}),
	}

	if context.CouchbaseVersion == "4.5.1" {
		exporter.bucketEpTapRebalanceCount = newGaugeVec("bucket_ep_tap_rebalance_count", "Number of internal rebalancing TAP queues", []string{"bucket"})
		exporter.bucketEpTapRebalanceQlen = newGaugeVec("bucket_ep_tap_rebalance_qlen", "Number of items in the rebalance TAP queues", []string{"bucket"})
		exporter.bucketEpTapRebalanceQueueBackfillremaining = newGaugeVec("bucket_ep_tap_rebalance_queue_backfillremaining", "Number of items in the backfill queues of rebalancing TAP connections", []string{"bucket"})
		exporter.bucketEpTapRebalanceQueueBackoff = newGaugeVec("bucket_ep_tap_rebalance_queue_backoff", "Number of back-offs received per second while sending data over rebalancing TAP connections", []string{"bucket"})
		exporter.bucketEpTapRebalanceQueueDrain = newGaugeVec("bucket_ep_tap_rebalance_queue_drain", "Number of items per second being sent over rebalancing TAP connections, i.e. removed from queue", []string{"bucket"})
		exporter.bucketEpTapRebalanceQueueFill = newGaugeVec("bucket_ep_tap_rebalance_queue_fill", "Number of items per second being sent to queue", []string{"bucket"})
		exporter.bucketEpTapRebalanceQueueItemondisk = newGaugeVec("bucket_ep_tap_rebalance_queue_itemondisk", "Number of items still on disk to be loaded for rebalancing TAP connections", []string{"bucket"})
		exporter.bucketEpTapRebalanceTotalBacklogSize = newGaugeVec("bucket_ep_tap_rebalance_total_backlog_size", "Number of remaining items for rebalancing TAP connections", []string{"bucket"})
		exporter.bucketEpTapReplicaCount = newGaugeVec("bucket_ep_tap_replica_count", "Number of internal replication TAP queues", []string{"bucket"})
		exporter.bucketEpTapReplicaQlen = newGaugeVec("bucket_ep_tap_replica_qlen", "Number of items in the replication TAP queues", []string{"bucket"})
		exporter.bucketEpTapReplicaQueueBackfillremaining = newGaugeVec("bucket_ep_tap_replica_queue_backfillremaining", "Number of items in the backfill queues of replication TAP connections", []string{"bucket"})
		exporter.bucketEpTapReplicaQueueBackoff = newGaugeVec("bucket_ep_tap_replica_queue_backoff", "Number of back-offs received per second while sending data over replication TAP connections", []string{"bucket"})
		exporter.bucketEpTapReplicaQueueDrain = newGaugeVec("bucket_ep_tap_replica_queue_drain", "Total drained items in the replica queue", []string{"bucket"})
		exporter.bucketEpTapReplicaQueueFill = newGaugeVec("bucket_ep_tap_replica_queue_fill", "Number of items per second being sent to queue", []string{"bucket"})
		exporter.bucketEpTapReplicaQueueItemondisk = newGaugeVec("bucket_ep_tap_replica_queue_itemondisk", "Number of items still on disk to be loaded for replication TAP connections", []string{"bucket"})
		exporter.bucketEpTapReplicaTotalBacklogSize = newGaugeVec("bucket_ep_tap_replica_total_backlog_size", "Number of remaining items for replication TAP connections", []string{"bucket"})
		exporter.bucketEpTapTotalCount = newGaugeVec("bucket_ep_tap_total_count", "Total number of internal TAP queues", []string{"bucket"})
		exporter.bucketEpTapTotalQlen = newGaugeVec("bucket_ep_tap_total_qlen", "Total number of items in TAP queues", []string{"bucket"})
		exporter.bucketEpTapTotalQueueBackfillremaining = newGaugeVec("bucket_ep_tap_total_queue_backfillremaining", "Total number of items in the backfill queues of TAP connections", []string{"bucket"})
		exporter.bucketEpTapTotalQueueBackoff = newGaugeVec("bucket_ep_tap_total_queue_backoff", "Total number of back-offs received per second while sending data over TAP connections", []string{"bucket"})
		exporter.bucketEpTapTotalQueueDrain = newGaugeVec("bucket_ep_tap_total_queue_drain", "Total drained items in the queue", []string{"bucket"})
		exporter.bucketEpTapTotalQueueFill = newGaugeVec("bucket_ep_tap_total_queue_fill", "Total enqueued items in the queue", []string{"bucket"})
		exporter.bucketEpTapTotalQueueItemondisk = newGaugeVec("bucket_ep_tap_total_queue_itemondisk", "Total number of items still on disk to be loaded for TAP connections", []string{"bucket"})
		exporter.bucketEpTapTotalTotalBacklogSize = newGaugeVec("bucket_ep_tap_total_total_backlog_size", "Number of remaining items for replication", []string{"bucket"})
		exporter.bucketEpTapUserCount = newGaugeVec("bucket_ep_tap_user_count", "Number of internal user TAP queues", []string{"bucket"})
		exporter.bucketEpTapUserQlen = newGaugeVec("bucket_ep_tap_user_qlen", "Number of items in user TAP queues", []string{"bucket"})
		exporter.bucketEpTapUserQueueBackfillremaining = newGaugeVec("bucket_ep_tap_user_queue_backfillremaining", "Number of items in the backfill queues of user TAP connections", []string{"bucket"})
		exporter.bucketEpTapUserQueueBackoff = newGaugeVec("bucket_ep_tap_user_queue_backoff", "Number of back-offs received per second while sending data over user TAP connections", []string{"bucket"})
		exporter.bucketEpTapUserQueueDrain = newGaugeVec("bucket_ep_tap_user_queue_drain", "Number of items per second being sent over user TAP connections to this bucket, i.e. removed from queue", []string{"bucket"})
		exporter.bucketEpTapUserQueueFill = newGaugeVec("bucket_ep_tap_user_queue_fill", "Number of items per second being sent to queue", []string{"bucket"})
		exporter.bucketEpTapUserQueueItemondisk = newGaugeVec("bucket_ep_tap_user_queue_itemondisk", "Number of items still on disk to be loaded for client TAP connections", []string{"bucket"})
		exporter.bucketEpTapUserTotalBacklogSize = newGaugeVec("bucket_ep_tap_user_total_backlog_size", "Number of remaining items for client TAP connections", []string{"bucket"})
	}
	if context.CouchbaseVersion == "5.1.1" {
		exporter.bucketAvgActiveTimestampDrift = newGaugeVec("bucket_avg_active_timestamp_drift", "Average active timestamp drift", []string{"bucket"})
		exporter.bucketAvgReplicaTimestampDrift = newGaugeVec("bucket_avg_replica_timestamp_drift", "Average replica timestamp drift", []string{"bucket"})
		exporter.bucketEpActiveAheadExceptions = newGaugeVec("bucket_ep_active_ahead_exceptions", "Sum total of all active vBuckets drift_ahead_threshold_exceeded counter", []string{"bucket"})
		exporter.bucketEpActiveHlcDrift = newGaugeVec("bucket_ep_active_hlc_drift", "Total absolute drift for all active vBuckets", []string{"bucket"})
		exporter.bucketEpActiveHlcDriftCount = newGaugeVec("bucket_ep_active_hlc_drift_count", "Number of updates applied to ep_active_hlc_drift", []string{"bucket"})
		exporter.bucketEpClockCasDriftThresholdExceeded = newGaugeVec("bucket_ep_clock_cas_drift_threshold_exceeded", "Ep clock cas drift threshold exceeded", []string{"bucket"})
		exporter.bucketEpReplicaAheadExceptions = newGaugeVec("bucket_ep_replica_ahead_exceptions", "Sum total of all replica vBuckets' drift_ahead_threshold_exceeded counter", []string{"bucket"})
		exporter.bucketEpReplicaHlcDrift = newGaugeVec("bucket_ep_replica_hlc_drift", "Total abosulte drift for all replica vBuckets", []string{"bucket"})
		exporter.bucketEpReplicaHlcDriftCount = newGaugeVec("bucket_ep_replica_hlc_drift_count", "Number of updates applied to ep_replica_hlc_drift", []string{"bucket"})
	}

	return exporter, nil
}

// Describe describes exported metrics.
func (e *BucketStatsExporter) Describe(ch chan<- *p.Desc) {

	e.bucketCouchTotalDiskSize.Describe(ch)
	e.bucketCouchDocsFragmentation.Describe(ch)
	e.bucketCouchViewsFragmentation.Describe(ch)
	e.bucketHitRatio.Describe(ch)
	e.bucketEpCacheMissRate.Describe(ch)
	e.bucketEpResidentItemsRate.Describe(ch)
	e.bucketVbAvgActiveQueueAge.Describe(ch)
	e.bucketVbAvgReplicaQueueAge.Describe(ch)
	e.bucketVbAvgPendingQueueAge.Describe(ch)
	e.bucketVbAvgTotalQueueAge.Describe(ch)
	e.bucketVbActiveResidentItemsRatio.Describe(ch)
	e.bucketVbReplicaResidentItemsRatio.Describe(ch)
	e.bucketVbPendingResidentItemsRatio.Describe(ch)
	e.bucketAvgDiskUpdateTime.Describe(ch)
	e.bucketAvgDiskCommitTime.Describe(ch)
	e.bucketAvgBgWaitTime.Describe(ch)
	e.bucketEpDcpViewsIndexesCount.Describe(ch)
	e.bucketEpDcpViewsIndexesItemsRemaining.Describe(ch)
	e.bucketEpDcpViewsIndexesProducerCount.Describe(ch)
	e.bucketEpDcpViewsIndexesTotalBacklogSize.Describe(ch)
	e.bucketEpDcpViewsIndexesItemsSent.Describe(ch)
	e.bucketEpDcpViewsIndexesTotalBytes.Describe(ch)
	e.bucketEpDcpViewsIndexesBackoff.Describe(ch)
	e.bucketBgWaitCount.Describe(ch)
	e.bucketBgWaitTotal.Describe(ch)
	e.bucketBytesRead.Describe(ch)
	e.bucketBytesWritten.Describe(ch)
	e.bucketCasBadval.Describe(ch)
	e.bucketCasHits.Describe(ch)
	e.bucketCasMisses.Describe(ch)
	e.bucketCmdGet.Describe(ch)
	e.bucketCmdSet.Describe(ch)
	e.bucketCouchDocsActualDiskSize.Describe(ch)
	e.bucketCouchDocsDataSize.Describe(ch)
	e.bucketCouchDocsDiskSize.Describe(ch)
	e.bucketCouchSpatialDataSize.Describe(ch)
	e.bucketCouchSpatialDiskSize.Describe(ch)
	e.bucketCouchSpatialOps.Describe(ch)
	e.bucketCouchViewsActualDiskSize.Describe(ch)
	e.bucketCouchViewsDataSize.Describe(ch)
	e.bucketCouchViewsDiskSize.Describe(ch)
	e.bucketCouchViewsOps.Describe(ch)
	e.bucketCurrConnections.Describe(ch)
	e.bucketCurrItems.Describe(ch)
	e.bucketCurrItemsTot.Describe(ch)
	e.bucketDecrHits.Describe(ch)
	e.bucketDecrMisses.Describe(ch)
	e.bucketDeleteHits.Describe(ch)
	e.bucketDeleteMisses.Describe(ch)
	e.bucketDiskCommitCount.Describe(ch)
	e.bucketDiskCommitTotal.Describe(ch)
	e.bucketDiskUpdateCount.Describe(ch)
	e.bucketDiskUpdateTotal.Describe(ch)
	e.bucketDiskWriteQueue.Describe(ch)
	e.bucketEpBgFetched.Describe(ch)
	e.bucketEpDcp2IBackoff.Describe(ch)
	e.bucketEpDcp2ICount.Describe(ch)
	e.bucketEpDcp2IItemsRemaining.Describe(ch)
	e.bucketEpDcp2IItemsSent.Describe(ch)
	e.bucketEpDcp2IProducerCount.Describe(ch)
	e.bucketEpDcp2ITotalBacklogSize.Describe(ch)
	e.bucketEpDcp2ITotalBytes.Describe(ch)
	e.bucketEpDcpFtsBackoff.Describe(ch)
	e.bucketEpDcpFtsCount.Describe(ch)
	e.bucketEpDcpFtsItemsRemaining.Describe(ch)
	e.bucketEpDcpFtsItemsSent.Describe(ch)
	e.bucketEpDcpFtsProducerCount.Describe(ch)
	e.bucketEpDcpFtsTotalBacklogSize.Describe(ch)
	e.bucketEpDcpFtsTotalBytes.Describe(ch)
	e.bucketEpDcpOtherBackoff.Describe(ch)
	e.bucketEpDcpOtherCount.Describe(ch)
	e.bucketEpDcpOtherItemsRemaining.Describe(ch)
	e.bucketEpDcpOtherItemsSent.Describe(ch)
	e.bucketEpDcpOtherProducerCount.Describe(ch)
	e.bucketEpDcpOtherTotalBacklogSize.Describe(ch)
	e.bucketEpDcpOtherTotalBytes.Describe(ch)
	e.bucketEpDcpReplicaBackoff.Describe(ch)
	e.bucketEpDcpReplicaCount.Describe(ch)
	e.bucketEpDcpReplicaItemsRemaining.Describe(ch)
	e.bucketEpDcpReplicaItemsSent.Describe(ch)
	e.bucketEpDcpReplicaProducerCount.Describe(ch)
	e.bucketEpDcpReplicaTotalBacklogSize.Describe(ch)
	e.bucketEpDcpReplicaTotalBytes.Describe(ch)
	e.bucketEpDcpViewsBackoff.Describe(ch)
	e.bucketEpDcpViewsCount.Describe(ch)
	e.bucketEpDcpViewsItemsRemaining.Describe(ch)
	e.bucketEpDcpViewsItemsSent.Describe(ch)
	e.bucketEpDcpViewsProducerCount.Describe(ch)
	e.bucketEpDcpViewsTotalBacklogSize.Describe(ch)
	e.bucketEpDcpViewsTotalBytes.Describe(ch)
	e.bucketEpDcpXdcrBackoff.Describe(ch)
	e.bucketEpDcpXdcrCount.Describe(ch)
	e.bucketEpDcpXdcrItemsRemaining.Describe(ch)
	e.bucketEpDcpXdcrItemsSent.Describe(ch)
	e.bucketEpDcpXdcrProducerCount.Describe(ch)
	e.bucketEpDcpXdcrTotalBacklogSize.Describe(ch)
	e.bucketEpDcpXdcrTotalBytes.Describe(ch)
	e.bucketEpDiskqueueDrain.Describe(ch)
	e.bucketEpDiskqueueFill.Describe(ch)
	e.bucketEpDiskqueueItems.Describe(ch)
	e.bucketEpFlusherTodo.Describe(ch)
	e.bucketEpItemCommitFailed.Describe(ch)
	e.bucketEpKvSize.Describe(ch)
	e.bucketEpMaxSize.Describe(ch)
	e.bucketEpMemHighWat.Describe(ch)
	e.bucketEpMemLowWat.Describe(ch)
	e.bucketEpMetaDataMemory.Describe(ch)
	e.bucketEpNumNonResident.Describe(ch)
	e.bucketEpNumOpsDelMeta.Describe(ch)
	e.bucketEpNumOpsDelRetMeta.Describe(ch)
	e.bucketEpNumOpsGetMeta.Describe(ch)
	e.bucketEpNumOpsSetMeta.Describe(ch)
	e.bucketEpNumOpsSetRetMeta.Describe(ch)
	e.bucketEpNumValueEjects.Describe(ch)
	e.bucketEpOomErrors.Describe(ch)
	e.bucketEpOpsCreate.Describe(ch)
	e.bucketEpOpsUpdate.Describe(ch)
	e.bucketEpOverhead.Describe(ch)
	e.bucketEpQueueSize.Describe(ch)
	e.bucketEpTmpOomErrors.Describe(ch)
	e.bucketEpVbTotal.Describe(ch)
	e.bucketEvictions.Describe(ch)
	e.bucketGetHits.Describe(ch)
	e.bucketGetMisses.Describe(ch)
	e.bucketIncrHits.Describe(ch)
	e.bucketIncrMisses.Describe(ch)
	e.bucketMemUsed.Describe(ch)
	e.bucketMisses.Describe(ch)
	e.bucketOps.Describe(ch)
	e.bucketVbActiveEject.Describe(ch)
	e.bucketVbActiveItmMemory.Describe(ch)
	e.bucketVbActiveMetaDataMemory.Describe(ch)
	e.bucketVbActiveNum.Describe(ch)
	e.bucketVbActiveNumNonResident.Describe(ch)
	e.bucketVbActiveOpsCreate.Describe(ch)
	e.bucketVbActiveOpsUpdate.Describe(ch)
	e.bucketVbActiveQueueAge.Describe(ch)
	e.bucketVbActiveQueueDrain.Describe(ch)
	e.bucketVbActiveQueueFill.Describe(ch)
	e.bucketVbActiveQueueSize.Describe(ch)
	e.bucketVbPendingCurrItems.Describe(ch)
	e.bucketVbPendingEject.Describe(ch)
	e.bucketVbPendingItmMemory.Describe(ch)
	e.bucketVbPendingMetaDataMemory.Describe(ch)
	e.bucketVbPendingNum.Describe(ch)
	e.bucketVbPendingNumNonResident.Describe(ch)
	e.bucketVbPendingOpsCreate.Describe(ch)
	e.bucketVbPendingOpsUpdate.Describe(ch)
	e.bucketVbPendingQueueAge.Describe(ch)
	e.bucketVbPendingQueueDrain.Describe(ch)
	e.bucketVbPendingQueueFill.Describe(ch)
	e.bucketVbPendingQueueSize.Describe(ch)
	e.bucketVbReplicaCurrItems.Describe(ch)
	e.bucketVbReplicaEject.Describe(ch)
	e.bucketVbReplicaItmMemory.Describe(ch)
	e.bucketVbReplicaMetaDataMemory.Describe(ch)
	e.bucketVbReplicaNum.Describe(ch)
	e.bucketVbReplicaNumNonResident.Describe(ch)
	e.bucketVbReplicaOpsCreate.Describe(ch)
	e.bucketVbReplicaOpsUpdate.Describe(ch)
	e.bucketVbReplicaQueueAge.Describe(ch)
	e.bucketVbReplicaQueueDrain.Describe(ch)
	e.bucketVbReplicaQueueFill.Describe(ch)
	e.bucketVbReplicaQueueSize.Describe(ch)
	e.bucketVbTotalQueueAge.Describe(ch)
	e.bucketXdcOps.Describe(ch)
	e.bucketCPUIdleMs.Describe(ch)
	e.bucketCPULocalMs.Describe(ch)
	e.bucketCPUUtilizationRate.Describe(ch)
	e.bucketHibernatedRequests.Describe(ch)
	e.bucketHibernatedWaked.Describe(ch)
	e.bucketMemActualFree.Describe(ch)
	e.bucketMemActualUsed.Describe(ch)
	e.bucketMemFree.Describe(ch)
	e.bucketMemTotal.Describe(ch)
	e.bucketMemUsedSys.Describe(ch)
	e.bucketRestRequests.Describe(ch)
	e.bucketSwapTotal.Describe(ch)
	e.bucketSwapUsed.Describe(ch)

	if e.context.CouchbaseVersion == "4.5.1" {
		e.bucketEpTapRebalanceCount.Describe(ch)
		e.bucketEpTapRebalanceQlen.Describe(ch)
		e.bucketEpTapRebalanceQueueBackfillremaining.Describe(ch)
		e.bucketEpTapRebalanceQueueBackoff.Describe(ch)
		e.bucketEpTapRebalanceQueueDrain.Describe(ch)
		e.bucketEpTapRebalanceQueueFill.Describe(ch)
		e.bucketEpTapRebalanceQueueItemondisk.Describe(ch)
		e.bucketEpTapRebalanceTotalBacklogSize.Describe(ch)
		e.bucketEpTapReplicaCount.Describe(ch)
		e.bucketEpTapReplicaQlen.Describe(ch)
		e.bucketEpTapReplicaQueueBackfillremaining.Describe(ch)
		e.bucketEpTapReplicaQueueBackoff.Describe(ch)
		e.bucketEpTapReplicaQueueDrain.Describe(ch)
		e.bucketEpTapReplicaQueueFill.Describe(ch)
		e.bucketEpTapReplicaQueueItemondisk.Describe(ch)
		e.bucketEpTapReplicaTotalBacklogSize.Describe(ch)
		e.bucketEpTapTotalCount.Describe(ch)
		e.bucketEpTapTotalQlen.Describe(ch)
		e.bucketEpTapTotalQueueBackfillremaining.Describe(ch)
		e.bucketEpTapTotalQueueBackoff.Describe(ch)
		e.bucketEpTapTotalQueueDrain.Describe(ch)
		e.bucketEpTapTotalQueueFill.Describe(ch)
		e.bucketEpTapTotalQueueItemondisk.Describe(ch)
		e.bucketEpTapTotalTotalBacklogSize.Describe(ch)
		e.bucketEpTapUserCount.Describe(ch)
		e.bucketEpTapUserQlen.Describe(ch)
		e.bucketEpTapUserQueueBackfillremaining.Describe(ch)
		e.bucketEpTapUserQueueBackoff.Describe(ch)
		e.bucketEpTapUserQueueDrain.Describe(ch)
		e.bucketEpTapUserQueueFill.Describe(ch)
		e.bucketEpTapUserQueueItemondisk.Describe(ch)
		e.bucketEpTapUserTotalBacklogSize.Describe(ch)
	}
	if e.context.CouchbaseVersion == "5.1.1" {
		e.bucketAvgActiveTimestampDrift.Describe(ch)
		e.bucketAvgReplicaTimestampDrift.Describe(ch)
		e.bucketEpActiveAheadExceptions.Describe(ch)
		e.bucketEpActiveHlcDrift.Describe(ch)
		e.bucketEpActiveHlcDriftCount.Describe(ch)
		e.bucketEpClockCasDriftThresholdExceeded.Describe(ch)
		e.bucketEpReplicaAheadExceptions.Describe(ch)
		e.bucketEpReplicaHlcDrift.Describe(ch)
		e.bucketEpReplicaHlcDriftCount.Describe(ch)
	}
}

// Collect fetches data for each exported metric.
func (e *BucketStatsExporter) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()
	e.scrape()
	e.bucketCouchTotalDiskSize.Collect(ch)
	e.bucketCouchDocsFragmentation.Collect(ch)
	e.bucketCouchViewsFragmentation.Collect(ch)
	e.bucketHitRatio.Collect(ch)
	e.bucketEpCacheMissRate.Collect(ch)
	e.bucketEpResidentItemsRate.Collect(ch)
	e.bucketVbAvgActiveQueueAge.Collect(ch)
	e.bucketVbAvgReplicaQueueAge.Collect(ch)
	e.bucketVbAvgPendingQueueAge.Collect(ch)
	e.bucketVbAvgTotalQueueAge.Collect(ch)
	e.bucketVbActiveResidentItemsRatio.Collect(ch)
	e.bucketVbReplicaResidentItemsRatio.Collect(ch)
	e.bucketVbPendingResidentItemsRatio.Collect(ch)
	e.bucketAvgDiskUpdateTime.Collect(ch)
	e.bucketAvgDiskCommitTime.Collect(ch)
	e.bucketAvgBgWaitTime.Collect(ch)
	e.bucketEpDcpViewsIndexesCount.Collect(ch)
	e.bucketEpDcpViewsIndexesItemsRemaining.Collect(ch)
	e.bucketEpDcpViewsIndexesProducerCount.Collect(ch)
	e.bucketEpDcpViewsIndexesTotalBacklogSize.Collect(ch)
	e.bucketEpDcpViewsIndexesItemsSent.Collect(ch)
	e.bucketEpDcpViewsIndexesTotalBytes.Collect(ch)
	e.bucketEpDcpViewsIndexesBackoff.Collect(ch)
	e.bucketBgWaitCount.Collect(ch)
	e.bucketBgWaitTotal.Collect(ch)
	e.bucketBytesRead.Collect(ch)
	e.bucketBytesWritten.Collect(ch)
	e.bucketCasBadval.Collect(ch)
	e.bucketCasHits.Collect(ch)
	e.bucketCasMisses.Collect(ch)
	e.bucketCmdGet.Collect(ch)
	e.bucketCmdSet.Collect(ch)
	e.bucketCouchDocsActualDiskSize.Collect(ch)
	e.bucketCouchDocsDataSize.Collect(ch)
	e.bucketCouchDocsDiskSize.Collect(ch)
	e.bucketCouchSpatialDataSize.Collect(ch)
	e.bucketCouchSpatialDiskSize.Collect(ch)
	e.bucketCouchSpatialOps.Collect(ch)
	e.bucketCouchViewsActualDiskSize.Collect(ch)
	e.bucketCouchViewsDataSize.Collect(ch)
	e.bucketCouchViewsDiskSize.Collect(ch)
	e.bucketCouchViewsOps.Collect(ch)
	e.bucketCurrConnections.Collect(ch)
	e.bucketCurrItems.Collect(ch)
	e.bucketCurrItemsTot.Collect(ch)
	e.bucketDecrHits.Collect(ch)
	e.bucketDecrMisses.Collect(ch)
	e.bucketDeleteHits.Collect(ch)
	e.bucketDeleteMisses.Collect(ch)
	e.bucketDiskCommitCount.Collect(ch)
	e.bucketDiskCommitTotal.Collect(ch)
	e.bucketDiskUpdateCount.Collect(ch)
	e.bucketDiskUpdateTotal.Collect(ch)
	e.bucketDiskWriteQueue.Collect(ch)
	e.bucketEpBgFetched.Collect(ch)
	e.bucketEpDcp2IBackoff.Collect(ch)
	e.bucketEpDcp2ICount.Collect(ch)
	e.bucketEpDcp2IItemsRemaining.Collect(ch)
	e.bucketEpDcp2IItemsSent.Collect(ch)
	e.bucketEpDcp2IProducerCount.Collect(ch)
	e.bucketEpDcp2ITotalBacklogSize.Collect(ch)
	e.bucketEpDcp2ITotalBytes.Collect(ch)
	e.bucketEpDcpFtsBackoff.Collect(ch)
	e.bucketEpDcpFtsCount.Collect(ch)
	e.bucketEpDcpFtsItemsRemaining.Collect(ch)
	e.bucketEpDcpFtsItemsSent.Collect(ch)
	e.bucketEpDcpFtsProducerCount.Collect(ch)
	e.bucketEpDcpFtsTotalBacklogSize.Collect(ch)
	e.bucketEpDcpFtsTotalBytes.Collect(ch)
	e.bucketEpDcpOtherBackoff.Collect(ch)
	e.bucketEpDcpOtherCount.Collect(ch)
	e.bucketEpDcpOtherItemsRemaining.Collect(ch)
	e.bucketEpDcpOtherItemsSent.Collect(ch)
	e.bucketEpDcpOtherProducerCount.Collect(ch)
	e.bucketEpDcpOtherTotalBacklogSize.Collect(ch)
	e.bucketEpDcpOtherTotalBytes.Collect(ch)
	e.bucketEpDcpReplicaBackoff.Collect(ch)
	e.bucketEpDcpReplicaCount.Collect(ch)
	e.bucketEpDcpReplicaItemsRemaining.Collect(ch)
	e.bucketEpDcpReplicaItemsSent.Collect(ch)
	e.bucketEpDcpReplicaProducerCount.Collect(ch)
	e.bucketEpDcpReplicaTotalBacklogSize.Collect(ch)
	e.bucketEpDcpReplicaTotalBytes.Collect(ch)
	e.bucketEpDcpViewsBackoff.Collect(ch)
	e.bucketEpDcpViewsCount.Collect(ch)
	e.bucketEpDcpViewsItemsRemaining.Collect(ch)
	e.bucketEpDcpViewsItemsSent.Collect(ch)
	e.bucketEpDcpViewsProducerCount.Collect(ch)
	e.bucketEpDcpViewsTotalBacklogSize.Collect(ch)
	e.bucketEpDcpViewsTotalBytes.Collect(ch)
	e.bucketEpDcpXdcrBackoff.Collect(ch)
	e.bucketEpDcpXdcrCount.Collect(ch)
	e.bucketEpDcpXdcrItemsRemaining.Collect(ch)
	e.bucketEpDcpXdcrItemsSent.Collect(ch)
	e.bucketEpDcpXdcrProducerCount.Collect(ch)
	e.bucketEpDcpXdcrTotalBacklogSize.Collect(ch)
	e.bucketEpDcpXdcrTotalBytes.Collect(ch)
	e.bucketEpDiskqueueDrain.Collect(ch)
	e.bucketEpDiskqueueFill.Collect(ch)
	e.bucketEpDiskqueueItems.Collect(ch)
	e.bucketEpFlusherTodo.Collect(ch)
	e.bucketEpItemCommitFailed.Collect(ch)
	e.bucketEpKvSize.Collect(ch)
	e.bucketEpMaxSize.Collect(ch)
	e.bucketEpMemHighWat.Collect(ch)
	e.bucketEpMemLowWat.Collect(ch)
	e.bucketEpMetaDataMemory.Collect(ch)
	e.bucketEpNumNonResident.Collect(ch)
	e.bucketEpNumOpsDelMeta.Collect(ch)
	e.bucketEpNumOpsDelRetMeta.Collect(ch)
	e.bucketEpNumOpsGetMeta.Collect(ch)
	e.bucketEpNumOpsSetMeta.Collect(ch)
	e.bucketEpNumOpsSetRetMeta.Collect(ch)
	e.bucketEpNumValueEjects.Collect(ch)
	e.bucketEpOomErrors.Collect(ch)
	e.bucketEpOpsCreate.Collect(ch)
	e.bucketEpOpsUpdate.Collect(ch)
	e.bucketEpOverhead.Collect(ch)
	e.bucketEpQueueSize.Collect(ch)
	e.bucketEpTmpOomErrors.Collect(ch)
	e.bucketEpVbTotal.Collect(ch)
	e.bucketEvictions.Collect(ch)
	e.bucketGetHits.Collect(ch)
	e.bucketGetMisses.Collect(ch)
	e.bucketIncrHits.Collect(ch)
	e.bucketIncrMisses.Collect(ch)
	e.bucketMemUsed.Collect(ch)
	e.bucketMisses.Collect(ch)
	e.bucketOps.Collect(ch)
	e.bucketVbActiveEject.Collect(ch)
	e.bucketVbActiveItmMemory.Collect(ch)
	e.bucketVbActiveMetaDataMemory.Collect(ch)
	e.bucketVbActiveNum.Collect(ch)
	e.bucketVbActiveNumNonResident.Collect(ch)
	e.bucketVbActiveOpsCreate.Collect(ch)
	e.bucketVbActiveOpsUpdate.Collect(ch)
	e.bucketVbActiveQueueAge.Collect(ch)
	e.bucketVbActiveQueueDrain.Collect(ch)
	e.bucketVbActiveQueueFill.Collect(ch)
	e.bucketVbActiveQueueSize.Collect(ch)
	e.bucketVbPendingCurrItems.Collect(ch)
	e.bucketVbPendingEject.Collect(ch)
	e.bucketVbPendingItmMemory.Collect(ch)
	e.bucketVbPendingMetaDataMemory.Collect(ch)
	e.bucketVbPendingNum.Collect(ch)
	e.bucketVbPendingNumNonResident.Collect(ch)
	e.bucketVbPendingOpsCreate.Collect(ch)
	e.bucketVbPendingOpsUpdate.Collect(ch)
	e.bucketVbPendingQueueAge.Collect(ch)
	e.bucketVbPendingQueueDrain.Collect(ch)
	e.bucketVbPendingQueueFill.Collect(ch)
	e.bucketVbPendingQueueSize.Collect(ch)
	e.bucketVbReplicaCurrItems.Collect(ch)
	e.bucketVbReplicaEject.Collect(ch)
	e.bucketVbReplicaItmMemory.Collect(ch)
	e.bucketVbReplicaMetaDataMemory.Collect(ch)
	e.bucketVbReplicaNum.Collect(ch)
	e.bucketVbReplicaNumNonResident.Collect(ch)
	e.bucketVbReplicaOpsCreate.Collect(ch)
	e.bucketVbReplicaOpsUpdate.Collect(ch)
	e.bucketVbReplicaQueueAge.Collect(ch)
	e.bucketVbReplicaQueueDrain.Collect(ch)
	e.bucketVbReplicaQueueFill.Collect(ch)
	e.bucketVbReplicaQueueSize.Collect(ch)
	e.bucketVbTotalQueueAge.Collect(ch)
	e.bucketXdcOps.Collect(ch)
	e.bucketCPUIdleMs.Collect(ch)
	e.bucketCPULocalMs.Collect(ch)
	e.bucketCPUUtilizationRate.Collect(ch)
	e.bucketHibernatedRequests.Collect(ch)
	e.bucketHibernatedWaked.Collect(ch)
	e.bucketMemActualFree.Collect(ch)
	e.bucketMemActualUsed.Collect(ch)
	e.bucketMemFree.Collect(ch)
	e.bucketMemTotal.Collect(ch)
	e.bucketMemUsedSys.Collect(ch)
	e.bucketRestRequests.Collect(ch)
	e.bucketSwapTotal.Collect(ch)
	e.bucketSwapUsed.Collect(ch)
	if e.context.CouchbaseVersion == "4.5.1" {
		e.bucketEpTapRebalanceCount.Collect(ch)
		e.bucketEpTapRebalanceQlen.Collect(ch)
		e.bucketEpTapRebalanceQueueBackfillremaining.Collect(ch)
		e.bucketEpTapRebalanceQueueBackoff.Collect(ch)
		e.bucketEpTapRebalanceQueueDrain.Collect(ch)
		e.bucketEpTapRebalanceQueueFill.Collect(ch)
		e.bucketEpTapRebalanceQueueItemondisk.Collect(ch)
		e.bucketEpTapRebalanceTotalBacklogSize.Collect(ch)
		e.bucketEpTapReplicaCount.Collect(ch)
		e.bucketEpTapReplicaQlen.Collect(ch)
		e.bucketEpTapReplicaQueueBackfillremaining.Collect(ch)
		e.bucketEpTapReplicaQueueBackoff.Collect(ch)
		e.bucketEpTapReplicaQueueDrain.Collect(ch)
		e.bucketEpTapReplicaQueueFill.Collect(ch)
		e.bucketEpTapReplicaQueueItemondisk.Collect(ch)
		e.bucketEpTapReplicaTotalBacklogSize.Collect(ch)
		e.bucketEpTapTotalCount.Collect(ch)
		e.bucketEpTapTotalQlen.Collect(ch)
		e.bucketEpTapTotalQueueBackfillremaining.Collect(ch)
		e.bucketEpTapTotalQueueBackoff.Collect(ch)
		e.bucketEpTapTotalQueueDrain.Collect(ch)
		e.bucketEpTapTotalQueueFill.Collect(ch)
		e.bucketEpTapTotalQueueItemondisk.Collect(ch)
		e.bucketEpTapTotalTotalBacklogSize.Collect(ch)
		e.bucketEpTapUserCount.Collect(ch)
		e.bucketEpTapUserQlen.Collect(ch)
		e.bucketEpTapUserQueueBackfillremaining.Collect(ch)
		e.bucketEpTapUserQueueBackoff.Collect(ch)
		e.bucketEpTapUserQueueDrain.Collect(ch)
		e.bucketEpTapUserQueueFill.Collect(ch)
		e.bucketEpTapUserQueueItemondisk.Collect(ch)
		e.bucketEpTapUserTotalBacklogSize.Collect(ch)
	}
	if e.context.CouchbaseVersion == "5.1.1" {
		e.bucketAvgActiveTimestampDrift.Collect(ch)
		e.bucketAvgReplicaTimestampDrift.Collect(ch)
		e.bucketEpActiveAheadExceptions.Collect(ch)
		e.bucketEpActiveHlcDrift.Collect(ch)
		e.bucketEpActiveHlcDriftCount.Collect(ch)
		e.bucketEpClockCasDriftThresholdExceeded.Collect(ch)
		e.bucketEpReplicaAheadExceptions.Collect(ch)
		e.bucketEpReplicaHlcDrift.Collect(ch)
		e.bucketEpReplicaHlcDriftCount.Collect(ch)
	}
	mutex.Unlock()
}

func (e *BucketStatsExporter) scrape() {
	req, err := http.NewRequest("GET", e.context.URI+bucketRoute, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.context.Username, e.context.Password)
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Error(req.URL.Path + ": " + res.Status)
		return
	}

	var buckets []BucketName
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &buckets)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.context.URI + bucketRoute + " - data: " + string(body))

	for _, bucket := range buckets {

		req, err := http.NewRequest("GET", e.context.URI+bucketRoute+"/"+bucket.Name+"/stats", nil)
		if err != nil {
			log.Error(err.Error())
			return
		}
		req.SetBasicAuth(e.context.Username, e.context.Password)
		client := http.Client{Timeout: 10 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if res.StatusCode != 200 {
			log.Error(req.URL.Path + ": " + res.Status)
			return
		}

		var bucketStats BucketStatsData
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
		err = json.Unmarshal([]byte(body), &bucketStats)
		if err != nil {
			log.Error(err.Error())
		}

		log.Debug("GET " + e.context.URI + bucketRoute + "/" + bucket.Name + "/stats" + " - data: " + string(body))

		e.bucketCouchTotalDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchTotalDiskSize[len(bucketStats.Op.Samples.CouchTotalDiskSize)-1]))
		e.bucketCouchDocsFragmentation.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchDocsFragmentation[len(bucketStats.Op.Samples.CouchDocsFragmentation)-1]))
		e.bucketCouchViewsFragmentation.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchViewsFragmentation[len(bucketStats.Op.Samples.CouchViewsFragmentation)-1]))
		e.bucketHitRatio.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.HitRatio[len(bucketStats.Op.Samples.HitRatio)-1]))
		e.bucketEpCacheMissRate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpCacheMissRate[len(bucketStats.Op.Samples.EpCacheMissRate)-1]))
		e.bucketEpResidentItemsRate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpResidentItemsRate[len(bucketStats.Op.Samples.EpResidentItemsRate)-1]))
		e.bucketVbAvgActiveQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbAvgActiveQueueAge[len(bucketStats.Op.Samples.VbAvgActiveQueueAge)-1]))
		e.bucketVbAvgReplicaQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbAvgReplicaQueueAge[len(bucketStats.Op.Samples.VbAvgReplicaQueueAge)-1]))
		e.bucketVbAvgPendingQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbAvgPendingQueueAge[len(bucketStats.Op.Samples.VbAvgPendingQueueAge)-1]))
		e.bucketVbAvgTotalQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbAvgTotalQueueAge[len(bucketStats.Op.Samples.VbAvgTotalQueueAge)-1]))
		e.bucketVbActiveResidentItemsRatio.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveResidentItemsRatio[len(bucketStats.Op.Samples.VbActiveResidentItemsRatio)-1]))
		e.bucketVbReplicaResidentItemsRatio.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaResidentItemsRatio[len(bucketStats.Op.Samples.VbReplicaResidentItemsRatio)-1]))
		e.bucketVbPendingResidentItemsRatio.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingResidentItemsRatio[len(bucketStats.Op.Samples.VbPendingResidentItemsRatio)-1]))
		e.bucketAvgDiskUpdateTime.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.AvgDiskUpdateTime[len(bucketStats.Op.Samples.AvgDiskUpdateTime)-1]))
		e.bucketAvgDiskCommitTime.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.AvgDiskCommitTime[len(bucketStats.Op.Samples.AvgDiskCommitTime)-1]))
		e.bucketAvgBgWaitTime.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.AvgBgWaitTime[len(bucketStats.Op.Samples.AvgBgWaitTime)-1]))
		e.bucketEpDcpViewsIndexesCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesCount[len(bucketStats.Op.Samples.EpDcpViewsIndexesCount)-1]))
		e.bucketEpDcpViewsIndexesItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesItemsRemaining[len(bucketStats.Op.Samples.EpDcpViewsIndexesItemsRemaining)-1]))
		e.bucketEpDcpViewsIndexesProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesProducerCount[len(bucketStats.Op.Samples.EpDcpViewsIndexesProducerCount)-1]))
		e.bucketEpDcpViewsIndexesTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpViewsIndexesTotalBacklogSize)-1]))
		e.bucketEpDcpViewsIndexesItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesItemsSent[len(bucketStats.Op.Samples.EpDcpViewsIndexesItemsSent)-1]))
		e.bucketEpDcpViewsIndexesTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesTotalBytes[len(bucketStats.Op.Samples.EpDcpViewsIndexesTotalBytes)-1]))
		e.bucketEpDcpViewsIndexesBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsIndexesBackoff[len(bucketStats.Op.Samples.EpDcpViewsIndexesBackoff)-1]))
		e.bucketBgWaitCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.BgWaitCount[len(bucketStats.Op.Samples.BgWaitCount)-1]))
		e.bucketBgWaitTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.BgWaitTotal[len(bucketStats.Op.Samples.BgWaitTotal)-1]))
		e.bucketBytesRead.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.BytesRead[len(bucketStats.Op.Samples.BytesRead)-1]))
		e.bucketBytesWritten.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.BytesWritten[len(bucketStats.Op.Samples.BytesWritten)-1]))
		e.bucketCasBadval.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CasBadval[len(bucketStats.Op.Samples.CasBadval)-1]))
		e.bucketCasHits.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CasHits[len(bucketStats.Op.Samples.CasHits)-1]))
		e.bucketCasMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CasMisses[len(bucketStats.Op.Samples.CasMisses)-1]))
		e.bucketCmdGet.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CmdGet[len(bucketStats.Op.Samples.CmdGet)-1]))
		e.bucketCmdSet.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CmdSet[len(bucketStats.Op.Samples.CmdSet)-1]))
		e.bucketCouchDocsActualDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchDocsActualDiskSize[len(bucketStats.Op.Samples.CouchDocsActualDiskSize)-1]))
		e.bucketCouchDocsDataSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchDocsDataSize[len(bucketStats.Op.Samples.CouchDocsDataSize)-1]))
		e.bucketCouchDocsDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchDocsDiskSize[len(bucketStats.Op.Samples.CouchDocsDiskSize)-1]))
		e.bucketCouchSpatialDataSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchSpatialDataSize[len(bucketStats.Op.Samples.CouchSpatialDataSize)-1]))
		e.bucketCouchSpatialDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchSpatialDiskSize[len(bucketStats.Op.Samples.CouchSpatialDiskSize)-1]))
		e.bucketCouchSpatialOps.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchSpatialOps[len(bucketStats.Op.Samples.CouchSpatialOps)-1]))
		e.bucketCouchViewsActualDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchViewsActualDiskSize[len(bucketStats.Op.Samples.CouchViewsActualDiskSize)-1]))
		e.bucketCouchViewsDataSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchViewsDataSize[len(bucketStats.Op.Samples.CouchViewsDataSize)-1]))
		e.bucketCouchViewsDiskSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchViewsDiskSize[len(bucketStats.Op.Samples.CouchViewsDiskSize)-1]))
		e.bucketCouchViewsOps.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CouchViewsOps[len(bucketStats.Op.Samples.CouchViewsOps)-1]))
		e.bucketCurrConnections.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CurrConnections[len(bucketStats.Op.Samples.CurrConnections)-1]))
		e.bucketCurrItems.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CurrItems[len(bucketStats.Op.Samples.CurrItems)-1]))
		e.bucketCurrItemsTot.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CurrItemsTot[len(bucketStats.Op.Samples.CurrItemsTot)-1]))
		e.bucketDecrHits.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DecrHits[len(bucketStats.Op.Samples.DecrHits)-1]))
		e.bucketDecrMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DecrMisses[len(bucketStats.Op.Samples.DecrMisses)-1]))
		e.bucketDeleteHits.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DeleteHits[len(bucketStats.Op.Samples.DeleteHits)-1]))
		e.bucketDeleteMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DeleteMisses[len(bucketStats.Op.Samples.DeleteMisses)-1]))
		e.bucketDiskCommitCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DiskCommitCount[len(bucketStats.Op.Samples.DiskCommitCount)-1]))
		e.bucketDiskCommitTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DiskCommitTotal[len(bucketStats.Op.Samples.DiskCommitTotal)-1]))
		e.bucketDiskUpdateCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DiskUpdateCount[len(bucketStats.Op.Samples.DiskUpdateCount)-1]))
		e.bucketDiskUpdateTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DiskUpdateTotal[len(bucketStats.Op.Samples.DiskUpdateTotal)-1]))
		e.bucketDiskWriteQueue.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.DiskWriteQueue[len(bucketStats.Op.Samples.DiskWriteQueue)-1]))
		e.bucketEpBgFetched.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpBgFetched[len(bucketStats.Op.Samples.EpBgFetched)-1]))
		e.bucketEpDcp2IBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2IBackoff[len(bucketStats.Op.Samples.EpDcp2IBackoff)-1]))
		e.bucketEpDcp2ICount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2ICount[len(bucketStats.Op.Samples.EpDcp2ICount)-1]))
		e.bucketEpDcp2IItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2IItemsRemaining[len(bucketStats.Op.Samples.EpDcp2IItemsRemaining)-1]))
		e.bucketEpDcp2IItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2IItemsSent[len(bucketStats.Op.Samples.EpDcp2IItemsSent)-1]))
		e.bucketEpDcp2IProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2IProducerCount[len(bucketStats.Op.Samples.EpDcp2IProducerCount)-1]))
		e.bucketEpDcp2ITotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2ITotalBacklogSize[len(bucketStats.Op.Samples.EpDcp2ITotalBacklogSize)-1]))
		e.bucketEpDcp2ITotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcp2ITotalBytes[len(bucketStats.Op.Samples.EpDcp2ITotalBytes)-1]))
		e.bucketEpDcpFtsBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsBackoff[len(bucketStats.Op.Samples.EpDcpFtsBackoff)-1]))
		e.bucketEpDcpFtsCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsCount[len(bucketStats.Op.Samples.EpDcpFtsCount)-1]))
		e.bucketEpDcpFtsItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsItemsRemaining[len(bucketStats.Op.Samples.EpDcpFtsItemsRemaining)-1]))
		e.bucketEpDcpFtsItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsItemsSent[len(bucketStats.Op.Samples.EpDcpFtsItemsSent)-1]))
		e.bucketEpDcpFtsProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsProducerCount[len(bucketStats.Op.Samples.EpDcpFtsProducerCount)-1]))
		e.bucketEpDcpFtsTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpFtsTotalBacklogSize)-1]))
		e.bucketEpDcpFtsTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpFtsTotalBytes[len(bucketStats.Op.Samples.EpDcpFtsTotalBytes)-1]))
		e.bucketEpDcpOtherBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherBackoff[len(bucketStats.Op.Samples.EpDcpOtherBackoff)-1]))
		e.bucketEpDcpOtherCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherCount[len(bucketStats.Op.Samples.EpDcpOtherCount)-1]))
		e.bucketEpDcpOtherItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherItemsRemaining[len(bucketStats.Op.Samples.EpDcpOtherItemsRemaining)-1]))
		e.bucketEpDcpOtherItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherItemsSent[len(bucketStats.Op.Samples.EpDcpOtherItemsSent)-1]))
		e.bucketEpDcpOtherProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherProducerCount[len(bucketStats.Op.Samples.EpDcpOtherProducerCount)-1]))
		e.bucketEpDcpOtherTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpOtherTotalBacklogSize)-1]))
		e.bucketEpDcpOtherTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpOtherTotalBytes[len(bucketStats.Op.Samples.EpDcpOtherTotalBytes)-1]))
		e.bucketEpDcpReplicaBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaBackoff[len(bucketStats.Op.Samples.EpDcpReplicaBackoff)-1]))
		e.bucketEpDcpReplicaCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaCount[len(bucketStats.Op.Samples.EpDcpReplicaCount)-1]))
		e.bucketEpDcpReplicaItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaItemsRemaining[len(bucketStats.Op.Samples.EpDcpReplicaItemsRemaining)-1]))
		e.bucketEpDcpReplicaItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaItemsSent[len(bucketStats.Op.Samples.EpDcpReplicaItemsSent)-1]))
		e.bucketEpDcpReplicaProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaProducerCount[len(bucketStats.Op.Samples.EpDcpReplicaProducerCount)-1]))
		e.bucketEpDcpReplicaTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpReplicaTotalBacklogSize)-1]))
		e.bucketEpDcpReplicaTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpReplicaTotalBytes[len(bucketStats.Op.Samples.EpDcpReplicaTotalBytes)-1]))
		e.bucketEpDcpViewsBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsBackoff[len(bucketStats.Op.Samples.EpDcpViewsBackoff)-1]))
		e.bucketEpDcpViewsCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsCount[len(bucketStats.Op.Samples.EpDcpViewsCount)-1]))
		e.bucketEpDcpViewsItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsItemsRemaining[len(bucketStats.Op.Samples.EpDcpViewsItemsRemaining)-1]))
		e.bucketEpDcpViewsItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsItemsSent[len(bucketStats.Op.Samples.EpDcpViewsItemsSent)-1]))
		e.bucketEpDcpViewsProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsProducerCount[len(bucketStats.Op.Samples.EpDcpViewsProducerCount)-1]))
		e.bucketEpDcpViewsTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpViewsTotalBacklogSize)-1]))
		e.bucketEpDcpViewsTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpViewsTotalBytes[len(bucketStats.Op.Samples.EpDcpViewsTotalBytes)-1]))
		e.bucketEpDcpXdcrBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrBackoff[len(bucketStats.Op.Samples.EpDcpXdcrBackoff)-1]))
		e.bucketEpDcpXdcrCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrCount[len(bucketStats.Op.Samples.EpDcpXdcrCount)-1]))
		e.bucketEpDcpXdcrItemsRemaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrItemsRemaining[len(bucketStats.Op.Samples.EpDcpXdcrItemsRemaining)-1]))
		e.bucketEpDcpXdcrItemsSent.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrItemsSent[len(bucketStats.Op.Samples.EpDcpXdcrItemsSent)-1]))
		e.bucketEpDcpXdcrProducerCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrProducerCount[len(bucketStats.Op.Samples.EpDcpXdcrProducerCount)-1]))
		e.bucketEpDcpXdcrTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrTotalBacklogSize[len(bucketStats.Op.Samples.EpDcpXdcrTotalBacklogSize)-1]))
		e.bucketEpDcpXdcrTotalBytes.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDcpXdcrTotalBytes[len(bucketStats.Op.Samples.EpDcpXdcrTotalBytes)-1]))
		e.bucketEpDiskqueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDiskqueueDrain[len(bucketStats.Op.Samples.EpDiskqueueDrain)-1]))
		e.bucketEpDiskqueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDiskqueueFill[len(bucketStats.Op.Samples.EpDiskqueueFill)-1]))
		e.bucketEpDiskqueueItems.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpDiskqueueItems[len(bucketStats.Op.Samples.EpDiskqueueItems)-1]))
		e.bucketEpFlusherTodo.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpFlusherTodo[len(bucketStats.Op.Samples.EpFlusherTodo)-1]))
		e.bucketEpItemCommitFailed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpItemCommitFailed[len(bucketStats.Op.Samples.EpItemCommitFailed)-1]))
		e.bucketEpKvSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpKvSize[len(bucketStats.Op.Samples.EpKvSize)-1]))
		e.bucketEpMaxSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpMaxSize[len(bucketStats.Op.Samples.EpMaxSize)-1]))
		e.bucketEpMemHighWat.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpMemHighWat[len(bucketStats.Op.Samples.EpMemHighWat)-1]))
		e.bucketEpMemLowWat.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpMemLowWat[len(bucketStats.Op.Samples.EpMemLowWat)-1]))
		e.bucketEpMetaDataMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpMetaDataMemory[len(bucketStats.Op.Samples.EpMetaDataMemory)-1]))
		e.bucketEpNumNonResident.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumNonResident[len(bucketStats.Op.Samples.EpNumNonResident)-1]))
		e.bucketEpNumOpsDelMeta.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumOpsDelMeta[len(bucketStats.Op.Samples.EpNumOpsDelMeta)-1]))
		e.bucketEpNumOpsDelRetMeta.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumOpsDelRetMeta[len(bucketStats.Op.Samples.EpNumOpsDelRetMeta)-1]))
		e.bucketEpNumOpsGetMeta.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumOpsGetMeta[len(bucketStats.Op.Samples.EpNumOpsGetMeta)-1]))
		e.bucketEpNumOpsSetMeta.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumOpsSetMeta[len(bucketStats.Op.Samples.EpNumOpsSetMeta)-1]))
		e.bucketEpNumOpsSetRetMeta.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumOpsSetRetMeta[len(bucketStats.Op.Samples.EpNumOpsSetRetMeta)-1]))
		e.bucketEpNumValueEjects.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpNumValueEjects[len(bucketStats.Op.Samples.EpNumValueEjects)-1]))
		e.bucketEpOomErrors.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpOomErrors[len(bucketStats.Op.Samples.EpOomErrors)-1]))
		e.bucketEpOpsCreate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpOpsCreate[len(bucketStats.Op.Samples.EpOpsCreate)-1]))
		e.bucketEpOpsUpdate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpOpsUpdate[len(bucketStats.Op.Samples.EpOpsUpdate)-1]))
		e.bucketEpOverhead.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpOverhead[len(bucketStats.Op.Samples.EpOverhead)-1]))
		e.bucketEpQueueSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpQueueSize[len(bucketStats.Op.Samples.EpQueueSize)-1]))
		e.bucketEpTmpOomErrors.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTmpOomErrors[len(bucketStats.Op.Samples.EpTmpOomErrors)-1]))
		e.bucketEpVbTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpVbTotal[len(bucketStats.Op.Samples.EpVbTotal)-1]))
		e.bucketEvictions.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.Evictions[len(bucketStats.Op.Samples.Evictions)-1]))
		e.bucketGetHits.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.GetHits[len(bucketStats.Op.Samples.GetHits)-1]))
		e.bucketGetMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.GetMisses[len(bucketStats.Op.Samples.GetMisses)-1]))
		e.bucketIncrHits.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.IncrHits[len(bucketStats.Op.Samples.IncrHits)-1]))
		e.bucketIncrMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.IncrMisses[len(bucketStats.Op.Samples.IncrMisses)-1]))
		e.bucketMemUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemUsed[len(bucketStats.Op.Samples.MemUsed)-1]))
		e.bucketMisses.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.Misses[len(bucketStats.Op.Samples.Misses)-1]))
		e.bucketOps.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.Ops[len(bucketStats.Op.Samples.Ops)-1]))
		e.bucketVbActiveEject.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveEject[len(bucketStats.Op.Samples.VbActiveEject)-1]))
		e.bucketVbActiveItmMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveItmMemory[len(bucketStats.Op.Samples.VbActiveItmMemory)-1]))
		e.bucketVbActiveMetaDataMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveMetaDataMemory[len(bucketStats.Op.Samples.VbActiveMetaDataMemory)-1]))
		e.bucketVbActiveNum.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveNum[len(bucketStats.Op.Samples.VbActiveNum)-1]))
		e.bucketVbActiveNumNonResident.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveNumNonResident[len(bucketStats.Op.Samples.VbActiveNumNonResident)-1]))
		e.bucketVbActiveOpsCreate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveOpsCreate[len(bucketStats.Op.Samples.VbActiveOpsCreate)-1]))
		e.bucketVbActiveOpsUpdate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveOpsUpdate[len(bucketStats.Op.Samples.VbActiveOpsUpdate)-1]))
		e.bucketVbActiveQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveQueueAge[len(bucketStats.Op.Samples.VbActiveQueueAge)-1]))
		e.bucketVbActiveQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveQueueDrain[len(bucketStats.Op.Samples.VbActiveQueueDrain)-1]))
		e.bucketVbActiveQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveQueueFill[len(bucketStats.Op.Samples.VbActiveQueueFill)-1]))
		e.bucketVbActiveQueueSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbActiveQueueSize[len(bucketStats.Op.Samples.VbActiveQueueSize)-1]))
		e.bucketVbPendingCurrItems.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingCurrItems[len(bucketStats.Op.Samples.VbPendingCurrItems)-1]))
		e.bucketVbPendingEject.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingEject[len(bucketStats.Op.Samples.VbPendingEject)-1]))
		e.bucketVbPendingItmMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingItmMemory[len(bucketStats.Op.Samples.VbPendingItmMemory)-1]))
		e.bucketVbPendingMetaDataMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingMetaDataMemory[len(bucketStats.Op.Samples.VbPendingMetaDataMemory)-1]))
		e.bucketVbPendingNum.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingNum[len(bucketStats.Op.Samples.VbPendingNum)-1]))
		e.bucketVbPendingNumNonResident.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingNumNonResident[len(bucketStats.Op.Samples.VbPendingNumNonResident)-1]))
		e.bucketVbPendingOpsCreate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingOpsCreate[len(bucketStats.Op.Samples.VbPendingOpsCreate)-1]))
		e.bucketVbPendingOpsUpdate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingOpsUpdate[len(bucketStats.Op.Samples.VbPendingOpsUpdate)-1]))
		e.bucketVbPendingQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingQueueAge[len(bucketStats.Op.Samples.VbPendingQueueAge)-1]))
		e.bucketVbPendingQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingQueueDrain[len(bucketStats.Op.Samples.VbPendingQueueDrain)-1]))
		e.bucketVbPendingQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingQueueFill[len(bucketStats.Op.Samples.VbPendingQueueFill)-1]))
		e.bucketVbPendingQueueSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbPendingQueueSize[len(bucketStats.Op.Samples.VbPendingQueueSize)-1]))
		e.bucketVbReplicaCurrItems.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaCurrItems[len(bucketStats.Op.Samples.VbReplicaCurrItems)-1]))
		e.bucketVbReplicaEject.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaEject[len(bucketStats.Op.Samples.VbReplicaEject)-1]))
		e.bucketVbReplicaItmMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaItmMemory[len(bucketStats.Op.Samples.VbReplicaItmMemory)-1]))
		e.bucketVbReplicaMetaDataMemory.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaMetaDataMemory[len(bucketStats.Op.Samples.VbReplicaMetaDataMemory)-1]))
		e.bucketVbReplicaNum.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaNum[len(bucketStats.Op.Samples.VbReplicaNum)-1]))
		e.bucketVbReplicaNumNonResident.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaNumNonResident[len(bucketStats.Op.Samples.VbReplicaNumNonResident)-1]))
		e.bucketVbReplicaOpsCreate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaOpsCreate[len(bucketStats.Op.Samples.VbReplicaOpsCreate)-1]))
		e.bucketVbReplicaOpsUpdate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaOpsUpdate[len(bucketStats.Op.Samples.VbReplicaOpsUpdate)-1]))
		e.bucketVbReplicaQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaQueueAge[len(bucketStats.Op.Samples.VbReplicaQueueAge)-1]))
		e.bucketVbReplicaQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaQueueDrain[len(bucketStats.Op.Samples.VbReplicaQueueDrain)-1]))
		e.bucketVbReplicaQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaQueueFill[len(bucketStats.Op.Samples.VbReplicaQueueFill)-1]))
		e.bucketVbReplicaQueueSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbReplicaQueueSize[len(bucketStats.Op.Samples.VbReplicaQueueSize)-1]))
		e.bucketVbTotalQueueAge.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.VbTotalQueueAge[len(bucketStats.Op.Samples.VbTotalQueueAge)-1]))
		e.bucketXdcOps.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.XdcOps[len(bucketStats.Op.Samples.XdcOps)-1]))
		e.bucketCPUIdleMs.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CPUIdleMs[len(bucketStats.Op.Samples.CPUIdleMs)-1]))
		e.bucketCPULocalMs.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CPULocalMs[len(bucketStats.Op.Samples.CPULocalMs)-1]))
		e.bucketCPUUtilizationRate.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.CPUUtilizationRate[len(bucketStats.Op.Samples.CPUUtilizationRate)-1]))
		e.bucketHibernatedRequests.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.HibernatedRequests[len(bucketStats.Op.Samples.HibernatedRequests)-1]))
		e.bucketHibernatedWaked.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.HibernatedWaked[len(bucketStats.Op.Samples.HibernatedWaked)-1]))
		e.bucketMemActualFree.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemActualFree[len(bucketStats.Op.Samples.MemActualFree)-1]))
		e.bucketMemActualUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemActualUsed[len(bucketStats.Op.Samples.MemActualUsed)-1]))
		e.bucketMemFree.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemFree[len(bucketStats.Op.Samples.MemFree)-1]))
		e.bucketMemTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemTotal[len(bucketStats.Op.Samples.MemTotal)-1]))
		e.bucketMemUsedSys.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.MemUsedSys[len(bucketStats.Op.Samples.MemUsedSys)-1]))
		e.bucketRestRequests.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.RestRequests[len(bucketStats.Op.Samples.RestRequests)-1]))
		e.bucketSwapTotal.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.SwapTotal[len(bucketStats.Op.Samples.SwapTotal)-1]))
		e.bucketSwapUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.SwapUsed[len(bucketStats.Op.Samples.SwapUsed)-1]))

		if e.context.CouchbaseVersion == "4.5.1" {
			e.bucketEpTapRebalanceCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceCount[len(bucketStats.Op.Samples.EpTapRebalanceCount)-1]))
			e.bucketEpTapRebalanceQlen.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQlen[len(bucketStats.Op.Samples.EpTapRebalanceQlen)-1]))
			e.bucketEpTapRebalanceQueueBackfillremaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQueueBackfillremaining[len(bucketStats.Op.Samples.EpTapRebalanceQueueBackfillremaining)-1]))
			e.bucketEpTapRebalanceQueueBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQueueBackoff[len(bucketStats.Op.Samples.EpTapRebalanceQueueBackoff)-1]))
			e.bucketEpTapRebalanceQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQueueDrain[len(bucketStats.Op.Samples.EpTapRebalanceQueueDrain)-1]))
			e.bucketEpTapRebalanceQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQueueFill[len(bucketStats.Op.Samples.EpTapRebalanceQueueFill)-1]))
			e.bucketEpTapRebalanceQueueItemondisk.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceQueueItemondisk[len(bucketStats.Op.Samples.EpTapRebalanceQueueItemondisk)-1]))
			e.bucketEpTapRebalanceTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapRebalanceTotalBacklogSize[len(bucketStats.Op.Samples.EpTapRebalanceTotalBacklogSize)-1]))
			e.bucketEpTapReplicaCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaCount[len(bucketStats.Op.Samples.EpTapReplicaCount)-1]))
			e.bucketEpTapReplicaQlen.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQlen[len(bucketStats.Op.Samples.EpTapReplicaQlen)-1]))
			e.bucketEpTapReplicaQueueBackfillremaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQueueBackfillremaining[len(bucketStats.Op.Samples.EpTapReplicaQueueBackfillremaining)-1]))
			e.bucketEpTapReplicaQueueBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQueueBackoff[len(bucketStats.Op.Samples.EpTapReplicaQueueBackoff)-1]))
			e.bucketEpTapReplicaQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQueueDrain[len(bucketStats.Op.Samples.EpTapReplicaQueueDrain)-1]))
			e.bucketEpTapReplicaQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQueueFill[len(bucketStats.Op.Samples.EpTapReplicaQueueFill)-1]))
			e.bucketEpTapReplicaQueueItemondisk.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaQueueItemondisk[len(bucketStats.Op.Samples.EpTapReplicaQueueItemondisk)-1]))
			e.bucketEpTapReplicaTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapReplicaTotalBacklogSize[len(bucketStats.Op.Samples.EpTapReplicaTotalBacklogSize)-1]))
			e.bucketEpTapTotalCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalCount[len(bucketStats.Op.Samples.EpTapTotalCount)-1]))
			e.bucketEpTapTotalQlen.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQlen[len(bucketStats.Op.Samples.EpTapTotalQlen)-1]))
			e.bucketEpTapTotalQueueBackfillremaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQueueBackfillremaining[len(bucketStats.Op.Samples.EpTapTotalQueueBackfillremaining)-1]))
			e.bucketEpTapTotalQueueBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQueueBackoff[len(bucketStats.Op.Samples.EpTapTotalQueueBackoff)-1]))
			e.bucketEpTapTotalQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQueueDrain[len(bucketStats.Op.Samples.EpTapTotalQueueDrain)-1]))
			e.bucketEpTapTotalQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQueueFill[len(bucketStats.Op.Samples.EpTapTotalQueueFill)-1]))
			e.bucketEpTapTotalQueueItemondisk.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalQueueItemondisk[len(bucketStats.Op.Samples.EpTapTotalQueueItemondisk)-1]))
			e.bucketEpTapTotalTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapTotalTotalBacklogSize[len(bucketStats.Op.Samples.EpTapTotalTotalBacklogSize)-1]))
			e.bucketEpTapUserCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserCount[len(bucketStats.Op.Samples.EpTapUserCount)-1]))
			e.bucketEpTapUserQlen.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQlen[len(bucketStats.Op.Samples.EpTapUserQlen)-1]))
			e.bucketEpTapUserQueueBackfillremaining.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQueueBackfillremaining[len(bucketStats.Op.Samples.EpTapUserQueueBackfillremaining)-1]))
			e.bucketEpTapUserQueueBackoff.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQueueBackoff[len(bucketStats.Op.Samples.EpTapUserQueueBackoff)-1]))
			e.bucketEpTapUserQueueDrain.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQueueDrain[len(bucketStats.Op.Samples.EpTapUserQueueDrain)-1]))
			e.bucketEpTapUserQueueFill.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQueueFill[len(bucketStats.Op.Samples.EpTapUserQueueFill)-1]))
			e.bucketEpTapUserQueueItemondisk.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserQueueItemondisk[len(bucketStats.Op.Samples.EpTapUserQueueItemondisk)-1]))
			e.bucketEpTapUserTotalBacklogSize.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpTapUserTotalBacklogSize[len(bucketStats.Op.Samples.EpTapUserTotalBacklogSize)-1]))
		}
		if e.context.CouchbaseVersion == "5.1.1" {
			e.bucketAvgActiveTimestampDrift.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.AvgActiveTimestampDrift[len(bucketStats.Op.Samples.AvgActiveTimestampDrift)-1]))
			e.bucketAvgReplicaTimestampDrift.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.AvgReplicaTimestampDrift[len(bucketStats.Op.Samples.AvgReplicaTimestampDrift)-1]))
			e.bucketEpActiveAheadExceptions.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpActiveAheadExceptions[len(bucketStats.Op.Samples.EpActiveAheadExceptions)-1]))
			e.bucketEpActiveHlcDrift.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpActiveHlcDrift[len(bucketStats.Op.Samples.EpActiveHlcDrift)-1]))
			e.bucketEpActiveHlcDriftCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpActiveHlcDriftCount[len(bucketStats.Op.Samples.EpActiveHlcDriftCount)-1]))
			e.bucketEpClockCasDriftThresholdExceeded.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpClockCasDriftThresholdExceeded[len(bucketStats.Op.Samples.EpClockCasDriftThresholdExceeded)-1]))
			e.bucketEpReplicaAheadExceptions.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpReplicaAheadExceptions[len(bucketStats.Op.Samples.EpReplicaAheadExceptions)-1]))
			e.bucketEpReplicaHlcDrift.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpReplicaHlcDrift[len(bucketStats.Op.Samples.EpReplicaHlcDrift)-1]))
			e.bucketEpReplicaHlcDriftCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucketStats.Op.Samples.EpReplicaHlcDriftCount[len(bucketStats.Op.Samples.EpReplicaHlcDriftCount)-1]))
		}
	}

}
