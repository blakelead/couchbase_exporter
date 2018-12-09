// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Task (/pools/default/tasks)
type Task struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	ID     string `json:"id"`
}

// XDCRExporter describes the exporter object.
type XDCRExporter struct {
	context Context
	route   string
	metrics map[string]*p.Desc
}

// NewXDCRExporter instantiates the Exporter with the URI and metrics.
func NewXDCRExporter(context Context) (*XDCRExporter, error) {
	xdcrMetrics, err := GetMetricsFromFile("xdcr")
	if err != nil {
		return &XDCRExporter{}, err
	}
	metrics := make(map[string]*p.Desc, len(xdcrMetrics.List))
	for _, metric := range xdcrMetrics.List {
		fqName := p.BuildFQName("cb", xdcrMetrics.Name, metric.Name)
		metrics[metric.ID] = p.NewDesc(fqName, metric.Description, metric.Labels, nil)
	}
	return &XDCRExporter{
		context: context,
		route:   xdcrMetrics.Route,
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *XDCRExporter) Describe(ch chan<- *p.Desc) {
	for _, metric := range e.metrics {
		ch <- metric
	}
}

// Collect fetches data for each exported metric.
func (e *XDCRExporter) Collect(ch chan<- p.Metric) {
	// get task list where xdcr are listed
	body, err := Fetch(e.context, "/pools/default/tasks")
	if err != nil {
		log.Error("Error when retrieving XDCR data. XDCR metrics won't be scraped")
		return
	}
	var tasks []Task
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		log.Error("Could not unmarshal tasks data")
		return
	}

	// create urls from task list for each xdcr metric
	var routes []string
	for _, task := range tasks {
		if task.Type == "xdcr" {
			longID := strings.Split(task.ID, "/") // id is in the form uuid/src/dest
			for id := range e.metrics {
				route := fmt.Sprintf("%s/%s/stats/replications%%2F%s%%2F%s%%2F%s%%2F%s", e.route, longID[1], longID[0], longID[1], longID[2], id)
				routes = append(routes, route)
			}
		}
	}

	// fetch all bodies from urls created above
	bodies := MultiFetch(e.context, routes)

	for route, body := range bodies {
		longID := strings.Split(route, "%2F")
		uuid, src, dest, id := longID[1], longID[2], longID[3], longID[4]

		var xdcr map[string]interface{}
		err := json.Unmarshal(body, &xdcr)
		if err != nil {
			log.Error("Could not unmarshal XDCR data for remote " + uuid + " and metric " + id)
			continue
		}
		for node, values := range xdcr["nodeStats"].(map[string]interface{}) {
			list := values.([]interface{})
			if len(list) == 0 {
				log.Warn("No value found for " + id + " metric in remote " + uuid)
				continue
			}
			var metric *p.Desc
			for mid, m := range e.metrics {
				if mid == id {
					metric = m
					break
				}
			}

			var value float64
			switch v := list[len(list)-1].(type) {
			case float64:
				value = v
			case int:
				value = float64(v)
			}

			ch <- p.MustNewConstMetric(metric, p.GaugeValue, value, node, uuid, src, dest)
		}
	}
}
