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
	"path/filepath"
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
	Timeout          time.Duration
	CouchbaseVersion string
}

// Exporters regroups all exporters structs
type Exporters struct {
	Cluster     *ClusterExporter
	Node        *NodeExporter
	Bucket      *BucketExporter
	BucketStats *BucketStatsExporter
	XDCR        *XDCRExporter
}

// InitExporters instantiates the Exporter with the URI and metrics.
func InitExporters(context Context, scrapeCluster, scrapeNode, scrapeBucket, scrapeXDCR bool) {
	if scrapeCluster {
		clusterExporter, err := NewClusterExporter(context)
		if err != nil {
			log.Error("Error during creation of cluster exporter. Cluster metrics won't be scraped")
			return
		}
		p.MustRegister(clusterExporter)
		log.Debug("Cluster exporter registered")
	}
	if scrapeNode {
		nodeExporter, err := NewNodeExporter(context)
		if err != nil {
			log.Error("Error during creation of node exporter. Node metrics won't be scraped")
			return
		}
		p.MustRegister(nodeExporter)
		log.Debug("Node exporter registered")
	}
	if scrapeBucket {
		bucketExporter, err := NewBucketExporter(context)
		if err != nil {
			log.Error("Error during creation of bucket exporter. Bucket metrics won't be scraped")
			return
		}
		p.MustRegister(bucketExporter)
		log.Debug("Bucket exporter registered")
		if err != nil {
			log.Error("Error during creation of bucketstats exporter. Bucket stats metrics won't be scraped")
			return
		}
		bucketStatsExporter, _ := NewBucketStatsExporter(context)
		p.MustRegister(bucketStatsExporter)
		log.Debug("Bucketstats exporter registered")
	}
	if scrapeXDCR {
		XDCRExporter, err := NewXDCRExporter(context)
		if err != nil {
			log.Error("Error during creation of XDCR exporter. XDCR metrics won't be scraped")
			return
		}
		p.MustRegister(XDCRExporter)
		log.Debug("XDCR exporter registered")
	}
}

// Fetch is a helper function that fetches data from Couchbase API
func Fetch(context Context, route string) ([]byte, error) {
	start := time.Now()
	req, err := http.NewRequest("GET", context.URI+route, nil)
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}
	req.SetBasicAuth(context.Username, context.Password)
	client := http.Client{Timeout: context.Timeout}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		log.Error(req.Method + " " + req.URL.Path + ": " + res.Status)
		return []byte{}, err
	}

	log.Debug("Get " + context.URI + route + " (" + time.Since(start).String() + ")")

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}
	return body, nil
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
			body, err := Fetch(context, route)
			if err != nil {
				return
			}
			ch <- struct {
				route string
				body  []byte
			}{route, body}
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
func GetMetricsFromFile(metricType string, cbVersion string) (Metrics, error) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}
	exPath := filepath.Dir(ex)
	filename := exPath + "/metrics/" + metricType + "-"
	if _, err := os.Stat(filename + cbVersion + ".json"); err == nil {
		filename = filename + cbVersion + ".json"
	} else {
		filename = filename + "default.json"
	}
	rawMetrics, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Could not read file " + filename)
		return Metrics{}, err
	}
	var metrics Metrics
	err = json.Unmarshal(rawMetrics, &metrics)
	if err != nil {
		log.Error("Could not unmarshal file " + filename)
		return Metrics{}, err
	}
	log.Debug(filename + " loaded")
	return metrics, nil
}

// FlattenStruct flattens structure into a map
func FlattenStruct(obj interface{}, def ...string) map[string]interface{} {
	fields := make(map[string]interface{}, 0)
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	var prefix string
	if len(def) > 0 {
		prefix = def[0]
	}

	for i := 0; i < objType.NumField(); i++ {
		attrField := objValue.Type().Field(i)
		valueField := objValue.Field(i)
		var key bytes.Buffer
		key.WriteString(prefix + attrField.Name)

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
