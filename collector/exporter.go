// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Exporter describes the exporter object.
type Exporter struct {
	mutex sync.RWMutex
	uri   URI

	totalScrapes p.Counter

	up                           p.Gauge
	clusterRAMTotal              p.Gauge
	clusterRAMUsed               p.Gauge
	clusterRAMUsedByData         p.Gauge
	clusterRAMQuotaTotal         p.Gauge
	clusterRAMQuotaTotalPerNode  p.Gauge
	clusterRAMQuotaUsed          p.Gauge
	clusterRAMQuotaUsedPerNode   p.Gauge
	clusterDiskTotal             p.Gauge
	clusterDiskQuotaTotal        p.Gauge
	clusterDiskUsed              p.Gauge
	clusterDiskUsedByData        p.Gauge
	clusterDiskFree              p.Gauge
	clusterFtsRAMQuota           p.Gauge
	clusterIndexRAMQuota         p.Gauge
	clusterRAMQuota              p.Gauge
	clusterRebalanceStatus       p.Gauge
	clusterMaxBucketCount        p.Gauge
	clusterFailoverNodeCount     p.Gauge
	clusterRebalanceSuccessCount p.Gauge
	clusterRebalanceStartCount   p.Gauge
	clusterRebalanceFailCount    p.Gauge

	nodeRAMTotal                 p.Gauge
	nodeRAMUsed                  p.Gauge
	nodeRAMUsedByData            p.Gauge
	nodeRAMQuotaTotal            p.Gauge
	nodeRAMQuotaTotalPerNode     p.Gauge
	nodeRAMQuotaUsed             p.Gauge
	nodeRAMQuotaUsedPerNode      p.Gauge
	nodeDiskTotal                p.Gauge
	nodeDiskQuotaTotal           p.Gauge
	nodeDiskUsed                 p.Gauge
	nodeDiskUsedByData           p.Gauge
	nodeDiskFree                 p.Gauge
	nodeCPUUtilizationRate       p.Gauge
	nodeSwapTotal                p.Gauge
	nodeSwapUsed                 p.Gauge
	nodeCmdGet                   p.Gauge
	nodeCouchDocsActualDiskSize  p.Gauge
	nodeCouchDocsDataSize        p.Gauge
	nodeCouchSpatialDataSize     p.Gauge
	nodeCouchSpatialDiskSize     p.Gauge
	nodeCouchViewsActualDiskSize p.Gauge
	nodeCouchViewsDataSize       p.Gauge
	nodeCurrItems                p.Gauge
	nodeCurrItemsTot             p.Gauge
	nodeEpBgFetched              p.Gauge
	nodeGetHits                  p.Gauge
	nodeMemUsed                  p.Gauge
	nodeOps                      p.Gauge
	nodeVbReplicaCurrItems       p.Gauge
	nodeUptime                   p.Gauge
	nodesClusterMembership       p.Gauge
	nodesStatus                  p.Gauge
	nodeFtsRAMQuota              p.Gauge
	nodeIndexRAMQuota            p.Gauge
	nodeRAMQuota                 p.Gauge

	bucketProxyPort        *p.GaugeVec
	bucketReplicaIndex     *p.GaugeVec
	bucketReplicaNumber    *p.GaugeVec
	bucketThreadsNumber    *p.GaugeVec
	bucketRAMQuota         *p.GaugeVec
	bucketRawRAMQuota      *p.GaugeVec
	bucketQuotaPercentUsed *p.GaugeVec
	bucketOpsPerSec        *p.GaugeVec
	bucketDiskFetches      *p.GaugeVec
	bucketItemCount        *p.GaugeVec
	bucketDiskUsed         *p.GaugeVec
	bucketDataUsed         *p.GaugeVec
	bucketMemUsed          *p.GaugeVec
}

// URI is a custom url wrapper with credentials
type URI struct {
	URL      string
	Username string
	Password string
}

func newGauge(name string, help string) p.Gauge {
	return p.NewGauge(p.GaugeOpts{Namespace: "cb", Name: name, Help: help})
}

