// Copyright 2018 Blake Lead.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

// ClusterData (/pools/default)
type ClusterData struct {
	StorageTotals struct {
		RAM struct {
			Total             int `json:"total"`
			QuotaTotal        int `json:"quotaTotal"`
			QuotaUsed         int `json:"quotaUsed"`
			Used              int `json:"used"`
			UsedByData        int `json:"usedByData"`
			QuotaUsedPerNode  int `json:"quotaUsedPerNode"`
			QuotaTotalPerNode int `json:"quotaTotalPerNode"`
		} `json:"ram"`
		Hdd struct {
			Total      int64 `json:"total"`
			QuotaTotal int64 `json:"quotaTotal"`
			Used       int64 `json:"used"`
			UsedByData int   `json:"usedByData"`
			Free       int64 `json:"free"`
		} `json:"hdd"`
	} `json:"storageTotals"`
	FtsMemoryQuota   int           `json:"ftsMemoryQuota"`
	IndexMemoryQuota int           `json:"indexMemoryQuota"`
	MemoryQuota      int           `json:"memoryQuota"`
	Name             string        `json:"name"`
	Alerts           []interface{} `json:"alerts"`
	AlertsSilenceURL string        `json:"alertsSilenceURL"`
	Nodes            []struct {
		SystemStats struct {
			CPUUtilizationRate float64 `json:"cpu_utilization_rate"`
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
			MemUsed                  int `json:"mem_used"`
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
		ClusterMembership    string `json:"clusterMembership"`
		RecoveryType         string `json:"recoveryType"`
		Status               string `json:"status"`
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
	Buckets struct {
		URI                       string `json:"uri"`
		TerseBucketsBase          string `json:"terseBucketsBase"`
		TerseStreamingBucketsBase string `json:"terseStreamingBucketsBase"`
	} `json:"buckets"`
	RemoteClusters struct {
		URI         string `json:"uri"`
		ValidateURI string `json:"validateURI"`
	} `json:"remoteClusters"`
	Controllers struct {
		AddNode struct {
			URI string `json:"uri"`
		} `json:"addNode"`
		Rebalance struct {
			URI string `json:"uri"`
		} `json:"rebalance"`
		FailOver struct {
			URI string `json:"uri"`
		} `json:"failOver"`
		StartGracefulFailover struct {
			URI string `json:"uri"`
		} `json:"startGracefulFailover"`
		ReAddNode struct {
			URI string `json:"uri"`
		} `json:"reAddNode"`
		ReFailOver struct {
			URI string `json:"uri"`
		} `json:"reFailOver"`
		EjectNode struct {
			URI string `json:"uri"`
		} `json:"ejectNode"`
		SetRecoveryType struct {
			URI string `json:"uri"`
		} `json:"setRecoveryType"`
		SetAutoCompaction struct {
			URI         string `json:"uri"`
			ValidateURI string `json:"validateURI"`
		} `json:"setAutoCompaction"`
		ClusterLogsCollection struct {
			StartURI  string `json:"startURI"`
			CancelURI string `json:"cancelURI"`
		} `json:"clusterLogsCollection"`
		Replication struct {
			CreateURI   string `json:"createURI"`
			ValidateURI string `json:"validateURI"`
		} `json:"replication"`
	} `json:"controllers"`
	RebalanceStatus        string `json:"rebalanceStatus"`
	RebalanceProgressURI   string `json:"rebalanceProgressUri"`
	StopRebalanceURI       string `json:"stopRebalanceUri"`
	NodeStatusesURI        string `json:"nodeStatusesUri"`
	MaxBucketCount         int    `json:"maxBucketCount"`
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
	Tasks struct {
		URI string `json:"uri"`
	} `json:"tasks"`
	Counters struct {
	} `json:"counters"`
	IndexStatusURI      string `json:"indexStatusURI"`
	CheckPermissionsURI string `json:"checkPermissionsURI"`
	ServerGroupsURI     string `json:"serverGroupsUri"`
}
