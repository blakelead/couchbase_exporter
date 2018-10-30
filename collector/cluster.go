// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"sync"

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

// ClusterExporter describes the exporter object.
type ClusterExporter struct {
	context      Context
	route        string
	totalScrapes p.Counter
	metrics      []*GaugeVecStruct
}

// NewClusterExporter instantiates the Exporter with the URI and metrics.
func NewClusterExporter(context Context) (*ClusterExporter, error) {
	clusterMetrics := GetMetricsFromFile("cluster", context)
	var metrics []*GaugeVecStruct
	for _, m := range clusterMetrics.List {
		metrics = append(metrics, newGaugeVecStruct(m.Name, m.ID, m.Description, m.Labels))
	}
	return &ClusterExporter{
		context:      context,
		route:        clusterMetrics.Route,
		totalScrapes: newCounter("total_scrapes", "Total number of scrapes"),
		metrics:      metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *ClusterExporter) Describe(ch chan<- *p.Desc) {
	e.totalScrapes.Describe(ch)
	for _, m := range e.metrics {
		m.gaugeVec.Describe(ch)
	}
}

// Collect fetches data for each exported metric.
func (e *ClusterExporter) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()
	e.totalScrapes.Inc()

	body := Fetch(e.context, e.route)
	var cluster ClusterData
	err := json.Unmarshal(body, &cluster)
	if err != nil {
		log.Error(err.Error())
		return
	}

	flat := FlattenStruct(cluster, "")
	for _, m := range e.metrics {
		if value, ok := flat[m.id]; ok {
			switch value.(type) {
			case bool:
				var v float64
				if value.(bool) {
					v = 1
				}
				m.gaugeVec.With(p.Labels{}).Set(v)
			case string:
				var v float64
				if value.(string) != "none" {
					v = 1
				}
				m.gaugeVec.With(p.Labels{}).Set(v)
			case float64:
				m.gaugeVec.With(p.Labels{}).Set(value.(float64))
			}
		}
	}

	for _, m := range e.metrics {
		m.gaugeVec.Collect(ch)
	}
	mutex.Unlock()
}
