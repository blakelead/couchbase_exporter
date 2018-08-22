// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// NodeData (/nodes/self)
type NodeData struct {
	StorageTotals struct {
		RAM struct {
			Total      int `json:"total"`
			QuotaTotal int `json:"quotaTotal"`
			QuotaUsed  int `json:"quotaUsed"`
			Used       int `json:"used"`
			UsedByData int `json:"usedByData"`
		} `json:"ram"`
		Hdd struct {
			Total      int64 `json:"total"`
			QuotaTotal int64 `json:"quotaTotal"`
			Used       int64 `json:"used"`
			UsedByData int   `json:"usedByData"`
			Free       int64 `json:"free"`
		} `json:"hdd"`
	} `json:"storageTotals"`
	SystemStats struct {
		CPUUtilizationRate float64 `json:"cpu_utilization_rate"`
		SwapTotal          int     `json:"swap_total"`
		SwapUsed           int     `json:"swap_used"`
	} `json:"systemStats"`
	InterestingStats struct {
		CmdGet                   int `json:"cmd_get"`
		CouchDocsActualDiskSize  int `json:"couch_docs_actual_disk_size"`
		CouchDocsDataSize        int `json:"couch_docs_data_size"`
		CouchSpatialDataSize     int `json:"couch_spatial_data_size"`
		CouchSpatialDiskSize     int `json:"couch_spatial_disk_size"`
		CouchViewsActualDiskSize int `json:"couch_views_actual_disk_size"`
		CouchViewsDataSize       int `json:"couch_views_data_size"`
		CurrItems                int `json:"curr_items"`
		CurrItemsTot             int `json:"curr_items_tot"`
		EpBgFetched              int `json:"ep_bg_fetched"`
		GetHits                  int `json:"get_hits"`
		MemUsed                  int `json:"mem_used"`
		Ops                      int `json:"ops"`
		VbReplicaCurrItems       int `json:"vb_replica_curr_items"`
	} `json:"interestingStats"`
	Uptime               string   `json:"uptime"`
	McdMemoryReserved    int      `json:"mcdMemoryReserved"`
	McdMemoryAllocated   int      `json:"mcdMemoryAllocated"`
	ClusterMembership    string   `json:"clusterMembership"`
	RecoveryType         string   `json:"recoveryType"`
	Status               string   `json:"status"`
	Hostname             string   `json:"hostname"`
	ClusterCompatibility int      `json:"clusterCompatibility"`
	Version              string   `json:"version"`
	Os                   string   `json:"os"`
	Services             []string `json:"services"`
	FtsMemoryQuota       int      `json:"ftsMemoryQuota"`
	IndexMemoryQuota     int      `json:"indexMemoryQuota"`
	MemoryQuota          int      `json:"memoryQuota"`
}

// NodeExporter describes the exporter object.
type NodeExporter struct {
	uri                          URI
	up                           p.Gauge
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
}

// NewNodeExporter instantiates the Exporter with the URI and metrics.
func NewNodeExporter(uri URI) (*NodeExporter, error) {
	return &NodeExporter{
		uri:                          uri,
		up:                           newGauge("node_service_up", "Couchbase service healthcheck"),
		nodeRAMTotal:                 newGauge("node_ram_total_bytes", "Total memory available to the node"),
		nodeRAMUsed:                  newGauge("node_ram_usage_bytes", "Memory used by the node"),
		nodeRAMUsedByData:            newGauge("node_ram_used_by_data_bytes", "Memory used by data in the node"),
		nodeRAMQuotaTotal:            newGauge("node_ram_quota_total_bytes", "Memory quota allocated to the node"),
		nodeRAMQuotaUsed:             newGauge("node_ram_quota_used_bytes", "Memory quota used by the node"),
		nodeDiskTotal:                newGauge("node_disk_total_bytes", "Total disk space available to the node"),
		nodeDiskQuotaTotal:           newGauge("node_disk_quota_total_bytes", "Disk space quota for the node"),
		nodeDiskUsed:                 newGauge("node_disk_used_bytes", "Disk space used by the node"),
		nodeDiskUsedByData:           newGauge("node_disk_quota_used_bytes", "Disk space quota used by the node"),
		nodeDiskFree:                 newGauge("node_disk_free_bytes", "Free disk space in the node"),
		nodeCPUUtilizationRate:       newGauge("node_cpu_utilization_rate", "CPU utilization rate in percent"),
		nodeSwapTotal:                newGauge("node_swap_total_bytes", "Total swap space allocated to the node"),
		nodeSwapUsed:                 newGauge("node_swap_used_bytes", "Amount of swap space used by the node"),
		nodeCmdGet:                   newGauge("node_stats_cmd_get", "Number of get commands"),
		nodeCouchDocsActualDiskSize:  newGauge("node_stats_couch_docs_actual_disk_size", "Disk space used by Couchbase documents"),
		nodeCouchDocsDataSize:        newGauge("node_stats_couch_docs_data_size", "Couchbase documents data size in the node"),
		nodeCouchSpatialDataSize:     newGauge("node_stats_couch_spatial_data_size", "Data size for Couchbase spatial views"),
		nodeCouchSpatialDiskSize:     newGauge("node_stats_couch_spatial_disk_size", "Disk space used by Couchbase spatial views"),
		nodeCouchViewsActualDiskSize: newGauge("node_stats_couch_views_actual_disk_size", "Disk space used by Couchbase views"),
		nodeCouchViewsDataSize:       newGauge("node_stats_couch_views_data_size", "Data size for Couchbase views"),
		nodeCurrItems:                newGauge("node_stats_curr_items", "Number of current items"),
		nodeCurrItemsTot:             newGauge("node_stats_curr_items_tot", "Total number of items in the node"),
		nodeEpBgFetched:              newGauge("node_stats_ep_bg_fetched", "Number of background disk fetches"),
		nodeGetHits:                  newGauge("node_stats_get_hits", "Number of get hits"),
		nodeMemUsed:                  newGauge("node_stats_mem_used", "Memory used by the node"),
		nodeOps:                      newGauge("node_stats_ops", "Number of operations performed in the node"),
		nodeVbReplicaCurrItems:       newGauge("node_stats_vb_replica_curr_items", "Number of replicas in current items"),
		nodeUptime:                   newGauge("node_uptime_seconds", "Node uptime"),
		nodesClusterMembership:       newGauge("node_cluster_membership", "Status of node cluster membership. 1:active, 2:inactiveAdded, 3:inactiveFailed"),
		nodesStatus:                  newGauge("node_status", "Status of couchbase node. 1:healthy, 2:warmup"),
		nodeFtsRAMQuota:              newGauge("node_fts_ram_quota_bytes", "Memory quota allocated to full text search buckets"),
		nodeIndexRAMQuota:            newGauge("node_index_ram_quota_bytes", "Memory quota allocated to index buckets"),
		nodeRAMQuota:                 newGauge("node_data_ram_quota_bytes", "Memory quota allocated to data buckets"),
	}, nil
}

// Describe describes exported metrics.
func (e *NodeExporter) Describe(ch chan<- *p.Desc) {
	e.up.Describe(ch)
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
}

// Collect fetches data for each exported metric.
func (e *NodeExporter) Collect(ch chan<- p.Metric) {
	e.scrape()
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
}

func (e *NodeExporter) scrape() {
	e.up.Set(0)
	req, err := http.NewRequest("GET", e.uri.URL+"/nodes/self", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.uri.Username, e.uri.Password)
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

	e.up.Set(1)
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
