// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// BucketData (/pools/default/buckets)
type BucketData struct {
	Name                   string `json:"name"`
	BucketType             string `json:"bucketType"`
	AuthType               string `json:"authType"`
	ProxyPort              int    `json:"proxyPort"`
	ReplicaIndex           bool   `json:"replicaIndex"`
	NodeLocator            string `json:"nodeLocator"`
	AutoCompactionSettings bool   `json:"autoCompactionSettings"`
	UUID                   string `json:"uuid"`
	ReplicaNumber          int    `json:"replicaNumber"`
	ThreadsNumber          int    `json:"threadsNumber"`
	Quota                  struct {
		RAM    int `json:"ram"`
		RawRAM int `json:"rawRAM"`
	} `json:"quota"`
	BasicStats struct {
		QuotaPercentUsed float64 `json:"quotaPercentUsed"`
		OpsPerSec        int     `json:"opsPerSec"`
		DiskFetches      int     `json:"diskFetches"`
		ItemCount        int     `json:"itemCount"`
		DiskUsed         int     `json:"diskUsed"`
		DataUsed         int     `json:"dataUsed"`
		MemUsed          int     `json:"memUsed"`
	} `json:"basicStats"`
	EvictionPolicy        string   `json:"evictionPolicy"`
	TimeSynchronization   string   `json:"timeSynchronization"`
	BucketCapabilitiesVer string   `json:"bucketCapabilitiesVer"`
	BucketCapabilities    []string `json:"bucketCapabilities"`
}

// BucketExporter describes the exporter object.
type BucketExporter struct {
	uri                    URI
	bucketProxyPort        *p.GaugeVec
	bucketReplicaIndex     *p.GaugeVec
	bucketReplicaNumber    *p.GaugeVec
	bucketThreadsNumber    *p.GaugeVec
	bucketRAMQuota         *p.GaugeVec
	bucketRawRAMQuota      *p.GaugeVec
	bucketQuotaPercentUsed *p.GaugeVec
	bucketOpsPerSec        *p.GaugeVec
	bucketDiskFetches      *p.GaugeVec
	bucketItemCount        *p.GaugeVec
	bucketDiskUsed         *p.GaugeVec
	bucketDataUsed         *p.GaugeVec
	bucketMemUsed          *p.GaugeVec
}

// NewBucketExporter instantiates the Exporter with the URI and metrics.
func NewBucketExporter(uri URI) (*BucketExporter, error) {
	return &BucketExporter{
		uri:                    uri,
		bucketProxyPort:        newGaugeVec("bucket_proxy_port", "Bucket proxy port", []string{"bucket"}),
		bucketReplicaIndex:     newGaugeVec("bucket_replica_index", "Replica index for the bucket. 1:true", []string{"bucket"}),
		bucketReplicaNumber:    newGaugeVec("bucket_replica_number", "Number of replicas for the bucket", []string{"bucket"}),
		bucketThreadsNumber:    newGaugeVec("bucket_threads_number", "Bucket thread number", []string{"bucket"}),
		bucketRAMQuota:         newGaugeVec("bucket_ram_quota_bytes", "Memory used by the bucket", []string{"bucket"}),
		bucketRawRAMQuota:      newGaugeVec("bucket_raw_ram_quota_bytes", "Raw memory used by the bucket", []string{"bucket"}),
		bucketQuotaPercentUsed: newGaugeVec("bucket_ram_quota_percent_used", "Memory used by the bucket in percent", []string{"bucket"}),
		bucketOpsPerSec:        newGaugeVec("bucket_ops_per_second", "Number of operations per second in the bucket", []string{"bucket"}),
		bucketDiskFetches:      newGaugeVec("bucket_disk_fetches", "Disk fetches for the bucket", []string{"bucket"}),
		bucketItemCount:        newGaugeVec("bucket_item_count", "Number of items in the bucket", []string{"bucket"}),
		bucketDiskUsed:         newGaugeVec("bucket_disk_used_bytes", "Disk used by the bucket", []string{"bucket"}),
		bucketDataUsed:         newGaugeVec("bucket_data_used_bytes", "Data loaded in memory", []string{"bucket"}),
		bucketMemUsed:          newGaugeVec("bucket_ram_used_bytes", "Bucket RAM used", []string{"bucket"}),
	}, nil
}

// Describe describes exported metrics.
func (e *BucketExporter) Describe(ch chan<- *p.Desc) {
	e.bucketProxyPort.Describe(ch)
	e.bucketReplicaIndex.Describe(ch)
	e.bucketReplicaNumber.Describe(ch)
	e.bucketThreadsNumber.Describe(ch)
	e.bucketRAMQuota.Describe(ch)
	e.bucketRawRAMQuota.Describe(ch)
	e.bucketQuotaPercentUsed.Describe(ch)
	e.bucketOpsPerSec.Describe(ch)
	e.bucketDiskFetches.Describe(ch)
	e.bucketItemCount.Describe(ch)
	e.bucketDiskUsed.Describe(ch)
	e.bucketDataUsed.Describe(ch)
	e.bucketMemUsed.Describe(ch)
}

// Collect fetches data for each exported metric.
func (e *BucketExporter) Collect(ch chan<- p.Metric) {
	e.scrape()
	e.bucketProxyPort.Collect(ch)
	e.bucketReplicaIndex.Collect(ch)
	e.bucketReplicaNumber.Collect(ch)
	e.bucketThreadsNumber.Collect(ch)
	e.bucketRAMQuota.Collect(ch)
	e.bucketRawRAMQuota.Collect(ch)
	e.bucketQuotaPercentUsed.Collect(ch)
	e.bucketOpsPerSec.Collect(ch)
	e.bucketDiskFetches.Collect(ch)
	e.bucketItemCount.Collect(ch)
	e.bucketDiskUsed.Collect(ch)
	e.bucketDataUsed.Collect(ch)
	e.bucketMemUsed.Collect(ch)
}

func (e *BucketExporter) scrape() {
	req, err := http.NewRequest("GET", e.uri.URL+"/pools/default/buckets", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.SetBasicAuth(e.uri.Username, e.uri.Password)
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if res.StatusCode != 200 {
		log.Error(req.URL.Path + ": " + res.Status)
		return
	}

	var buckets []BucketData
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}
	err = json.Unmarshal([]byte(body), &buckets)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("GET " + e.uri.URL + "/pools/default/buckets" + " - data: " + string(body))

	for _, bucket := range buckets {
		var replicaIndex int
		if bucket.ReplicaIndex == true {
			replicaIndex = 1
		}
		e.bucketProxyPort.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ProxyPort))
		e.bucketReplicaIndex.With(p.Labels{"bucket": bucket.Name}).Set(float64(replicaIndex))
		e.bucketReplicaNumber.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ReplicaNumber))
		e.bucketThreadsNumber.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.ThreadsNumber))
		e.bucketRAMQuota.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.Quota.RAM))
		e.bucketRawRAMQuota.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.Quota.RawRAM))
		e.bucketQuotaPercentUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.QuotaPercentUsed))
		e.bucketOpsPerSec.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.OpsPerSec))
		e.bucketDiskFetches.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskFetches))
		e.bucketItemCount.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.ItemCount))
		e.bucketDiskUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskUsed))
		e.bucketDataUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.DiskUsed))
		e.bucketMemUsed.With(p.Labels{"bucket": bucket.Name}).Set(float64(bucket.BasicStats.MemUsed))

	}

}
