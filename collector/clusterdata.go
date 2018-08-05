// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

// ClusterData (/pools/default)
type ClusterData struct {
	StorageTotals struct {
		RAM struct {
			Total             int `json:"total"`             // done
			QuotaTotal        int `json:"quotaTotal"`        // done
			QuotaUsed         int `json:"quotaUsed"`         // done
			Used              int `json:"used"`              // done
			UsedByData        int `json:"usedByData"`        // done
			QuotaUsedPerNode  int `json:"quotaUsedPerNode"`  // done
			QuotaTotalPerNode int `json:"quotaTotalPerNode"` // done
		} `json:"ram"`
		Hdd struct {
			Total      int64 `json:"total"`      // done
			QuotaTotal int64 `json:"quotaTotal"` // done
			Used       int64 `json:"used"`       // done
			UsedByData int   `json:"usedByData"` // done
			Free       int64 `json:"free"`       // done
		} `json:"hdd"`
	} `json:"storageTotals"`
	FtsMemoryQuota   int           `json:"ftsMemoryQuota"`   // done
	IndexMemoryQuota int           `json:"indexMemoryQuota"` // done
	MemoryQuota      int           `json:"memoryQuota"`      // done
	Name             string        `json:"name"`
	Alerts           []interface{} `json:"alerts"`
	Nodes            []struct {
		SystemStats struct {
			CPUUtilizationRate float64 `json:"cpu_utilization_rate"` // done
			SwapTotal          int     `json:"swap_total"`
			SwapUsed           int     `json:"swap_used"`
			MemTotal           int     `json:"mem_total"`
			MemFree            int     `json:"mem_free"`
		} `json:"systemStats"`
		InterestingStats struct {
			CmdGet                   int `json:"cmd_get"`
			CouchDocsActualDiskSize  int `json:"couch_docs_actual_disk_size"`
			CouchDocsDataSize        int `json:"couch_docs_data_size"`
			CouchSpatialDataSize     int `json:"couch_spatial_data_size"`
			CouchSpatialDiskSize     int `json:"couch_spatial_disk_size"`
			CouchViewsActualDiskSize int `json:"couch_views_actual_disk_size"`
			CouchViewsDataSize       int `json:"couch_views_data_size"`
			CurrItems                int `json:"curr_items"`
			CurrItemsTot             int `json:"curr_items_tot"`
			EpBgFetched              int `json:"ep_bg_fetched"`
			GetHits                  int `json:"get_hits"`
			MemUsed                  int `json:"mem_used"` // done
			Ops                      int `json:"ops"`
			VbReplicaCurrItems       int `json:"vb_replica_curr_items"`
		} `json:"interestingStats"`
		Uptime               string `json:"uptime"`
		MemoryTotal          int    `json:"memoryTotal"`
		MemoryFree           int    `json:"memoryFree"`
		McdMemoryReserved    int    `json:"mcdMemoryReserved"`
		McdMemoryAllocated   int    `json:"mcdMemoryAllocated"`
		CouchAPIBase         string `json:"couchApiBase"`
		CouchAPIBaseHTTPS    string `json:"couchApiBaseHTTPS"`
		OtpCookie            string `json:"otpCookie"`
		ClusterMembership    string `json:"clusterMembership"` // done
		RecoveryType         string `json:"recoveryType"`
		Status               string `json:"status"` // done
		OtpNode              string `json:"otpNode"`
		ThisNode             bool   `json:"thisNode"`
		Hostname             string `json:"hostname"`
		ClusterCompatibility int    `json:"clusterCompatibility"`
		Version              string `json:"version"`
		Os                   string `json:"os"`
		Ports                struct {
			SslProxy  int `json:"sslProxy"`
			HTTPSMgmt int `json:"httpsMgmt"`
			HTTPSCAPI int `json:"httpsCAPI"`
			Proxy     int `json:"proxy"`
			Direct    int `json:"direct"`
		} `json:"ports"`
		Services []string `json:"services"`
	} `json:"nodes"`
	RebalanceStatus        string `json:"rebalanceStatus"` // done
	MaxBucketCount         int    `json:"maxBucketCount"`  // done
	AutoCompactionSettings struct {
		ParallelDBAndViewCompaction    bool `json:"parallelDBAndViewCompaction"`
		DatabaseFragmentationThreshold struct {
			Percentage int    `json:"percentage"`
			Size       string `json:"size"`
		} `json:"databaseFragmentationThreshold"`
		ViewFragmentationThreshold struct {
			Percentage int    `json:"percentage"`
			Size       string `json:"size"`
		} `json:"viewFragmentationThreshold"`
		IndexCompactionMode     string `json:"indexCompactionMode"`
		IndexCircularCompaction struct {
			DaysOfWeek string `json:"daysOfWeek"`
			Interval   struct {
				FromHour     int  `json:"fromHour"`
				ToHour       int  `json:"toHour"`
				FromMinute   int  `json:"fromMinute"`
				ToMinute     int  `json:"toMinute"`
				AbortOutside bool `json:"abortOutside"`
			} `json:"interval"`
		} `json:"indexCircularCompaction"`
		IndexFragmentationThreshold struct {
			Percentage int `json:"percentage"`
		} `json:"indexFragmentationThreshold"`
	} `json:"autoCompactionSettings"`
	Counters struct {
		FailoverNode     int `json:"failover_node"`     // done
		RebalanceSuccess int `json:"rebalance_success"` // done
		RebalanceStart   int `json:"rebalance_start"`   // done
		RebalanceFail    int `json:"rebalance_fail"`    // done
	} `json:"counters"`
}
