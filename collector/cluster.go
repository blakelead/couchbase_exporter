// Copyright 2019 Adel Abdelhak.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// ClusterData (/pools/default)
type ClusterData struct {
	StorageTotals struct {
		RAM struct {
			Total      float64 `json:"total"`
			QuotaTotal float64 `json:"quotaTotal"`
			QuotaUsed  float64 `json:"quotaUsed"`
			Used       float64 `json:"used"`
			UsedByData float64 `json:"usedByData"`
		} `json:"ram"`
		Hdd struct {
			Total      float64 `json:"total"`
			QuotaTotal float64 `json:"quotaTotal"`
			Used       float64 `json:"used"`
			UsedByData float64 `json:"usedByData"`
			Free       float64 `json:"free"`
		} `json:"hdd"`
	} `json:"storageTotals"`
	FtsMemoryQuota   float64 `json:"ftsMemoryQuota"`
	IndexMemoryQuota float64 `json:"indexMemoryQuota"`
	MemoryQuota      float64 `json:"memoryQuota"`
	RebalanceStatus  string  `json:"rebalanceStatus"`
	MaxBucketCount   float64 `json:"maxBucketCount"`
	Counters         struct {
		FailoverNode     float64 `json:"failover_node"`
		RebalanceSuccess float64 `json:"rebalance_success"`
		RebalanceStart   float64 `json:"rebalance_start"`
		RebalanceFail    float64 `json:"rebalance_fail"`
	} `json:"counters"`
	Balanced bool `json:"balanced"` // couchbase 5.1.1
}

// ClusterExporter encapsulates cluster metrics and context.
type ClusterExporter struct {
	context      Context
	route        string
	totalScrapes p.Counter
	metrics      map[string]*CustomDesc
}

// NewClusterExporter creates the ClusterExporter and fill it with metrics metadata from the metrics file.
func NewClusterExporter(context Context) (*ClusterExporter, error) {
	clusterMetrics, err := GetMetricsFromFile("cluster")
	if err != nil {
		return &ClusterExporter{}, err
	}
	// metrics is a map where the key is the metric ID and the value is a Prometheus Descriptor for that metric.
	metrics := make(map[string]*CustomDesc, len(clusterMetrics.List))
	for _, metric := range clusterMetrics.List {
		fqName := p.BuildFQName("cb", clusterMetrics.Name, metric.Name)
		metrics[metric.ID] = newCustomDesc(fqName, metric.Description, metric.Labels, metric.Type)
	}
	return &ClusterExporter{
		context: context,
		route:   clusterMetrics.Route,
		totalScrapes: p.NewCounter(p.CounterOpts{
			Name: p.BuildFQName("cb", clusterMetrics.Name, "scrapes_total"),
			Help: "Number of scrapes since the start of the exporter.",
		}),
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *ClusterExporter) Describe(ch chan<- *p.Desc) {
	ch <- e.totalScrapes.Desc()
	for _, metric := range e.metrics {
		ch <- metric.pDesc
	}
}

// Collect fetches data for each exported metric.
func (e *ClusterExporter) Collect(ch chan<- p.Metric) {
	e.totalScrapes.Inc()
	ch <- e.totalScrapes

	body, err := Fetch(e.context, e.route)
	if err != nil {
		log.Error("Error when retrieving cluster data. Cluster metrics won't be scraped")
		return
	}
	var cluster ClusterData
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		log.Error("Could not unmarshal cluster data")
		return
	}

	flat := FlattenStruct(cluster)
	for id, metric := range e.metrics {
		if value, ok := flat[id]; ok {
			switch value.(type) {
			case bool:
				var v float64
				// Rebalance status
				if value.(bool) {
					v = 1
				}
				ch <- p.MustNewConstMetric(metric.pDesc, getValueType(metric.mType), v)
			case string:
				var v float64
				if value.(string) != "none" {
					v = 1
				}
				ch <- p.MustNewConstMetric(metric.pDesc, getValueType(metric.mType), v)
			case float64:
				ch <- p.MustNewConstMetric(metric.pDesc, getValueType(metric.mType), value.(float64))
			}
		}
	}
}
