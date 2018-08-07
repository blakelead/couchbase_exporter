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
}

// URI is a custom url wrapper with credentials
type URI struct {
	URL      string
	Username string
	Password string
}

func newGaugeVec(name string, help string, labels []string) *p.GaugeVec {
	return p.NewGaugeVec(p.GaugeOpts{Namespace: "cb", Name: name, Help: help}, labels)
}

// NewExporter instantiates the Exporter with the URI and metrics.
func NewExporter(uri URI) (*Exporter, error) {
	return &Exporter{
		uri: uri,
		up: p.NewGauge(p.GaugeOpts{
			Namespace: "cb",
			Name:      "up",
			Help:      "State of last scrape.",
		}),
		totalScrapes: p.NewCounter(p.CounterOpts{
			Namespace: "cb",
			Name:      "total_scrapes",
			Help:      "Total number of scrapes.",
		}),

		nodesStatus:            newGaugeVec("node_status", "Status of couchbase node.", nil),
		nodesClusterMembership: newGaugeVec("node_cluster_membership", "Status of node cluster membership.", nil),
		nodeCPUUtilizationRate: newGaugeVec("node_cpu_utilization_rate", "CPU utilization rate.", nil),
		nodeRAMUsed:            newGaugeVec("node_ram_usage_bytes", "RAM used per node in bytes.", nil),

		clusterRAMTotal:              newGaugeVec("cluster_ram_total_bytes", "Total RAM in the cluster.", nil),
		clusterRAMUsed:               newGaugeVec("cluster_ram_used_bytes", "Used RAM in the cluster.", nil),
		clusterRAMUsedByData:         newGaugeVec("cluster_ram_used_by_data_bytes", "Used RAM by data in the cluster.", nil),
		clusterRAMQuotaTotal:         newGaugeVec("cluster_ram_quota_total_bytes", "Total quota RAM in the cluster.", nil),
		clusterRAMQuotaTotalPerNode:  newGaugeVec("cluster_ram_quota_total_per_node_bytes", "Total quota RAM per node in the cluster.", nil),
		clusterRAMQuotaUsed:          newGaugeVec("cluster_ram_quota_used_bytes", "Used quota RAM in the cluster.", nil),
		clusterRAMQuotaUsedPerNode:   newGaugeVec("cluster_ram_quota_used_per_node_bytes", "Used quota RAM per node in the cluster.", nil),
		clusterDiskTotal:             newGaugeVec("cluster_disk_total_bytes", "Total disk in the cluster.", nil),
		clusterDiskQuotaTotal:        newGaugeVec("cluster_disk_quota_total_bytes", "Disk quota in the cluster.", nil),
		clusterDiskUsed:              newGaugeVec("cluster_disk_used_bytes", "Used disk in the cluster.", nil),
		clusterDiskUsedByData:        newGaugeVec("cluster_disk_used_by_data_bytes", "Disk used by data in the cluster.", nil),
		clusterDiskFree:              newGaugeVec("cluster_disk_free_bytes", "Free disk in the cluster", nil),
		clusterFtsRAMQuota:           newGaugeVec("cluster_fts_ram_quota_bytes", "RAM quota for Full text search bucket.", nil),
		clusterIndexRAMQuota:         newGaugeVec("cluster_index_ram_quota_bytes", "RAM quota for Index bucket.", nil),
		clusterRAMQuota:              newGaugeVec("cluster_data_ram_quota_bytes", "RAM quota for Data bucket.", nil),
		clusterRebalanceStatus:       newGaugeVec("cluster_rebalance_status", "Occurrence of rebalancing in the cluster.", nil),
		clusterMaxBucketCount:        newGaugeVec("cluster_max_bucket_count", "Maximum number of buckets.", nil),
		clusterFailoverNodeCount:     newGaugeVec("cluster_failover_node_count", "Number of failovers since cluster is up.", nil),
		clusterRebalanceSuccessCount: newGaugeVec("cluster_rebalance_success_count", "Number of rebalance success since cluster is up.", nil),
		clusterRebalanceStartCount:   newGaugeVec("cluster_rebalance_start_count", "Number of rebalance start since cluster is up.", nil),
		clusterRebalanceFailCount:    newGaugeVec("cluster_rebalance_fail_count", "Number of rebalance failure since cluster is up.", nil),
	}, nil
}

// Describe describes exported metrics.
func (e *Exporter) Describe(ch chan<- *p.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()

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
}

// Collect fetches data for each exported metric.
func (e *Exporter) Collect(ch chan<- p.Metric) {
	e.mutex.Lock()

	e.totalScrapes.Inc()
	e.scrapeUp()
	e.scrapeNodes()

	ch <- e.up
	ch <- e.totalScrapes

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

	e.nodesStatus.Collect(ch)
	e.nodesClusterMembership.Collect(ch)
	e.nodeCPUUtilizationRate.Collect(ch)
	e.nodeRAMUsed.Collect(ch)

	e.mutex.Unlock()
}

func (e *Exporter) scrapeUp() {
	e.up.Set(0)
	req, err := http.NewRequest("HEAD", e.uri.URL, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode == 200 {
		e.up.Set(1)
	}
	log.Debug("HEAD " + e.uri.URL + " - response: " + res.Status)
	res.Body.Close()
}

func reset(e *Exporter) {
	e.nodesStatus.Reset()
	e.nodesClusterMembership.Reset()
	e.nodeCPUUtilizationRate.Reset()
	e.nodeRAMUsed.Reset()

	e.clusterRAMUsed.Reset()
	e.clusterRAMUsedByData.Reset()
	e.clusterRAMQuotaTotal.Reset()
	e.clusterRAMQuotaTotalPerNode.Reset()
	e.clusterRAMQuotaUsed.Reset()
	e.clusterRAMQuotaUsedPerNode.Reset()
	e.clusterDiskTotal.Reset()
	e.clusterDiskQuotaTotal.Reset()
	e.clusterDiskUsed.Reset()
	e.clusterDiskUsedByData.Reset()
	e.clusterDiskFree.Reset()
	e.clusterFtsRAMQuota.Reset()
	e.clusterIndexRAMQuota.Reset()
	e.clusterRAMQuota.Reset()
	e.clusterRebalanceStatus.Reset()
	e.clusterMaxBucketCount.Reset()
	e.clusterFailoverNodeCount.Reset()
	e.clusterRebalanceSuccessCount.Reset()
	e.clusterRebalanceStartCount.Reset()
	e.clusterRebalanceFailCount.Reset()
}

func (e *Exporter) scrapeNodes() {
	reset(e)

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
