// Copyright 2018 Blake Lead.
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
	uri                    URI
	up                     p.Gauge
	totalScrapes           p.Counter
	nodesStatus            *p.GaugeVec
	nodesClusterMembership *p.GaugeVec
	cpuUtilizationRate     *p.GaugeVec
	ramUsage               *p.GaugeVec
	clusterRAMTotal        *p.GaugeVec
	mutex                  sync.RWMutex
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
		nodesStatus:            newGaugeVec("node_status", "Status of couchbase node.", []string{"hostname"}),
		nodesClusterMembership: newGaugeVec("node_cluster_membership", "Status of node cluster membership.", []string{"hostname"}),
		cpuUtilizationRate:     newGaugeVec("node_cpu_utilization_rate", "CPU utilization rate.", []string{"hostname"}),
		ramUsage:               newGaugeVec("node_ram_usage_bytes", "RAM used per node in bytes.", []string{"hostname"}),
		clusterRAMTotal:        newGaugeVec("cluster_ram_total_bytes", "Total RAM in the cluster.", nil),
	}, nil
}

// Describe describes exported metrics.
func (e *Exporter) Describe(ch chan<- *p.Desc) {
	ch <- e.up.Desc()
	ch <- e.totalScrapes.Desc()
	e.nodesStatus.Describe(ch)
	e.nodesClusterMembership.Describe(ch)
	e.cpuUtilizationRate.Describe(ch)
	e.ramUsage.Describe(ch)
	e.clusterRAMTotal.Describe(ch)
}

// Collect fetches data for each exported metric.
func (e *Exporter) Collect(ch chan<- p.Metric) {
	e.mutex.Lock()
	e.totalScrapes.Inc()
	e.scrapeUp()
	e.scrapeNodes()
	ch <- e.up
	ch <- e.totalScrapes
	e.nodesStatus.Collect(ch)
	e.nodesClusterMembership.Collect(ch)
	e.cpuUtilizationRate.Collect(ch)
	e.ramUsage.Collect(ch)
	e.clusterRAMTotal.Collect(ch)
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

func (e *Exporter) scrapeNodes() {
	e.nodesStatus.Reset()
	e.nodesClusterMembership.Reset()
	e.cpuUtilizationRate.Reset()
	e.ramUsage.Reset()

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
	e.clusterRAMTotal.With(nil).Set(float64(data.StorageTotals.RAM.Total))
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
		// Node status
		e.nodesStatus.With(p.Labels{"hostname": n.Hostname}).Set(float64(status))
		// Node cluster membership
		e.nodesClusterMembership.With(p.Labels{"hostname": n.Hostname}).Set(float64(membership))
		// Node CPU usage
		e.cpuUtilizationRate.With(p.Labels{"hostname": n.Hostname}).Set(n.SystemStats.CPUUtilizationRate)
		// Node memory usage
		e.ramUsage.With(p.Labels{"hostname": n.Hostname}).Set(float64(n.InterestingStats.MemUsed))
	}
}

func getBucketData(e *Exporter, data *BucketData) {

}
