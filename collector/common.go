// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"sync"

	p "github.com/prometheus/client_golang/prometheus"
)

// URI is a custom url wrapper with credentials
type URI struct {
	URL      string
	Username string
	Password string
}

// Exporters regroups all exporters structs
type Exporters struct {
	Cluster *ClusterExporter
	Node    *NodeExporter
	Bucket  *BucketExporter
}

func newCounter(name string, help string) p.Counter {
	return p.NewCounter(p.CounterOpts{Namespace: "cb", Name: name, Help: help})
}

func newGauge(name string, help string) p.Gauge {
	return p.NewGauge(p.GaugeOpts{Namespace: "cb", Name: name, Help: help})
}

func newGaugeVec(name string, help string, labels []string) *p.GaugeVec {
	return p.NewGaugeVec(p.GaugeOpts{Namespace: "cb", Name: name, Help: help}, labels)
}

// NewExporters instantiates the Exporter with the URI and metrics.
func NewExporters(uri URI) (*Exporters, error) {
	clusterExporter, _ := NewClusterExporter(uri)
	nodeExporter, _ := NewNodeExporter(uri)
	bucketExporter, _ := NewBucketExporter(uri)

	return &Exporters{clusterExporter, nodeExporter, bucketExporter}, nil
}

// Describe describes exported metrics.
func (e *Exporters) Describe(ch chan<- *p.Desc) {
	e.Cluster.Describe(ch)
	e.Node.Describe(ch)
	e.Bucket.Describe(ch)
}

// Collect fetches data for each exported metric.
func (e *Exporters) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()

	e.Cluster.Collect(ch)
	e.Node.Collect(ch)
	e.Bucket.Collect(ch)

	mutex.Unlock()
}
