// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Metrics is a structure that describes our metrics files
type Metrics struct {
	Name  string `json:"name"`
	Route string `json:"route"`
	List  []struct {
		Name        string   `json:"name"`
		ID          string   `json:"id"`
		Description string   `json:"description"`
		Labels      []string `json:"labels"`
	} `json:"list"`
}

// Context is a custom url wrapper with credentials
type Context struct {
	URI              string
	Username         string
	Password         string
	CouchbaseVersion string
	ConfDir          string
}

// Exporters regroups all exporters structs
type Exporters struct {
	Cluster     *ClusterExporter
	Node        *NodeExporter
	Bucket      *BucketExporter
	BucketStats *BucketStatsExporter
	XDCR        *XDCRExporter
}

// GaugeVecStruct is a wrapper for GaugeVec
type GaugeVecStruct struct {
	id       string
	gaugeVec p.GaugeVec
}

func newCounter(name string, help string) p.Counter {
	return p.NewCounter(p.CounterOpts{Namespace: "cb", Name: name, Help: help})
}

func newGaugeVec(name string, help string, labels []string) *p.GaugeVec {
	return p.NewGaugeVec(p.GaugeOpts{Namespace: "cb", Name: name, Help: help}, labels)
}

func newGaugeVecStruct(name string, id string, help string, labels []string) *GaugeVecStruct {
	return &GaugeVecStruct{
		id:       id,
		gaugeVec: *p.NewGaugeVec(p.GaugeOpts{Namespace: "cb", Name: name, Help: help}, labels),
	}
}

// NewExporters instantiates the Exporter with the URI and metrics.
func NewExporters(context Context) (*Exporters, error) {
	clusterExporter, _ := NewClusterExporter(context)
	nodeExporter, _ := NewNodeExporter(context)
	bucketExporter, _ := NewBucketExporter(context)
	bucketStatsExporter, _ := NewBucketStatsExporter(context)
	XDCRExporter, _ := NewXDCRExporter(context)

	return &Exporters{
			clusterExporter,
			nodeExporter,
			bucketExporter,
			bucketStatsExporter,
			XDCRExporter},
		nil
}

// Fetch is a helper function that fetches data from Couchbase API
func Fetch(context Context, route string) []byte {
	start := time.Now()
	req, err := http.NewRequest("GET", context.URI+route, nil)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	req.SetBasicAuth(context.Username, context.Password)
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if res.StatusCode != 200 {
		log.Error(req.Method + " " + req.URL.Path + ": " + res.Status)
		return nil
	}

	secs := time.Since(start).String()
	log.Debug("Get " + context.URI + route + " (" + secs + ")")

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return body
}

// MultiFetch is like Fetch but makes multiple requests concurrently
func MultiFetch(context Context, routes []string) map[string][]byte {
	ch := make(chan struct {
		route string
		body  []byte
	}, len(routes))
	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(route string) {
			defer wg.Done()
			ch <- struct {
				route string
				body  []byte
			}{route, Fetch(context, route)}
		}(route)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	bodies := make(map[string][]byte, 0)
	for b := range ch {
		bodies[b.route] = b.body
	}
	return bodies
}

// GetMetricsFromFile checks if metric file exist and convert it to Metrics structure
func GetMetricsFromFile(metricType string, context Context) Metrics {
	filename := context.ConfDir + "/" + metricType + "-"
	if _, err := os.Stat(filename + context.CouchbaseVersion + ".json"); os.IsNotExist(err) {
		filename = filename + "default.json"
	} else {
		filename = filename + context.CouchbaseVersion + ".json"
	}
	rawMetrics, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err.Error())
	}
	var metrics Metrics
	err = json.Unmarshal(rawMetrics, &metrics)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Loaded " + filename)
	return metrics
}

// FlattenStruct flattens structure into a map
func FlattenStruct(obj interface{}, def string) map[string]interface{} {
	fields := make(map[string]interface{}, 0)
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	for i := 0; i < objType.NumField(); i++ {
		attrField := objValue.Type().Field(i)
		valueField := objValue.Field(i)
		var key bytes.Buffer
		key.WriteString(def + attrField.Name)

		switch valueField.Kind() {
		case reflect.Struct:
			tmpMap := FlattenStruct(valueField.Interface(), attrField.Name+".")
			for k, v := range tmpMap {
				fields[k] = v
			}
		case reflect.Float64:
			fields[key.String()] = valueField.Float()
		case reflect.String:
			fields[key.String()] = valueField.String()
		case reflect.Int64:
			fields[key.String()] = valueField.Int()
		case reflect.Bool:
			fields[key.String()] = valueField.Bool()
		default:
			fields[key.String()] = valueField.Interface()
		}
	}
	return fields
}
