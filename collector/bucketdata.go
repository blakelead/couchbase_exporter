// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package collector

// BucketData (/pools/default/buckets/`bucket_name`)
type BucketData struct {
	Name              string `json:"name"`
	BucketType        string `json:"bucketType"`
	AuthType          string `json:"authType"`
	SaslPassword      string `json:"saslPassword"`
	ProxyPort         int    `json:"proxyPort"`
	ReplicaIndex      bool   `json:"replicaIndex"`
	URI               string `json:"uri"`
	StreamingURI      string `json:"streamingUri"`
	LocalRandomKeyURI string `json:"localRandomKeyUri"`
	Controllers       struct {
		CompactAll    string `json:"compactAll"`
		CompactDB     string `json:"compactDB"`
		PurgeDeletes  string `json:"purgeDeletes"`
		StartRecovery string `json:"startRecovery"`
	} `json:"controllers"`
	Nodes []struct {
		CouchAPIBaseHTTPS string `json:"couchApiBaseHTTPS"`
		CouchAPIBase      string `json:"couchApiBase"`
		SystemStats       struct {
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
		Replication          int    `json:"replication"`
		ClusterMembership    string `json:"clusterMembership"`
		RecoveryType         string `json:"recoveryType"`
		Status               string `json:"status"`
		OtpNode              string `json:"otpNode"`
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
		ThisNode bool     `json:"thisNode,omitempty"`
	} `json:"nodes"`
	Stats struct {
		URI              string `json:"uri"`
		DirectoryURI     string `json:"directoryURI"`
		NodeStatsListURI string `json:"nodeStatsListURI"`
	} `json:"stats"`
	Ddocs struct {
		URI string `json:"uri"`
	} `json:"ddocs"`
	NodeLocator            string `json:"nodeLocator"`
	AutoCompactionSettings bool   `json:"autoCompactionSettings"`
	UUID                   string `json:"uuid"`
	VBucketServerMap       struct {
		HashAlgorithm string   `json:"hashAlgorithm"`
		NumReplicas   int      `json:"numReplicas"`
		ServerList    []string `json:"serverList"`
		VBucketMap    [][]int  `json:"vBucketMap"`
	} `json:"vBucketServerMap"`
	ReplicaNumber int `json:"replicaNumber"`
	ThreadsNumber int `json:"threadsNumber"`
	Quota         struct {
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