func newGaugeVec(name string, help string, labels []string) *p.GaugeVec {
	return p.NewGaugeVec(p.GaugeOpts{Namespace: "cb", Name: name, Help: help}, labels)
}

// NewExporter instantiates the Exporter with the URI and metrics.
func NewExporter(uri URI) (*Exporter, error) {
	return &Exporter{
		uri: uri,

		totalScrapes: p.NewCounter(p.CounterOpts{Namespace: "cb", Name: "total_scrapes", Help: "Total number of scrapes."}),

		up:                           newGauge("up", "State of cluster."),
		clusterRAMTotal:              newGauge("cluster_ram_total_bytes", "Total RAM in the cluster."),
		clusterRAMUsed:               newGauge("cluster_ram_used_bytes", "Used RAM in the cluster."),
		clusterRAMUsedByData:         newGauge("cluster_ram_used_by_data_bytes", "Used RAM by data in the cluster."),
		clusterRAMQuotaTotal:         newGauge("cluster_ram_quota_total_bytes", "Total quota RAM in the cluster."),
		clusterRAMQuotaTotalPerNode:  newGauge("cluster_ram_quota_total_per_node_bytes", "Total quota RAM per node in the cluster."),
		clusterRAMQuotaUsed:          newGauge("cluster_ram_quota_used_bytes", "Used quota RAM in the cluster."),
		clusterRAMQuotaUsedPerNode:   newGauge("cluster_ram_quota_used_per_node_bytes", "Used quota RAM per node in the cluster."),
		clusterDiskTotal:             newGauge("cluster_disk_total_bytes", "Total disk in the cluster."),
		clusterDiskQuotaTotal:        newGauge("cluster_disk_quota_total_bytes", "Disk quota in the cluster."),
		clusterDiskUsed:              newGauge("cluster_disk_used_bytes", "Used disk in the cluster."),
		clusterDiskUsedByData:        newGauge("cluster_disk_used_by_data_bytes", "Disk used by data in the cluster."),
		clusterDiskFree:              newGauge("cluster_disk_free_bytes", "Free disk in the cluster"),
		clusterFtsRAMQuota:           newGauge("cluster_fts_ram_quota_bytes", "RAM quota for Full text search bucket."),
		clusterIndexRAMQuota:         newGauge("cluster_index_ram_quota_bytes", "RAM quota for Index bucket."),
		clusterRAMQuota:              newGauge("cluster_data_ram_quota_bytes", "RAM quota for Data bucket."),
		clusterRebalanceStatus:       newGauge("cluster_rebalance_status", "Occurrence of rebalancing in the cluster."),
		clusterMaxBucketCount:        newGauge("cluster_max_bucket_count", "Maximum number of buckets."),
		clusterFailoverNodeCount:     newGauge("cluster_failover_node_count", "Number of failovers since cluster is up."),
		clusterRebalanceSuccessCount: newGauge("cluster_rebalance_success_count", "Number of rebalance success since cluster is up."),
		clusterRebalanceStartCount:   newGauge("cluster_rebalance_start_count", "Number of rebalance start since cluster is up."),
		clusterRebalanceFailCount:    newGauge("cluster_rebalance_fail_count", "Number of rebalance failure since cluster is up."),

		nodeRAMTotal:                 newGauge("node_ram_total_bytes", "Node total RAM."),
		nodeRAMUsed:                  newGauge("node_ram_usage_bytes", "Node used RAM."),
		nodeRAMUsedByData:            newGauge("node_ram_used_by_data_bytes", "Node RAM used by data."),
		nodeRAMQuotaTotal:            newGauge("node_ram_quota_total_bytes", "Node RAM quota total."),
		nodeRAMQuotaUsed:             newGauge("node_ram_quota_used_bytes", "Node RAM quota used."),
		nodeDiskTotal:                newGauge("node_disk_total_bytes", "Node Total disk."),
		nodeDiskQuotaTotal:           newGauge("node_disk_quota_total_bytes", "Node disk quota total."),
		nodeDiskUsed:                 newGauge("node_disk_used_bytes", "Node used disk."),
		nodeDiskUsedByData:           newGauge("node_disk_quota_used_bytes", "Node disk quota total."),
		nodeDiskFree:                 newGauge("node_disk_free_bytes", "Node free disk."),
		nodeCPUUtilizationRate:       newGauge("node_cpu_utilization_rate", "CPU utilization rate."),
		nodeSwapTotal:                newGauge("node_swap_total_bytes", "Node total swap."),
		nodeSwapUsed:                 newGauge("node_swap_used_bytes", "Node used swap."),
		nodeCmdGet:                   newGauge("node_stats_cmd_get", "Node stats: cmd_get."),
		nodeCouchDocsActualDiskSize:  newGauge("node_stats_couch_docs_actual_disk_size", "Node stats: couch_docs_actual_disk_size."),
		nodeCouchDocsDataSize:        newGauge("node_stats_couch_docs_data_size", "Node stats: couch_docs_data_size."),
		nodeCouchSpatialDataSize:     newGauge("node_stats_couch_spatial_data_size", "Node stats: couch_spatial_data_size."),
		nodeCouchSpatialDiskSize:     newGauge("node_stats_couch_spatial_disk_size", "Node stats: couch_spatial_disk_size."),
		nodeCouchViewsActualDiskSize: newGauge("node_stats_couch_views_actual_disk_size", "Node stats: couch_views_actual_disk_size"),
		nodeCouchViewsDataSize:       newGauge("node_stats_couch_views_data_size", "Node stats: couch_views_data_size."),
		nodeCurrItems:                newGauge("node_stats_curr_items", "Node stats: curr_items."),
		nodeCurrItemsTot:             newGauge("node_stats_curr_items_tot", "Node stats: curr_items_tot."),
		nodeEpBgFetched:              newGauge("node_stats_ep_bg_fetched", "Node stats: ep_bg_fetched."),
		nodeGetHits:                  newGauge("node_stats_get_hits", "Node stats: get_hits."),
		nodeMemUsed:                  newGauge("node_stats_mem_used", "Node stats: mem_used."),
		nodeOps:                      newGauge("node_stats_ops", "Node stats: ops."),
		nodeVbReplicaCurrItems:       newGauge("node_stats_vb_replica_curr_items", "Node stats: vb_replica_curr_items."),
		nodeUptime:                   newGauge("node_uptime_seconds", "Node uptime."),
		nodesClusterMembership:       newGauge("node_cluster_membership", "Status of node cluster membership."),
		nodesStatus:                  newGauge("node_status", "Status of couchbase node."),
		nodeFtsRAMQuota:              newGauge("node_fts_ram_quota_bytes", "Node quota for Full text search bucket."),
		nodeIndexRAMQuota:            newGauge("node_index_ram_quota_bytes", "Node quota for Index bucket."),
		nodeRAMQuota:                 newGauge("node_data_ram_quota_bytes", "Node quota for Data bucket."),

		bucketProxyPort:        newGaugeVec("bucket_proxy_port", "Bucket proxy port.", []string{"bucket"}),
		bucketReplicaIndex:     newGaugeVec("bucket_replica_index", "Bucket replica index.", []string{"bucket"}),
		bucketReplicaNumber:    newGaugeVec("bucket_replica_number", "Bucket replica number.", []string{"bucket"}),
		bucketThreadsNumber:    newGaugeVec("bucket_threads_number", "Bucket thread number.", []string{"bucket"}),
		bucketRAMQuota:         newGaugeVec("bucket_ram_quota_bytes", "Bucket RAM quota.", []string{"bucket"}),
		bucketRawRAMQuota:      newGaugeVec("bucket_raw_ram_quota_bytes", "Bucket raw RAM quota.", []string{"bucket"}),
		bucketQuotaPercentUsed: newGaugeVec("bucket_quota_percent_used", "Bucket quota usage.", []string{"bucket"}),
		bucketOpsPerSec:        newGaugeVec("bucket_ops_per_second", "Bucket operations per second.", []string{"bucket"}),
		bucketDiskFetches:      newGaugeVec("bucket_disk_fetches", "Bucket disk fetches.", []string{"bucket"}),
		bucketItemCount:        newGaugeVec("bucket_item_count", "Bucket item count.", []string{"bucket"}),
		bucketDiskUsed:         newGaugeVec("bucket_disk_used_bytes", "Bucket disk used.", []string{"bucket"}),
		bucketDataUsed:         newGaugeVec("bucket_data_used_bytes", "Bucket data used.", []string{"bucket"}),
		bucketMemUsed:          newGaugeVec("bucket_ram_used_bytes", "Bucket RAM used.", []string{"bucket"}),
	}, nil
}

