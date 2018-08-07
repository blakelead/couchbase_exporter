// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

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
