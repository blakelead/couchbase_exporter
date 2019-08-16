// Copyright 2019 Adel Abdelhak.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE.txt file.

package collector

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
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

// Metrics is a structure that describes the metrics files
// that hold  all metrics  informations used  for scraping
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

// Context is a custom url wrapper with credentials and
// booleans about which metrics types should be scraped
type Context struct {
	URI             string
	Username        string
	Password        string
	Timeout         time.Duration
	ScrapeCluster   bool
	ScrapeNode      bool
	ScrapeBucket    bool
	ScrapeXDCR      bool
	TLSEnabled      bool
	TLSSkipInsecure bool
	TLSCACert       string
	TLSClientCert   string
	TLSClientKey    string
}

// Exporters structure contains all exporters
type Exporters struct {
	Cluster     *ClusterExporter
	Node        *NodeExporter
	Bucket      *BucketExporter
	BucketStats *BucketStatsExporter
	XDCR        *XDCRExporter
}

// InitExporters instantiates the Exporters
func InitExporters(c Context) {
	if c.ScrapeCluster {
		clusterExporter, err := NewClusterExporter(c)
		if err != nil {
			log.Error("Error during creation of cluster exporter. Cluster metrics won't be scraped")
		} else {
			p.MustRegister(clusterExporter)
			log.Info("Cluster exporter registered")
		}
	}
	if c.ScrapeNode {
		nodeExporter, err := NewNodeExporter(c)
		if err != nil {
			log.Error("Error during creation of node exporter. Node metrics won't be scraped")
		} else {
			p.MustRegister(nodeExporter)
			log.Info("Node exporter registered")
		}
	}
	if c.ScrapeBucket {
		bucketExporter, err := NewBucketExporter(c)
		if err != nil {
			log.Error("Error during creation of bucket exporter. Bucket metrics won't be scraped")
		} else {
			p.MustRegister(bucketExporter)
			log.Info("Bucket exporter registered")
		}
		bucketStatsExporter, err := NewBucketStatsExporter(c)
		if err != nil {
			log.Error("Error during creation of bucketstats exporter. Bucket stats metrics won't be scraped")
		} else {
			p.MustRegister(bucketStatsExporter)
			log.Info("Bucketstats exporter registered")
		}
	}
	if c.ScrapeXDCR {
		xdcrExporter, err := NewXDCRExporter(c)
		if err != nil {
			log.Error("Error during creation of XDCR exporter. XDCR metrics won't be scraped")
		} else {
			p.MustRegister(xdcrExporter)
			log.Info("XDCR exporter registered")
		}
	}
}

// Fetch is a helper function that fetches data from Couchbase API
func Fetch(c Context, route string) ([]byte, error) {
	start := time.Now()

	req, err := http.NewRequest("GET", c.URI+route, nil)
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}

	tlsClientConfig := &tls.Config{}
	if c.TLSEnabled {
		tlsClientConfig, err = createTLSClientConfig(c)
		if err != nil {
			log.Error(err)
		}
	}

	req.SetBasicAuth(c.Username, c.Password)
	client := http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsClientConfig,
		},
	}

	res, err := client.Do(req)

	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error(req.Method + " " + req.URL.Path + ": " + res.Status)
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error())
		return []byte{}, err
	}

	log.Debug("Get " + c.URI + route + " (" + time.Since(start).String() + ")")

	return body, nil
}

// MultiFetch is like Fetch but makes multiple requests concurrently
func MultiFetch(c Context, routes []string) map[string][]byte {
	ch := make(chan struct {
		route string
		body  []byte
	}, len(routes))

	var wg sync.WaitGroup
	for _, route := range routes {
		wg.Add(1)
		go func(route string) {
			defer wg.Done()
			body, err := Fetch(c, route)
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

	bodies := make(map[string][]byte, len(ch))
	for b := range ch {
		bodies[b.route] = b.body
	}

	return bodies
}

// GetMetricsFromFile checks if metric file exist and convert it to Metrics structure
func GetMetricsFromFile(metricType string) (Metrics, error) {
	absPath, err := os.Executable()
	if err != nil {
		log.Error("An unknown error occurred: ", err)
		return Metrics{}, err
	}

	filename := filepath.Dir(absPath) + string(os.PathSeparator) + "metrics" + string(os.PathSeparator) + metricType + ".json"
	if _, err := os.Stat(filename); err != nil {
		log.Error("Could not find metrics file ", filename)
		return Metrics{}, err
	}

	rawMetrics, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Could not read file ", filename)
		return Metrics{}, err
	}

	var metrics Metrics
	err = json.Unmarshal(rawMetrics, &metrics)
	if err != nil {
		log.Error("Could not unmarshal file ", filename)
		return Metrics{}, err
	}

	log.Debug(filename, " loaded")

	return metrics, nil
}

// FlattenStruct flattens structure into a Go map
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

// createTLSClientConfig loads certificates and create TLS config
func createTLSClientConfig(c Context) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(c.TLSCACert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	keyPair, err := tls.LoadX509KeyPair(c.TLSClientCert, c.TLSClientKey)
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		Certificates:       []tls.Certificate{keyPair},
		ClientCAs:          certPool,
		InsecureSkipVerify: c.TLSSkipInsecure,
	}

	return &config, nil
}