// Describe describes exported metrics.
func (e *Exporter) Describe(ch chan<- *p.Desc) {
	e.totalScrapes.Describe(ch)

	e.up.Describe(ch)
	e.clusterRAMTotal.Describe(ch)
	e.clusterRAMUsed.Describe(ch)
	e.clusterRAMUsedByData.Describe(ch)
	e.clusterRAMQuotaTotal.Describe(ch)
	e.clusterRAMQuotaTotalPerNode.Describe(ch)
	e.clusterRAMQuotaUsed.Describe(ch)
	e.clusterRAMQuotaUsedPerNode.Describe(ch)
	e.clusterDiskTotal.Describe(ch)
	e.clusterDiskQuotaTotal.Describe(ch)
	e.clusterDiskUsed.Describe(ch)
	e.clusterDiskUsedByData.Describe(ch)
	e.clusterDiskFree.Describe(ch)
	e.clusterFtsRAMQuota.Describe(ch)
	e.clusterIndexRAMQuota.Describe(ch)
	e.clusterRAMQuota.Describe(ch)
	e.clusterRebalanceStatus.Describe(ch)
	e.clusterMaxBucketCount.Describe(ch)
	e.clusterFailoverNodeCount.Describe(ch)
	e.clusterRebalanceSuccessCount.Describe(ch)
	e.clusterRebalanceStartCount.Describe(ch)
	e.clusterRebalanceFailCount.Describe(ch)

	e.nodeRAMTotal.Describe(ch)
	e.nodeRAMUsed.Describe(ch)
	e.nodeRAMUsedByData.Describe(ch)
	e.nodeRAMQuotaTotal.Describe(ch)
	e.nodeRAMQuotaUsed.Describe(ch)
	e.nodeDiskTotal.Describe(ch)
	e.nodeDiskQuotaTotal.Describe(ch)
	e.nodeDiskUsed.Describe(ch)
	e.nodeDiskUsedByData.Describe(ch)
	e.nodeDiskFree.Describe(ch)
	e.nodeCPUUtilizationRate.Describe(ch)
	e.nodeSwapTotal.Describe(ch)
	e.nodeSwapUsed.Describe(ch)
	e.nodeCmdGet.Describe(ch)
	e.nodeCouchDocsActualDiskSize.Describe(ch)
	e.nodeCouchDocsDataSize.Describe(ch)
	e.nodeCouchSpatialDataSize.Describe(ch)
	e.nodeCouchSpatialDiskSize.Describe(ch)
	e.nodeCouchViewsActualDiskSize.Describe(ch)
	e.nodeCouchViewsDataSize.Describe(ch)
	e.nodeCurrItems.Describe(ch)
	e.nodeCurrItemsTot.Describe(ch)
	e.nodeEpBgFetched.Describe(ch)
	e.nodeGetHits.Describe(ch)
	e.nodeMemUsed.Describe(ch)
	e.nodeOps.Describe(ch)
	e.nodeVbReplicaCurrItems.Describe(ch)
	e.nodeUptime.Describe(ch)
	e.nodesClusterMembership.Describe(ch)
	e.nodesStatus.Describe(ch)
	e.nodeFtsRAMQuota.Describe(ch)
	e.nodeIndexRAMQuota.Describe(ch)
	e.nodeRAMQuota.Describe(ch)

	e.bucketProxyPort.Describe(ch)
	e.bucketReplicaIndex.Describe(ch)
	e.bucketReplicaNumber.Describe(ch)
	e.bucketThreadsNumber.Describe(ch)
	e.bucketRAMQuota.Describe(ch)
	e.bucketRawRAMQuota.Describe(ch)
	e.bucketQuotaPercentUsed.Describe(ch)
	e.bucketOpsPerSec.Describe(ch)
	e.bucketDiskFetches.Describe(ch)
	e.bucketItemCount.Describe(ch)
	e.bucketDiskUsed.Describe(ch)
	e.bucketDataUsed.Describe(ch)
	e.bucketMemUsed.Describe(ch)

}

