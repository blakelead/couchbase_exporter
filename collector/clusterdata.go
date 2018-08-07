// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

// ClusterData (/pools/default)
type ClusterData struct {
	StorageTotals struct {
		RAM struct {
			Total      int `json:"total"`
			QuotaTotal int `json:"quotaTotal"`
			QuotaUsed  int `json:"quotaUsed"`
			Used       int `json:"used"`
			UsedByData int `json:"usedByData"`
		} `json:"ram"`
		Hdd struct {
			Total      int64 `json:"total"`
			QuotaTotal int64 `json:"quotaTotal"`
			Used       int64 `json:"used"`
			UsedByData int   `json:"usedByData"`
			Free       int64 `json:"free"`
		} `json:"hdd"`
	} `json:"storageTotals"`
	FtsMemoryQuota   int    `json:"ftsMemoryQuota"`
	IndexMemoryQuota int    `json:"indexMemoryQuota"`
	MemoryQuota      int    `json:"memoryQuota"`
	RebalanceStatus  string `json:"rebalanceStatus"`
	MaxBucketCount   int    `json:"maxBucketCount"`
	Counters         struct {
		FailoverNode     int `json:"failover_node"`
		RebalanceSuccess int `json:"rebalance_success"`
		RebalanceStart   int `json:"rebalance_start"`
		RebalanceFail    int `json:"rebalance_fail"`
	} `json:"counters"`
}
