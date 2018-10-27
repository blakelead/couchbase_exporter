// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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
	metrics []*GaugeVecStruct
}

// NewXDCRExporter instantiates the Exporter with the URI and metrics.
func NewXDCRExporter(context Context) (*XDCRExporter, error) {
	xdcrMetrics := GetMetricsFromFile("xdcr", context)
	var metrics []*GaugeVecStruct
	for _, m := range xdcrMetrics.List {
		metrics = append(metrics, newGaugeVecStruct(m.Name, m.ID, m.Description, m.Labels))
	}
	return &XDCRExporter{
		context: context,
		route:   xdcrMetrics.Route,
		metrics: metrics,
	}, nil
}

// Describe describes exported metrics.
func (e *XDCRExporter) Describe(ch chan<- *p.Desc) {
	for _, m := range e.metrics {
		m.gaugeVec.Describe(ch)
	}
}

// Collect fetches data for each exported metric.
func (e *XDCRExporter) Collect(ch chan<- p.Metric) {
	var mutex sync.RWMutex
	mutex.Lock()

	// get task list where xdcr are listed
	body := Fetch(e.context, "/pools/default/tasks")
	var tasks []Task
	err := json.Unmarshal(body, &tasks)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// create urls from task list for each xdcr metric
	var routes []string
	for _, task := range tasks {
		if task.Type == "xdcr" {
			info := strings.Split(task.ID, "/") // id is in the form uuid/src/dest
			for _, m := range e.metrics {
				route := fmt.Sprintf("%s/%s/stats/replications%%2F%s%%2F%s%%2F%s%%2F%s", e.route, info[1], info[0], info[1], info[2], m.id)
				routes = append(routes, route)
			}
		}
	}

	// fetch all bodies from urls created above
	bodies := MultiFetch(e.context, routes)

	for route, body := range bodies {
		info := strings.Split(route, "%2F")
		uuid, src, dest, id := info[1], info[2], info[3], info[4]

		var xdcr map[string]interface{}
		err := json.Unmarshal(body, &xdcr)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		for node, values := range xdcr["nodeStats"].(map[string]interface{}) {
			list := values.([]interface{})
			if len(list) == 0 {
				log.Warn("No value found for " + id + " metric in remote " + uuid)
				continue
			}
			var metric *GaugeVecStruct
			for _, m := range e.metrics {
				if m.id == id {
					metric = m
					break
				}
			}
			v := list[len(list)-1].(float64)
			metric.gaugeVec.With((p.Labels{"node": node, "remote_cluster_id": uuid, "source_bucket": src, "destination_bucket": dest})).Set(v)
		}
	}

	for _, m := range e.metrics {
		m.gaugeVec.Collect(ch)
	}
	mutex.Unlock()
}
