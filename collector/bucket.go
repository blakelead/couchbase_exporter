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

// BucketData (/pools/default/buckets)
type BucketData struct {
	Name       string `json:"name"`
	BasicStats struct {
		QuotaPercentUsed float64 `json:"quotaPercentUsed"`
		OpsPerSec        float64 `json:"opsPerSec"`
		DiskFetches      float64 `json:"diskFetches"`
		ItemCount        float64 `json:"itemCount"`
		DiskUsed         float64 `json:"diskUsed"`
		DataUsed         float64 `json:"dataUsed"`
		MemUsed          float64 `json:"memUsed"`
	} `json:"basicStats"`
}

// BucketExporter describes the exporter object.
type BucketExporter struct {
	context Context
	route   string
	metrics []*GaugeVecStruct
}

// NewBucketExporter instantiates the Exporter with the URI and metrics.
func NewBucketExporter(context Context) (*BucketExporter, error) {
	bucketMetrics := GetMetricsFromFile("bucket", context)
	var metrics []*GaugeVecStruct
	for _, m := range bucketMetrics.List {
		metrics = append(metrics, newGaugeVecStruct(m.Name, m.ID, m.Description, m.Labels))
	}
	return &BucketExporter{
		context: context,
		route:   bucketMetrics.Route,
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *BucketExporter) Describe(ch chan<- *p.Desc) {
	for _, m := range e.metrics {
		m.gaugeVec.Describe(ch)
	}
}

// Collect fetches data for each exported metric.
func (e *BucketExporter) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()

	body := Fetch(e.context, e.route)
	var buckets []BucketData
	err := json.Unmarshal(body, &buckets)
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, bucket := range buckets {
		flat := FlattenStruct(bucket, "")
		for _, m := range e.metrics {
			if value, ok := flat[m.id].(float64); ok {
				m.gaugeVec.With(p.Labels{"bucket": bucket.Name}).Set(value)
			}
		}
	}
	for _, m := range e.metrics {
		m.gaugeVec.Collect(ch)
	}

	mutex.Unlock()
}
