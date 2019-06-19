// Copyright 2019 Adel Abdelhak.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// XDCRExporter encapsulates XDCR metrics and context.
type XDCRExporter struct {
	context    Context
	route      string
	errorCount *p.GaugeVec
	metrics    map[string]*p.Desc
}

// NewXDCRExporter creates the XDCRExporter and fill it with metrics metadata from the metrics file.
func NewXDCRExporter(c Context) (*XDCRExporter, error) {
	xdcrMetrics, err := GetMetricsFromFile("xdcr")
	if err != nil {
		return &XDCRExporter{}, err
	}
	// metrics is a map where the key is the metric ID and the value is a Prometheus Descriptor for that metric.
	metrics := make(map[string]*p.Desc, len(xdcrMetrics.List))
	for _, metric := range xdcrMetrics.List {
		fqName := p.BuildFQName("cb", xdcrMetrics.Name, metric.Name)
		metrics[metric.ID] = p.NewDesc(fqName, metric.Description, metric.Labels, nil)
	}
	return &XDCRExporter{
		context: c,
		route:   xdcrMetrics.Route,
		errorCount: p.NewGaugeVec(p.GaugeOpts{
			Name: p.BuildFQName("cb", xdcrMetrics.Name, "error_count"),
			Help: "Number of XDCR errors",
		}, []string{"remote_cluster_id", "remote_cluster_name", "source_bucket", "destination_bucket"}),
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics
func (e *XDCRExporter) Describe(ch chan<- *p.Desc) {
	e.errorCount.Describe(ch)
	for _, metric := range e.metrics {
		ch <- metric
	}
}

// Collect fetches data for each exported metric
func (e *XDCRExporter) Collect(ch chan<- p.Metric) {
	// Get task list to retrieve active XDCR links.
	body, err := Fetch(e.context, "/pools/default/tasks")
	if err != nil {
		log.Error("Could not retrieve tasks data: XDCR metrics won't be scraped")
		return
	}
	var tasks []struct {
		Type   string   `json:"type"`
		Status string   `json:"status"`
		ID     string   `json:"id"`
		Errors []string `json:"errors"`
	}
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		log.Error("Could not unmarshal tasks data: XDCR metrics won't be scraped")
		return
	}

	var routes []string
	remoteClusters := make(map[string]string, 0)
	errorsCount := make(map[string]int, 0)
	for _, task := range tasks {
		if task.Type == "xdcr" {
			// Create URL for each XDCR metric.
			taskID := strings.Split(task.ID, "/")
			if len(taskID) < 3 {
				log.Error("Task ID doesn't have the expected format (uuid/src/dest): ", taskID)
				continue
			}
			uuid, src, dest := taskID[0], taskID[1], taskID[2]
			for metricID := range e.metrics {
				route := fmt.Sprintf("%s/%s/stats/replications%%2F%s%%2F%s%%2F%s%%2F%s", e.route, src, uuid, src, dest, metricID)
				routes = append(routes, route)
			}

			// Get error count based on number of error messages in tasks endpoint.
			errorsCount[uuid] = len(task.Errors)

			// Associate remote clusters names with uuid for labelling.
			body, err = Fetch(e.context, "/pools/default/remoteClusters")
			if err != nil {
				log.Error("Could not retrieve remote clusters data")
			}
			var tmpRemoteClusters []struct {
				Name string `json:"name"`
				UUID string `json:"uuid"`
			}
			err = json.Unmarshal(body, &tmpRemoteClusters)
			if err != nil {
				log.Error("Could not unmarshal remote clusters data")
			}
			for _, rc := range tmpRemoteClusters {
				remoteClusters[rc.UUID] = rc.Name
			}
		}
	}

	// Get hostname of the node.
	body, err = Fetch(e.context, "/nodes/self")
	if err != nil {
		log.Error("Could not retrieve node data: XDCR metrics won't be scraped")
		return
	}
	var node struct {
		Hostname string `json:"hostname"`
	}
	err = json.Unmarshal(body, &node)
	if err != nil {
		log.Error("Could not unmarshal node data: XDCR metrics won't be scraped")
		return
	}

	// Fetch all bodies from urls created above.
	bodies := MultiFetch(e.context, routes)

	var currentUUID string
	for route, body := range bodies {
		// Split back url to get uuid src & dest buckets and metric name.
		longID := strings.Split(route, "%2F")
		uuid, src, dest, metricID := longID[1], longID[2], longID[3], longID[4]

		// Get node stats object.
		var xdcr struct {
			NodeStats map[string]interface{} `json:"nodeStats"`
		}
		err := json.Unmarshal(body, &xdcr)
		if err != nil {
			log.Error("Could not unmarshal XDCR data for remote " + uuid + " and metric " + metricID)
			continue
		}
		if _, ok := xdcr.NodeStats[node.Hostname].([]interface{}); !ok {
			continue
		}

		list := xdcr.NodeStats[node.Hostname].([]interface{})
		if len(list) == 0 {
			log.Debug("No value found for " + metricID + " metric in remote " + uuid)
			continue
		}

		var value float64
		switch v := list[len(list)-1].(type) {
		case float64:
			value = v
		case int:
			value = float64(v)
		}

		if currentUUID != uuid {
			currentUUID = uuid
			e.errorCount.WithLabelValues(uuid, remoteClusters[uuid], src, dest).Set(float64(errorsCount[uuid]))
			e.errorCount.Collect(ch)
		}

		ch <- p.MustNewConstMetric(e.metrics[metricID], p.GaugeValue, value, uuid, remoteClusters[uuid], src, dest)
	}
}
