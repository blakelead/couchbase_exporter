// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	p "github.com/prometheus/client_golang/prometheus"
)

var (
	clusterRoute = "/pools/default"
	nodeRoute    = "/nodes/self"
	bucketRoute  = "/pools/default/buckets"
)

// Context is a custom url wrapper with credentials
type Context struct {
	URI              string
	Username         string
	Password         string
	CouchbaseVersion string
}

// Exporters regroups all exporters structs
type Exporters struct {
	Cluster     *ClusterExporter
	Node        *NodeExporter
	Bucket      *BucketExporter
	BucketStats *BucketStatsExporter
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
func NewExporters(context Context) (*Exporters, error) {
	clusterExporter, _ := NewClusterExporter(context)
	nodeExporter, _ := NewNodeExporter(context)
	bucketExporter, _ := NewBucketExporter(context)
	bucketStatsExporter, _ := NewBucketStatsExporter(context)

	return &Exporters{clusterExporter, nodeExporter, bucketExporter, bucketStatsExporter}, nil
}
