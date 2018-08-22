// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// ClusterData (/pools/default)
type ClusterData struct {
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
	FtsMemoryQuota   int    `json:"ftsMemoryQuota"`
	IndexMemoryQuota int    `json:"indexMemoryQuota"`
	MemoryQuota      int    `json:"memoryQuota"`
	RebalanceStatus  string `json:"rebalanceStatus"`
	MaxBucketCount   int    `json:"maxBucketCount"`
	Counters         struct {
		FailoverNode     int `json:"failover_node"`
		RebalanceSuccess int `json:"rebalance_success"`
		RebalanceStart   int `json:"rebalance_start"`
		RebalanceFail    int `json:"rebalance_fail"`
	} `json:"counters"`
}

// ClusterExporter describes the exporter object.
type ClusterExporter struct {
	uri                          URI
	totalScrapes                 p.Counter
	clusterRAMTotal              p.Gauge
	clusterRAMUsed               p.Gauge
	clusterRAMUsedByData         p.Gauge
	clusterRAMQuotaTotal         p.Gauge
	clusterRAMQuotaUsed          p.Gauge
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
}

// NewClusterExporter instantiates the Exporter with the URI and metrics.
func NewClusterExporter(uri URI) (*ClusterExporter, error) {
	return &ClusterExporter{
		uri:                          uri,
		totalScrapes:                 newCounter("total_scrapes", "Total number of scrapes"),
		clusterRAMTotal:              newGauge("cluster_ram_total_bytes", "Total memory available to the cluster"),
		clusterRAMUsed:               newGauge("cluster_ram_used_bytes", "Memory used by the cluster"),
		clusterRAMUsedByData:         newGauge("cluster_ram_used_by_data_bytes", "Memory used by the data in the cluster"),
		clusterRAMQuotaTotal:         newGauge("cluster_ram_quota_total_bytes", "Total memory allocated to Couchbase in the cluster"),
		clusterRAMQuotaUsed:          newGauge("cluster_ram_quota_used_bytes", "Memory quota used by the cluster"),
		clusterDiskTotal:             newGauge("cluster_disk_total_bytes", "Total disk space available to the cluster"),
		clusterDiskUsed:              newGauge("cluster_disk_used_bytes", "Disk space used by the cluster"),
		clusterDiskQuotaTotal:        newGauge("cluster_disk_quota_total_bytes", "Disk space quota for the cluster"),
		clusterDiskUsedByData:        newGauge("cluster_disk_used_by_data_bytes", "Disk space used by the data in the cluster"),
		clusterDiskFree:              newGauge("cluster_disk_free_bytes", "Free disk space in the cluster"),
		clusterFtsRAMQuota:           newGauge("cluster_fts_ram_quota_bytes", "Memory quota allocated to full text search buckets"),
		clusterIndexRAMQuota:         newGauge("cluster_index_ram_quota_bytes", "Memory quota allocated to Index buckets"),
		clusterRAMQuota:              newGauge("cluster_data_ram_quota_bytes", "Memory quota allocated to Data buckets"),
		clusterRebalanceStatus:       newGauge("cluster_rebalance_status", "Rebalance status. 1:rebalancing"),
		clusterMaxBucketCount:        newGauge("cluster_max_bucket_count", "Maximum number of buckets allowed"),
		clusterFailoverNodeCount:     newGauge("cluster_failover_node_count", "Number of failovers since cluster is up"),
		clusterRebalanceSuccessCount: newGauge("cluster_rebalance_success_count", "Number of rebalance successes since cluster is up"),
		clusterRebalanceStartCount:   newGauge("cluster_rebalance_start_count", "Number of rebalance starts since cluster is up"),
		clusterRebalanceFailCount:    newGauge("cluster_rebalance_fail_count", "Number of rebalance fails since cluster is up"),
	}, nil
}

// Describe describes exported metrics.
func (e *ClusterExporter) Describe(ch chan<- *p.Desc) {

	e.totalScrapes.Describe(ch)
	e.clusterRAMTotal.Describe(ch)
	e.clusterRAMUsed.Describe(ch)
	e.clusterRAMUsedByData.Describe(ch)
	e.clusterRAMQuotaTotal.Describe(ch)
	e.clusterRAMQuotaUsed.Describe(ch)
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
}

// Collect fetches data for each exported metric.
func (e *ClusterExporter) Collect(ch chan<- p.Metric) {
	e.totalScrapes.Inc()
	e.scrape()
	e.totalScrapes.Collect(ch)
	e.clusterRAMTotal.Collect(ch)
	e.clusterRAMUsed.Collect(ch)
	e.clusterRAMUsedByData.Collect(ch)
	e.clusterRAMQuotaTotal.Collect(ch)
	e.clusterRAMQuotaUsed.Collect(ch)
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
}

func (e *ClusterExporter) scrape() {
	req, err := http.NewRequest("GET", e.uri.URL+"/pools/default", nil)
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

	rebalance := 1
	if cluster.RebalanceStatus == "node" {
		rebalance = 0
	}
	e.clusterRAMTotal.Set(float64(cluster.StorageTotals.RAM.Total))
	e.clusterRAMUsed.Set(float64(cluster.StorageTotals.RAM.Used))
	e.clusterRAMUsedByData.Set(float64(cluster.StorageTotals.RAM.UsedByData))
	e.clusterRAMQuotaTotal.Set(float64(cluster.StorageTotals.RAM.QuotaTotal))
	e.clusterRAMQuotaUsed.Set(float64(cluster.StorageTotals.RAM.QuotaUsed))
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