// Collect fetches data for each exported metric.
func (e *Exporter) Collect(ch chan<- p.Metric) {
	e.mutex.Lock()

	e.totalScrapes.Inc()
	e.scrapeClusterData()
	e.scrapeNodeData()
	e.scrapeBucketData()

	e.totalScrapes.Collect(ch)

	e.up.Collect(ch)
	e.clusterRAMTotal.Collect(ch)
	e.clusterRAMUsed.Collect(ch)
	e.clusterRAMUsedByData.Collect(ch)
	e.clusterRAMQuotaTotal.Collect(ch)
	e.clusterRAMQuotaTotalPerNode.Collect(ch)
	e.clusterRAMQuotaUsed.Collect(ch)
	e.clusterRAMQuotaUsedPerNode.Collect(ch)
	e.clusterDiskTotal.Collect(ch)
	e.clusterDiskQuotaTotal.Collect(ch)
	e.clusterDiskUsed.Collect(ch)
	e.clusterDiskUsedByData.Collect(ch)
	e.clusterDiskFree.Collect(ch)
	e.clusterFtsRAMQuota.Collect(ch)
	e.clusterIndexRAMQuota.Collect(ch)
	e.clusterRAMQuota.Collect(ch)
	e.clusterRebalanceStatus.Collect(ch)
	e.clusterMaxBucketCount.Collect(ch)
	e.clusterFailoverNodeCount.Collect(ch)
	e.clusterRebalanceSuccessCount.Collect(ch)
	e.clusterRebalanceStartCount.Collect(ch)
	e.clusterRebalanceFailCount.Collect(ch)

	e.nodeRAMTotal.Collect(ch)
	e.nodeRAMUsed.Collect(ch)
	e.nodeRAMUsedByData.Collect(ch)
	e.nodeRAMQuotaTotal.Collect(ch)
	e.nodeRAMQuotaUsed.Collect(ch)
	e.nodeDiskTotal.Collect(ch)
	e.nodeDiskQuotaTotal.Collect(ch)
	e.nodeDiskUsed.Collect(ch)
	e.nodeDiskUsedByData.Collect(ch)
	e.nodeDiskFree.Collect(ch)
	e.nodeCPUUtilizationRate.Collect(ch)
	e.nodeSwapTotal.Collect(ch)
	e.nodeSwapUsed.Collect(ch)
	e.nodeCmdGet.Collect(ch)
	e.nodeCouchDocsActualDiskSize.Collect(ch)
	e.nodeCouchDocsDataSize.Collect(ch)
	e.nodeCouchSpatialDataSize.Collect(ch)
	e.nodeCouchSpatialDiskSize.Collect(ch)
	e.nodeCouchViewsActualDiskSize.Collect(ch)
	e.nodeCouchViewsDataSize.Collect(ch)
	e.nodeCurrItems.Collect(ch)
	e.nodeCurrItemsTot.Collect(ch)
	e.nodeEpBgFetched.Collect(ch)
	e.nodeGetHits.Collect(ch)
	e.nodeMemUsed.Collect(ch)
	e.nodeOps.Collect(ch)
	e.nodeVbReplicaCurrItems.Collect(ch)
	e.nodeUptime.Collect(ch)
	e.nodesClusterMembership.Collect(ch)
	e.nodesStatus.Collect(ch)
	e.nodeFtsRAMQuota.Collect(ch)
	e.nodeIndexRAMQuota.Collect(ch)
	e.nodeRAMQuota.Collect(ch)

	e.bucketProxyPort.Collect(ch)
	e.bucketReplicaIndex.Collect(ch)
	e.bucketReplicaNumber.Collect(ch)
	e.bucketThreadsNumber.Collect(ch)
	e.bucketRAMQuota.Collect(ch)
	e.bucketRawRAMQuota.Collect(ch)
	e.bucketQuotaPercentUsed.Collect(ch)
	e.bucketOpsPerSec.Collect(ch)
	e.bucketDiskFetches.Collect(ch)
	e.bucketItemCount.Collect(ch)
	e.bucketDiskUsed.Collect(ch)
	e.bucketDataUsed.Collect(ch)
	e.bucketMemUsed.Collect(ch)

	e.mutex.Unlock()
}

