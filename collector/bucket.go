// Copyright 2019 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"

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

// BucketExporter encapsulates bucket metrics and context.
type BucketExporter struct {
	context Context
	route   string
	metrics map[string]*p.Desc
}

// NewBucketExporter creates the BucketExporter and fill it with metrics metadata from the metrics file.
func NewBucketExporter(context Context) (*BucketExporter, error) {
	bucketMetrics, err := GetMetricsFromFile("bucket")
	if err != nil {
		return &BucketExporter{}, err
	}
	// metrics is a map where the key is the metric ID and the value is a Prometheus Descriptor for that metric.
	metrics := make(map[string]*p.Desc, len(bucketMetrics.List))
	for _, metric := range bucketMetrics.List {
		fqName := p.BuildFQName("cb", bucketMetrics.Name, metric.Name)
		metrics[metric.ID] = p.NewDesc(fqName, metric.Description, metric.Labels, nil)
	}
	return &BucketExporter{
		context: context,
		route:   bucketMetrics.Route,
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *BucketExporter) Describe(ch chan<- *p.Desc) {
	for _, metric := range e.metrics {
		ch <- metric
	}
}

// Collect fetches data for each exported metric.
func (e *BucketExporter) Collect(ch chan<- p.Metric) {
	body, err := Fetch(e.context, e.route)
	if err != nil {
		log.Error("Error when retrieving buckets data. Buckets metrics won't be scraped")
		return
	}
	var buckets []BucketData
	err = json.Unmarshal(body, &buckets)
	if err != nil {
		log.Error("Could not unmarshal buckets data")
		return
	}

	for _, bucket := range buckets {
		flat := FlattenStruct(bucket)
		for id, metric := range e.metrics {
			if value, ok := flat[id].(float64); ok {
				ch <- p.MustNewConstMetric(metric, p.GaugeValue, value, bucket.Name)
			}
		}
	}
}
