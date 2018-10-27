// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"strconv"
	"sync"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// NodeData (/nodes/self)
type NodeData struct {
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
	SystemStats struct {
		CPUUtilizationRate float64 `json:"cpu_utilization_rate"`
		SwapTotal          float64 `json:"swap_total"`
		SwapUsed           float64 `json:"swap_used"`
	} `json:"systemStats"`
	InterestingStats struct {
		CmdGet                       float64 `json:"cmd_get"`
		CouchDocsActualDiskSize      float64 `json:"couch_docs_actual_disk_size"`
		CouchDocsDataSize            float64 `json:"couch_docs_data_size"`
		CouchSpatialDataSize         float64 `json:"couch_spatial_data_size"`
		CouchSpatialDiskSize         float64 `json:"couch_spatial_disk_size"`
		CouchViewsActualDiskSize     float64 `json:"couch_views_actual_disk_size"`
		CouchViewsDataSize           float64 `json:"couch_views_data_size"`
		CurrItems                    float64 `json:"curr_items"`
		CurrItemsTot                 float64 `json:"curr_items_tot"`
		EpBgFetched                  float64 `json:"ep_bg_fetched"`
		GetHits                      float64 `json:"get_hits"`
		MemUsed                      float64 `json:"mem_used"`
		Ops                          float64 `json:"ops"`
		VbReplicaCurrItems           float64 `json:"vb_replica_curr_items"`
		VbActiveNumNonResidentNumber float64 `json:"vb_active_num_non_residentNumber"` // couchbase 5.1.1
	} `json:"interestingStats"`
	Uptime            string  `json:"uptime"`
	ClusterMembership string  `json:"clusterMembership"`
	Status            string  `json:"status"`
	FtsMemoryQuota    float64 `json:"ftsMemoryQuota"`
	IndexMemoryQuota  float64 `json:"indexMemoryQuota"`
	MemoryQuota       float64 `json:"memoryQuota"`
}

// NodeExporter describes the exporter object.
type NodeExporter struct {
	context Context
	route   string
	up      *p.GaugeVec
	metrics []*GaugeVecStruct
}

// NewNodeExporter instantiates the Exporter with the URI and metrics.
func NewNodeExporter(context Context) (*NodeExporter, error) {
	nodeMetrics := GetMetricsFromFile("node", context)
	var metrics []*GaugeVecStruct
	for _, m := range nodeMetrics.List {
		metrics = append(metrics, newGaugeVecStruct(m.Name, m.ID, m.Description, m.Labels))
	}
	return &NodeExporter{
		context: context,
		route:   nodeMetrics.Route,
		up:      newGaugeVec("node_service_up", "Couchbase service healthcheck", []string{}),
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *NodeExporter) Describe(ch chan<- *p.Desc) {
	e.up.Describe(ch)
	for _, m := range e.metrics {
		m.gaugeVec.Describe(ch)
	}
}

// Collect fetches data for each exported metric.
func (e *NodeExporter) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()
	e.up.With(p.Labels{}).Set(0)

	body := Fetch(e.context, e.route)
	var node NodeData
	err := json.Unmarshal(body, &node)
	if err != nil {
		log.Error(err.Error())
		return
	}

	statusValues := map[string]float64{"healthy": 1, "active": 1, "warmup": 2, "inactiveHealthy": 2, "inactiveFailed": 3}

	e.up.With(p.Labels{}).Set(1)
	flat := FlattenStruct(node, "")
	for _, m := range e.metrics {
		if value, ok := flat[m.id]; ok {
			switch value.(type) {
			case string:
				var v float64
				uptime, err := strconv.Atoi(value.(string))
				if err == nil {
					v = float64(uptime)
				} else {
					v = statusValues[value.(string)]
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