func (e *Exporter) scrapeClusterData() {
	e.up.Set(0)

	req, err := http.NewRequest("GET", e.uri.URL+"/pools/default", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.uri.Username, e.uri.Password)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Error(req.URL.Path + ": " + res.Status)
		return
	}

	var cluster ClusterData
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &cluster)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.uri.URL + "/pools/default" + " - data: " + string(body))

	var rebalance int
	if cluster.RebalanceStatus == "running" {
		rebalance = 1
	}
	e.up.Set(1)
	e.clusterRAMTotal.Set(float64(cluster.StorageTotals.RAM.Total))
	e.clusterRAMUsed.Set(float64(cluster.StorageTotals.RAM.Used))
	e.clusterRAMUsedByData.Set(float64(cluster.StorageTotals.RAM.UsedByData))
	e.clusterRAMQuotaTotal.Set(float64(cluster.StorageTotals.RAM.QuotaTotal))
	e.clusterRAMQuotaTotalPerNode.Set(float64(cluster.StorageTotals.RAM.QuotaTotalPerNode))
	e.clusterRAMQuotaUsed.Set(float64(cluster.StorageTotals.RAM.QuotaUsed))
	e.clusterRAMQuotaUsedPerNode.Set(float64(cluster.StorageTotals.RAM.QuotaUsedPerNode))
	e.clusterDiskTotal.Set(float64(cluster.StorageTotals.Hdd.Total))
	e.clusterDiskQuotaTotal.Set(float64(cluster.StorageTotals.Hdd.QuotaTotal))
	e.clusterDiskUsed.Set(float64(cluster.StorageTotals.Hdd.Used))
	e.clusterDiskUsedByData.Set(float64(cluster.StorageTotals.Hdd.UsedByData))
	e.clusterDiskFree.Set(float64(cluster.StorageTotals.Hdd.Free))
	e.clusterFtsRAMQuota.Set(float64(cluster.FtsMemoryQuota * 1024 * 1024))
	e.clusterIndexRAMQuota.Set(float64(cluster.IndexMemoryQuota * 1024 * 1024))
	e.clusterRAMQuota.Set(float64(cluster.MemoryQuota * 1024 * 1024))
	e.clusterRebalanceStatus.Set(float64(rebalance))
	e.clusterMaxBucketCount.Set(float64(cluster.MaxBucketCount))
	e.clusterFailoverNodeCount.Set(float64(cluster.Counters.FailoverNode))
	e.clusterRebalanceSuccessCount.Set(float64(cluster.Counters.RebalanceSuccess))
	e.clusterRebalanceStartCount.Set(float64(cluster.Counters.RebalanceStart))
	e.clusterRebalanceFailCount.Set(float64(cluster.Counters.RebalanceFail))
}

