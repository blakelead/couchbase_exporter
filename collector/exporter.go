// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Exporter describes the exporter object.
type Exporter struct {
	mutex        sync.RWMutex
	uri          URI
	up           p.Gauge
	totalScrapes p.Counter

	clusterRAMTotal              *p.GaugeVec
	clusterRAMUsed               *p.GaugeVec
	clusterRAMUsedByData         *p.GaugeVec
	clusterRAMQuotaTotal         *p.GaugeVec
	clusterRAMQuotaTotalPerNode  *p.GaugeVec
	clusterRAMQuotaUsed          *p.GaugeVec
	clusterRAMQuotaUsedPerNode   *p.GaugeVec
	clusterDiskTotal             *p.GaugeVec
	clusterDiskQuotaTotal        *p.GaugeVec
	clusterDiskUsed              *p.GaugeVec
	clusterDiskUsedByData        *p.GaugeVec
	clusterDiskFree              *p.GaugeVec
	clusterFtsRAMQuota           *p.GaugeVec
	clusterIndexRAMQuota         *p.GaugeVec
	clusterRAMQuota              *p.GaugeVec
	clusterRebalanceStatus       *p.GaugeVec
	clusterMaxBucketCount        *p.GaugeVec
	clusterFailoverNodeCount     *p.GaugeVec
	clusterRebalanceSuccessCount *p.GaugeVec
	clusterRebalanceStartCount   *p.GaugeVec
	clusterRebalanceFailCount    *p.GaugeVec

	nodesStatus            *p.GaugeVec
	nodesClusterMembership *p.GaugeVec
	nodeCPUUtilizationRate *p.GaugeVec
	nodeRAMUsed            *p.GaugeVec
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

	e.nodesStatus.Describe(ch)
	e.nodesClusterMembership.Describe(ch)
	e.nodeCPUUtilizationRate.Describe(ch)
	e.nodeRAMUsed.Describe(ch)
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

	var data ClusterData
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.uri.URL + "/pools/default" + " - data: " + string(body))

	getClusterData(e, &data)
	getNodeData(e, &data)
	// getBucketData()
}

func getClusterData(e *Exporter, data *ClusterData) {
	var rebalance int
	if data.RebalanceStatus == "running" {
		rebalance = 1
	}
	e.clusterRAMTotal.With(nil).Set(float64(data.StorageTotals.RAM.Total))
	e.clusterRAMUsed.With(nil).Set(float64(data.StorageTotals.RAM.Used))
	e.clusterRAMUsedByData.With(nil).Set(float64(data.StorageTotals.RAM.UsedByData))
	e.clusterRAMQuotaTotal.With(nil).Set(float64(data.StorageTotals.RAM.QuotaTotal))
	e.clusterRAMQuotaTotalPerNode.With(nil).Set(float64(data.StorageTotals.RAM.QuotaTotalPerNode))
	e.clusterRAMQuotaUsed.With(nil).Set(float64(data.StorageTotals.RAM.QuotaUsed))
	e.clusterRAMQuotaUsedPerNode.With(nil).Set(float64(data.StorageTotals.RAM.QuotaUsedPerNode))
	e.clusterDiskTotal.With(nil).Set(float64(data.StorageTotals.Hdd.Total))
	e.clusterDiskQuotaTotal.With(nil).Set(float64(data.StorageTotals.Hdd.QuotaTotal))
	e.clusterDiskUsed.With(nil).Set(float64(data.StorageTotals.Hdd.Used))
	e.clusterDiskUsedByData.With(nil).Set(float64(data.StorageTotals.Hdd.UsedByData))
	e.clusterDiskFree.With(nil).Set(float64(data.StorageTotals.Hdd.Free))
	e.clusterFtsRAMQuota.With(nil).Set(float64(data.FtsMemoryQuota * 1024 * 1024))
	e.clusterIndexRAMQuota.With(nil).Set(float64(data.IndexMemoryQuota * 1024 * 1024))
	e.clusterRAMQuota.With(nil).Set(float64(data.MemoryQuota * 1024 * 1024))
	e.clusterRebalanceStatus.With(nil).Set(float64(rebalance))
	e.clusterMaxBucketCount.With(nil).Set(float64(data.MaxBucketCount))
	e.clusterFailoverNodeCount.With(nil).Set(float64(data.Counters.FailoverNode))
	e.clusterRebalanceSuccessCount.With(nil).Set(float64(data.Counters.RebalanceSuccess))
	e.clusterRebalanceStartCount.With(nil).Set(float64(data.Counters.RebalanceStart))
	e.clusterRebalanceFailCount.With(nil).Set(float64(data.Counters.RebalanceFail))
}

func getNodeData(e *Exporter, data *ClusterData) {
	for _, n := range data.Nodes {
		var status int
		if n.Status == "healthy" {
			status = 1
		}
		var membership int
		if n.ClusterMembership == "active" {
			membership = 1
		}
		e.nodesStatus.With(nil).Set(float64(status))
		e.nodesClusterMembership.With(nil).Set(float64(membership))
		e.nodeCPUUtilizationRate.With(nil).Set(n.SystemStats.CPUUtilizationRate)
		e.nodeRAMUsed.With(nil).Set(float64(n.InterestingStats.MemUsed))
	}
}

func getBucketData(e *Exporter, data *BucketData) {

}
