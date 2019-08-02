// Copyright 2019 Adel Abdelhak.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE.txt file.

package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/blakelead/couchbase_exporter/collector"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cl "github.com/blakelead/confloader"
	log "github.com/sirupsen/logrus"
)

// Options struct regroups runtime parameters
type Options struct {
	serverListenAddress string
	serverMetricsPath   string
	serverTimeout       time.Duration
	dbUsername          string
	dbPassword          string
	dbURI               string
	dbTimeout           time.Duration
	tlsEnabled          bool
	tlsSkipInsecure     bool
	tlsCACert           string
	tlsClientCert       string
	tlsClientKey        string
	logLevel            string
	logFormat           string
	scrapeCluster       bool
	scrapeNode          bool
	scrapeBucket        bool
	scrapeXDCR          bool
	configFile          string
}

var (
	version        = "devel"
	runtimeOptions = &Options{}
)

func main() {

	// Initialize global variables, initialize logger and display user defined values.
	initEnv()
	initLogger()
	displayInfo()

	// Context encapsulates connection and scraping details for the exporters.
	context := collector.Context{
		URI:             runtimeOptions.dbURI,
		Username:        runtimeOptions.dbUsername,
		Password:        runtimeOptions.dbPassword,
		Timeout:         runtimeOptions.dbTimeout,
		TLSEnabled:      runtimeOptions.tlsEnabled,
		TLSSkipInsecure: runtimeOptions.tlsSkipInsecure,
		TLSCACert:       runtimeOptions.tlsCACert,
		TLSClientCert:   runtimeOptions.tlsClientCert,
		TLSClientKey:    runtimeOptions.tlsClientKey,
		ScrapeCluster:   runtimeOptions.scrapeCluster,
		ScrapeNode:      runtimeOptions.scrapeNode,
		ScrapeBucket:    runtimeOptions.scrapeBucket,
		ScrapeXDCR:      runtimeOptions.scrapeXDCR,
	}

	// Exporters are initialized, meaning that metrics files are loaded and
	// Exporter objects are created and filled with metrics metadata.
	collector.InitExporters(context)

	// Handle metrics path with the prometheus handler.
	http.Handle(runtimeOptions.serverMetricsPath, promhttp.Handler())

	// Handle paths other than given metrics path.
	if runtimeOptions.serverMetricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<body><p>See <a href="` + runtimeOptions.serverMetricsPath + `">Metrics</a></p></body>`))
		})
	}

	// Create the server with provided listen address and server timeout.
	s := &http.Server{
		Addr:        runtimeOptions.serverListenAddress,
		ReadTimeout: runtimeOptions.serverTimeout,
	}

	// Start the server.
	log.Info("Started listening at ", runtimeOptions.serverListenAddress)
	log.Fatal(s.ListenAndServe())
}

// FlagPresent returns true if the flag name passed as the
// function parameter is found in the command-line flags.
func FlagPresent(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func initEnv() {
	// Default parameters.
	runtimeOptions.serverListenAddress = "127.0.0.1:9191"
	runtimeOptions.serverMetricsPath = "/metrics"
	runtimeOptions.serverTimeout = 10 * time.Second
	runtimeOptions.dbURI = "http://localhost:8091"
	runtimeOptions.dbTimeout = 10 * time.Second
	runtimeOptions.tlsEnabled = false
	runtimeOptions.tlsSkipInsecure = false
	runtimeOptions.tlsCACert = ""
	runtimeOptions.tlsClientCert = ""
	runtimeOptions.tlsClientKey = ""
	runtimeOptions.logLevel = "info"
	runtimeOptions.logFormat = "text"
	runtimeOptions.scrapeCluster = true
	runtimeOptions.scrapeNode = true
	runtimeOptions.scrapeBucket = true
	runtimeOptions.scrapeXDCR = true
	runtimeOptions.configFile = ""

	// Get command-line values.
	var cmdlineOptions = &Options{}
	flag.StringVar(&cmdlineOptions.configFile, "config.file", runtimeOptions.configFile, "Path to configuration file.")
	flag.StringVar(&cmdlineOptions.serverListenAddress, "web.listen-address", runtimeOptions.serverListenAddress, "Address to listen on for HTTP requests.")
	flag.StringVar(&cmdlineOptions.serverMetricsPath, "web.telemetry-path", runtimeOptions.serverMetricsPath, "Path under which to expose metrics.")
	flag.DurationVar(&cmdlineOptions.serverTimeout, "web.timeout", runtimeOptions.serverTimeout, "Server read timeout in seconds.")
	flag.StringVar(&cmdlineOptions.dbURI, "db.uri", runtimeOptions.dbURI, "Couchbase node URI with port.")
	flag.DurationVar(&cmdlineOptions.dbTimeout, "db.timeout", runtimeOptions.dbTimeout, "Couchbase client timeout in seconds.")
	flag.BoolVar(&cmdlineOptions.tlsEnabled, "tls.enabled", runtimeOptions.tlsEnabled, "If true, TLS is used when communicating with cluster.")
	flag.BoolVar(&cmdlineOptions.tlsSkipInsecure, "tls.skip-insecure", runtimeOptions.tlsSkipInsecure, "If true, certificate won't be verified.")
	flag.StringVar(&cmdlineOptions.tlsCACert, "tls.ca-cert", runtimeOptions.tlsCACert, "Root certificate of the cluster.")
	flag.StringVar(&cmdlineOptions.tlsClientCert, "tls.client-cert", runtimeOptions.tlsClientCert, "Client certificate.")
	flag.StringVar(&cmdlineOptions.tlsClientKey, "tls.client-key", runtimeOptions.tlsClientKey, "Client private key.")
	flag.StringVar(&cmdlineOptions.logLevel, "log.level", runtimeOptions.logLevel, "Log level: info, debug, warn, error, fatal.")
	flag.StringVar(&cmdlineOptions.logFormat, "log.format", runtimeOptions.logFormat, "Log format: text or json.")
	flag.BoolVar(&cmdlineOptions.scrapeCluster, "scrape.cluster", runtimeOptions.scrapeCluster, "If false, cluster metrics won't be scraped.")
	flag.BoolVar(&cmdlineOptions.scrapeNode, "scrape.node", runtimeOptions.scrapeNode, "If false, node metrics won't be scraped.")
	flag.BoolVar(&cmdlineOptions.scrapeBucket, "scrape.bucket", runtimeOptions.scrapeBucket, "If false, bucket metrics won't be scraped.")
	flag.BoolVar(&cmdlineOptions.scrapeXDCR, "scrape.xdcr", runtimeOptions.scrapeXDCR, "If false, XDCR metrics won't be scraped.")
	flag.Parse()

	configFileProvided := FlagPresent("config.file")
	configLocations := [4]string{cmdlineOptions.configFile, "config.json", "config.yml", "config.yaml"}
	for idx, configLocation := range configLocations {
		config, err := cl.Load(configLocation)
		if err != nil && configFileProvided && idx == 0 {
			log.Fatalf("Could not use provided configuration file (%s). Error: %s", configLocation, err)
		}
		if err != nil {
			continue
		}
		if config.GetString("web.listenAddress") != "" {
			runtimeOptions.serverListenAddress = config.GetString("web.listenAddress")
		}
		if config.GetString("web.telemetryPath") != "" {
			runtimeOptions.serverMetricsPath = config.GetString("web.telemetryPath")
		}
		if config.GetDuration("web.timeout") != 0*time.Second {
			runtimeOptions.serverTimeout = config.GetDuration("web.timeout")
		}
		if config.GetString("db.user") != "" {
			runtimeOptions.dbUsername = config.GetString("db.user")
		}
		if config.GetString("db.password") != "" {
			runtimeOptions.dbPassword = config.GetString("db.password")
		}
		if config.GetString("db.uri") != "" {
			runtimeOptions.dbURI = config.GetString("db.uri")
		}
		if config.GetDuration("db.timeout") != 0*time.Second {
			runtimeOptions.dbTimeout = config.GetDuration("db.timeout")
		}
		if config.GetBool("tls.enabled") != runtimeOptions.tlsEnabled {
			runtimeOptions.tlsEnabled = config.GetBool("tls.enabled")
		}
		if config.GetBool("tls.skip-insecure") != runtimeOptions.tlsSkipInsecure {
			runtimeOptions.tlsSkipInsecure = config.GetBool("tls.skip-insecure")
		}
		if config.GetString("tls.ca-cert") != "" {
			runtimeOptions.tlsCACert = config.GetString("tls.ca-cert")
		}
		if config.GetString("tls.client-cert") != "" {
			runtimeOptions.tlsClientCert = config.GetString("tls.client-cert")
		}
		if config.GetString("tls.client-key") != "" {
			runtimeOptions.tlsClientKey = config.GetString("tls.client-key")
		}
		if config.GetString("log.level") != "" {
			runtimeOptions.logLevel = config.GetString("log.level")
		}
		if config.GetString("log.format") != "" {
			runtimeOptions.logFormat = config.GetString("log.format")
		}
		if config.GetBool("scrape.cluster") != runtimeOptions.scrapeCluster {
			runtimeOptions.scrapeCluster = config.GetBool("scrape.cluster")
		}
		if config.GetBool("scrape.node") != runtimeOptions.scrapeNode {
			runtimeOptions.scrapeNode = config.GetBool("scrape.node")
		}
		if config.GetBool("scrape.bucket") != runtimeOptions.scrapeBucket {
			runtimeOptions.scrapeBucket = config.GetBool("scrape.bucket")
		}
		if config.GetBool("scrape.xdcr") != runtimeOptions.scrapeXDCR {
			runtimeOptions.scrapeXDCR = config.GetBool("scrape.xdcr")
		}

		// Stop on first encounter
		runtimeOptions.configFile = configLocation
		break
	}

	// Get environment variables values.
	if val, ok := os.LookupEnv("CB_EXPORTER_LISTEN_ADDR"); ok {
		runtimeOptions.serverListenAddress = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TELEMETRY_PATH"); ok {
		runtimeOptions.serverMetricsPath = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SERVER_TIMEOUT"); ok {
		runtimeOptions.serverTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_USER"); ok {
		runtimeOptions.dbUsername = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_PASSWORD"); ok {
		runtimeOptions.dbPassword = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_URI"); ok {
		runtimeOptions.dbURI = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_TIMEOUT"); ok {
		runtimeOptions.dbTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_ENABLED"); ok {
		runtimeOptions.tlsEnabled, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_SKIP_INSECURE"); ok {
		runtimeOptions.tlsSkipInsecure, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_CA_CERT"); ok {
		runtimeOptions.tlsCACert = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_CLIENT_CERT"); ok {
		runtimeOptions.tlsClientCert = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_CLIENT_KEY"); ok {
		runtimeOptions.tlsClientKey = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_LEVEL"); ok {
		runtimeOptions.logLevel = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_FORMAT"); ok {
		runtimeOptions.logFormat = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_CLUSTER"); ok {
		runtimeOptions.scrapeCluster, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_NODE"); ok {
		runtimeOptions.scrapeNode, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_BUCKET"); ok {
		runtimeOptions.scrapeBucket, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_XDCR"); ok {
		runtimeOptions.scrapeXDCR, _ = strconv.ParseBool(val)
	}

	// Command-line values
	if FlagPresent("web.listen-address") {
		runtimeOptions.serverListenAddress = cmdlineOptions.serverListenAddress
	}
	if FlagPresent("web.telemetry-path") {
		runtimeOptions.serverMetricsPath = cmdlineOptions.serverMetricsPath
	}
	if FlagPresent("web.timeout") {
		runtimeOptions.serverTimeout = cmdlineOptions.serverTimeout
	}
	if FlagPresent("db.username") {
		runtimeOptions.dbUsername = cmdlineOptions.dbUsername
	}
	if FlagPresent("db.password") {
		runtimeOptions.dbPassword = cmdlineOptions.dbPassword
	}
	if FlagPresent("db.uri") {
		runtimeOptions.dbURI = cmdlineOptions.dbURI
	}
	if FlagPresent("db.timeout") {
		runtimeOptions.dbTimeout = cmdlineOptions.dbTimeout
	}
	if FlagPresent("tls.enabled") {
		runtimeOptions.tlsEnabled = cmdlineOptions.tlsEnabled
	}
	if FlagPresent("tls.skip-insecure") {
		runtimeOptions.tlsSkipInsecure = cmdlineOptions.tlsSkipInsecure
	}
	if FlagPresent("tls.ca-cert") {
		runtimeOptions.tlsCACert = cmdlineOptions.tlsCACert
	}
	if FlagPresent("tls.client-cert") {
		runtimeOptions.tlsClientCert = cmdlineOptions.tlsClientCert
	}
	if FlagPresent("tls.client-key") {
		runtimeOptions.tlsClientKey = cmdlineOptions.tlsClientKey
	}
	if FlagPresent("log.level") {
		runtimeOptions.logLevel = cmdlineOptions.logLevel
	}
	if FlagPresent("log.format") {
		runtimeOptions.logFormat = cmdlineOptions.logFormat
	}
	if FlagPresent("scrape.cluster") {
		runtimeOptions.scrapeCluster = cmdlineOptions.scrapeCluster
	}
	if FlagPresent("scrape.node") {
		runtimeOptions.scrapeNode = cmdlineOptions.scrapeNode
	}
	if FlagPresent("scrape.bucket") {
		runtimeOptions.scrapeBucket = cmdlineOptions.scrapeBucket
	}
	if FlagPresent("scrape.xdcr") {
		runtimeOptions.scrapeXDCR = cmdlineOptions.scrapeXDCR
	}
}

func initLogger() {
	switch runtimeOptions.logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch runtimeOptions.logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			PrettyPrint: true,
		})
	default:
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
	log.SetOutput(os.Stdout)
}

func displayInfo() {
	log.Info("Couchbase Exporter Version: ", version)
	log.Info("config.file=", runtimeOptions.configFile)
	log.Info("web.listen-address=", runtimeOptions.serverListenAddress)
	log.Info("web.telemetry-path=", runtimeOptions.serverMetricsPath)
	log.Info("web.timeout=", runtimeOptions.serverTimeout)
	log.Info("db.uri=", runtimeOptions.dbURI)
	log.Info("db.timeout=", runtimeOptions.dbTimeout)
	log.Info("tls.skip-insecure=", runtimeOptions.tlsSkipInsecure)
	log.Info("tls.ca-cert=", runtimeOptions.tlsCACert)
	log.Info("tls.enabled=", runtimeOptions.tlsEnabled)
	log.Info("tls.client-cert=", runtimeOptions.tlsClientCert)
	log.Info("tls.client-key=", runtimeOptions.tlsClientKey)
	log.Info("log.level=", runtimeOptions.logLevel)
	log.Info("log.format=", runtimeOptions.logFormat)
	log.Info("scrape.cluster=", runtimeOptions.scrapeCluster)
	log.Info("scrape.node=", runtimeOptions.scrapeNode)
	log.Info("scrape.bucket=", runtimeOptions.scrapeBucket)
	log.Info("scrape.xdcr=", runtimeOptions.scrapeXDCR)
}