func (e *Exporter) scrapeNodeData() {
	req, err := http.NewRequest("GET", e.uri.URL+"/nodes/self", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.uri.Username, e.uri.Password)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Error(req.URL.Path + ": " + res.Status)
		return
	}

	var node NodeData
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &node)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.uri.URL + "/nodes/self" + " - data: " + string(body))

	var status int
	if node.Status == "healthy" {
		status = 1
	}
	var membership int
	if node.ClusterMembership == "active" {
		membership = 1
	}
	uptime, err := strconv.Atoi(node.Uptime)
	if err != nil {
		log.Error(err.Error())
	}

	e.nodeRAMTotal.Set(float64(node.StorageTotals.RAM.Total))
	e.nodeRAMUsed.Set(float64(node.StorageTotals.RAM.Used))
	e.nodeRAMUsedByData.Set(float64(node.StorageTotals.RAM.UsedByData))
	e.nodeRAMQuotaTotal.Set(float64(node.StorageTotals.RAM.QuotaTotal))
	e.nodeRAMQuotaUsed.Set(float64(node.StorageTotals.RAM.QuotaUsed))
	e.nodeDiskTotal.Set(float64(node.StorageTotals.Hdd.Total))
	e.nodeDiskQuotaTotal.Set(float64(node.StorageTotals.Hdd.QuotaTotal))
	e.nodeDiskUsed.Set(float64(node.StorageTotals.Hdd.Used))
	e.nodeDiskUsedByData.Set(float64(node.StorageTotals.Hdd.UsedByData))
	e.nodeDiskFree.Set(float64(node.StorageTotals.Hdd.Free))
	e.nodeCPUUtilizationRate.Set(float64(node.SystemStats.CPUUtilizationRate))
	e.nodeSwapTotal.Set(float64(node.SystemStats.SwapTotal))
	e.nodeSwapUsed.Set(float64(node.SystemStats.SwapUsed))
	e.nodeCmdGet.Set(float64(node.InterestingStats.CmdGet))
	e.nodeCouchDocsActualDiskSize.Set(float64(node.InterestingStats.CouchDocsActualDiskSize))
	e.nodeCouchDocsDataSize.Set(float64(node.InterestingStats.CouchDocsDataSize))
	e.nodeCouchSpatialDataSize.Set(float64(node.InterestingStats.CouchSpatialDataSize))
	e.nodeCouchSpatialDiskSize.Set(float64(node.InterestingStats.CouchSpatialDiskSize))
	e.nodeCouchViewsActualDiskSize.Set(float64(node.InterestingStats.CouchViewsActualDiskSize))
	e.nodeCouchViewsDataSize.Set(float64(node.InterestingStats.CouchViewsDataSize))
	e.nodeCurrItems.Set(float64(node.InterestingStats.CurrItems))
	e.nodeCurrItemsTot.Set(float64(node.InterestingStats.CurrItemsTot))
	e.nodeEpBgFetched.Set(float64(node.InterestingStats.EpBgFetched))
	e.nodeGetHits.Set(float64(node.InterestingStats.GetHits))
	e.nodeMemUsed.Set(float64(node.InterestingStats.MemUsed))
	e.nodeOps.Set(float64(node.InterestingStats.Ops))
	e.nodeVbReplicaCurrItems.Set(float64(node.InterestingStats.VbReplicaCurrItems))
	e.nodeUptime.Set(float64(uptime))
	e.nodesClusterMembership.Set(float64(membership))
	e.nodesStatus.Set(float64(status))
	e.nodeFtsRAMQuota.Set(float64(node.FtsMemoryQuota * 1024 * 1024))
	e.nodeIndexRAMQuota.Set(float64(node.IndexMemoryQuota * 1024 * 1024))
	e.nodeRAMQuota.Set(float64(node.FtsMemoryQuota * 1024 * 1024))
}

func (e *Exporter) scrapeBucketData() {
	req, err := http.NewRequest("GET", e.uri.URL+"/pools/default/buckets", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.uri.Username, e.uri.Password)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Error(req.URL.Path + ": " + res.Status)
		return
	}

	var buckets []BucketData
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &buckets)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.uri.URL + "/pools/default/buckets" + " - data: " + string(body))

	for _, bucket := range buckets {
		var replicaIndex int
		if bucket.ReplicaIndex == true {
			replicaIndex = 1
		}
		e.bucketProxyPort.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ProxyPort))
		e.bucketReplicaIndex.With(p.Labels{"bucket": bucket.Name}).Set(float64(replicaIndex))
		e.bucketReplicaNumber.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ReplicaNumber))
		e.bucketThreadsNumber.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ThreadsNumber))
		e.bucketRAMQuota.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.Quota.RAM))
		e.bucketRawRAMQuota.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.Quota.RawRAM))
		e.bucketQuotaPercentUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.QuotaPercentUsed))
		e.bucketOpsPerSec.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.OpsPerSec))
		e.bucketDiskFetches.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskFetches))
		e.bucketItemCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.ItemCount))
		e.bucketDiskUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskUsed))
		e.bucketDataUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskUsed))
		e.bucketMemUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.MemUsed))

	}

}
