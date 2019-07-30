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

type Options struct {
	serverListenAddress string
	serverMetricsPath   string
	serverTimeout       time.Duration
	dbUsername          string
	dbPassword          string
	dbURI               string
	dbTimeout           time.Duration
	logLevel            string
	logFormat           string
	scrapeCluster       bool
	scrapeNode          bool
	scrapeBucket        bool
	scrapeXDCR          bool
	configFile          string
}

var (
	// git describe
	version                  = "0.8.0"
	runtime_options *Options = &Options{}
)

func main() {

	// Initialize global variables, initialize logger and display user defined values.
	initEnv()
	initLogger()
	displayInfo()

	// Context encapsulates connection and scraping details for the exporters.
	context := collector.Context{
		URI:           runtime_options.dbURI,
		Username:      runtime_options.dbUsername,
		Password:      runtime_options.dbPassword,
		Timeout:       runtime_options.dbTimeout,
		ScrapeCluster: runtime_options.scrapeCluster,
		ScrapeNode:    runtime_options.scrapeNode,
		ScrapeBucket:  runtime_options.scrapeBucket,
		ScrapeXDCR:    runtime_options.scrapeXDCR,
	}

	// Exporters are initialized, meaning that metrics files are loaded and
	// Exporter objects are created and filled with metrics metadata.
	collector.InitExporters(context)

	// Handle metrics path with the prometheus handler.
	http.Handle(runtime_options.serverMetricsPath, promhttp.Handler())

	// Handle paths other than given metrics path.
	if runtime_options.serverMetricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<body><p>See <a href="` + runtime_options.serverMetricsPath + `">Metrics</a></p></body>`))
		})
	}

	// Create the server with provided listen address and server timeout.
	s := &http.Server{
		Addr:        runtime_options.serverListenAddress,
		ReadTimeout: runtime_options.serverTimeout,
	}

	// Start the server.
	log.Info("Started listening at ", runtime_options.serverListenAddress)
	log.Fatal(s.ListenAndServe())
}

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
	runtime_options.serverListenAddress = "127.0.0.1:9191"
	runtime_options.serverMetricsPath = "/metrics"
	runtime_options.serverTimeout = 10 * time.Second
	runtime_options.dbURI = "http://localhost:8091"
	runtime_options.dbTimeout = 10 * time.Second
	runtime_options.logLevel = "info"
	runtime_options.logFormat = "text"
	runtime_options.scrapeCluster = true
	runtime_options.scrapeNode = true
	runtime_options.scrapeBucket = true
	runtime_options.scrapeXDCR = true
	runtime_options.configFile = ""

	// Get command-line values.
	var cmd_line_options *Options = &Options{}
	flag.StringVar(&cmd_line_options.configFile, "config.file", runtime_options.configFile, "Path to configuration file.")
	flag.StringVar(&cmd_line_options.serverListenAddress, "web.listen-address", runtime_options.serverListenAddress, "Address to listen on for HTTP requests.")
	flag.StringVar(&cmd_line_options.serverMetricsPath, "web.telemetry-path", runtime_options.serverMetricsPath, "Path under which to expose metrics.")
	flag.DurationVar(&cmd_line_options.serverTimeout, "web.timeout", runtime_options.serverTimeout, "Server read timeout in seconds.")
	flag.StringVar(&cmd_line_options.dbURI, "db.uri", runtime_options.dbURI, "Couchbase node URI with port.")
	flag.DurationVar(&cmd_line_options.dbTimeout, "db.timeout", runtime_options.dbTimeout, "Couchbase client timeout in seconds.")
	flag.StringVar(&cmd_line_options.logLevel, "log.level", runtime_options.logLevel, "Log level: info, debug, warn, error, fatal.")
	flag.StringVar(&cmd_line_options.logFormat, "log.format", runtime_options.logFormat, "Log format: text or json.")
	flag.BoolVar(&cmd_line_options.scrapeCluster, "scrape.cluster", runtime_options.scrapeCluster, "If false, cluster metrics won't be scraped.")
	flag.BoolVar(&cmd_line_options.scrapeNode, "scrape.node", runtime_options.scrapeNode, "If false, node metrics won't be scraped.")
	flag.BoolVar(&cmd_line_options.scrapeBucket, "scrape.bucket", runtime_options.scrapeBucket, "If false, bucket metrics won't be scraped.")
	flag.BoolVar(&cmd_line_options.scrapeXDCR, "scrape.xdcr", runtime_options.scrapeXDCR, "If false, XDCR metrics won't be scraped.")
	flag.Parse()

	config_file_provided := FlagPresent("config.file")
	config_locations := [3]string{cmd_line_options.configFile, "config.json", "config.yml"}
	for idx, config_location := range config_locations {
		config, err := cl.Load(config_location)
		if err != nil && config_file_provided && idx == 0 {
			log.Fatalf("Could not use provided configuration file (%s). Error: %s", config_location, err)
		}
		if err != nil {
			continue
		}
		if config.GetString("web.listenAddress") != "" {
			runtime_options.serverListenAddress = config.GetString("web.listenAddress")
		}
		if config.GetString("web.telemetryPath") != "" {
			runtime_options.serverMetricsPath = config.GetString("web.telemetryPath")
		}
		if config.GetDuration("web.timeout") != 0*time.Second {
			runtime_options.serverTimeout = config.GetDuration("web.timeout")
		}
		if config.GetString("db.user") != "" {
			runtime_options.dbUsername = config.GetString("db.user")
		}
		if config.GetString("db.password") != "" {
			runtime_options.dbPassword = config.GetString("db.password")
		}
		if config.GetString("db.uri") != "" {
			runtime_options.dbURI = config.GetString("db.uri")
		}
		if config.GetDuration("db.timeout") != 0*time.Second {
			runtime_options.dbTimeout = config.GetDuration("db.timeout")
		}
		if config.GetString("log.level") != "" {
			runtime_options.logLevel = config.GetString("log.level")
		}
		if config.GetString("log.format") != "" {
			runtime_options.logFormat = config.GetString("log.format")
		}
		if config.GetBool("scrape.cluster") != runtime_options.scrapeCluster {
			runtime_options.scrapeCluster = config.GetBool("scrape.cluster")
		}
		if config.GetBool("scrape.node") != runtime_options.scrapeNode {
			runtime_options.scrapeNode = config.GetBool("scrape.node")
		}
		if config.GetBool("scrape.bucket") != runtime_options.scrapeBucket {
			runtime_options.scrapeBucket = config.GetBool("scrape.bucket")
		}
		if config.GetBool("scrape.xdcr") != runtime_options.scrapeXDCR {
			runtime_options.scrapeXDCR = config.GetBool("scrape.xdcr")
		}

		// Stop on first encounter
		runtime_options.configFile = config_location
		break
	}

	// Get environment variables values.
	if val, ok := os.LookupEnv("CB_EXPORTER_LISTEN_ADDR"); ok {
		runtime_options.serverListenAddress = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TELEMETRY_PATH"); ok {
		runtime_options.serverMetricsPath = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SERVER_TIMEOUT"); ok {
		runtime_options.serverTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_USER"); ok {
		runtime_options.dbUsername = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_PASSWORD"); ok {
		runtime_options.dbPassword = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_URI"); ok {
		runtime_options.dbURI = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_TIMEOUT"); ok {
		runtime_options.dbTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_LEVEL"); ok {
		runtime_options.logLevel = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_FORMAT"); ok {
		runtime_options.logFormat = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_CLUSTER"); ok {
		runtime_options.scrapeCluster, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_NODE"); ok {
		runtime_options.scrapeNode, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_BUCKET"); ok {
		runtime_options.scrapeBucket, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_XDCR"); ok {
		runtime_options.scrapeXDCR, _ = strconv.ParseBool(val)
	}

	// Command-line values
	if FlagPresent("web.listen-address") {
		runtime_options.serverListenAddress = cmd_line_options.serverListenAddress
	}
	if FlagPresent("web.telemetry-path") {
		runtime_options.serverMetricsPath = cmd_line_options.serverMetricsPath
	}
	if FlagPresent("web.timeout") {
		runtime_options.serverTimeout = cmd_line_options.serverTimeout
	}
	if FlagPresent("db.username") {
		runtime_options.dbUsername = cmd_line_options.dbUsername
	}
	if FlagPresent("db.password") {
		runtime_options.dbPassword = cmd_line_options.dbPassword
	}
	if FlagPresent("db.uri") {
		runtime_options.dbURI = cmd_line_options.dbURI
	}
	if FlagPresent("db.timeout") {
		runtime_options.dbTimeout = cmd_line_options.dbTimeout
	}
	if FlagPresent("log.level") {
		runtime_options.logLevel = cmd_line_options.logLevel
	}
	if FlagPresent("log.format") {
		runtime_options.logFormat = cmd_line_options.logFormat
	}
	if FlagPresent("scrape.cluster") {
		runtime_options.scrapeCluster = cmd_line_options.scrapeCluster
	}
	if FlagPresent("scrape.node") {
		runtime_options.scrapeNode = cmd_line_options.scrapeNode
	}
	if FlagPresent("scrape.bucket") {
		runtime_options.scrapeBucket = cmd_line_options.scrapeBucket
	}
	if FlagPresent("scrape.xdcr") {
		runtime_options.scrapeXDCR = cmd_line_options.scrapeXDCR
	}
}

func initLogger() {
	switch runtime_options.logLevel {
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
	switch runtime_options.logFormat {
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
	log.Info("config.file=", runtime_options.configFile)
	log.Info("web.listen-address=", runtime_options.serverListenAddress)
	log.Info("web.telemetry-path=", runtime_options.serverMetricsPath)
	log.Info("web.timeout=", runtime_options.serverTimeout)
	log.Info("db.uri=", runtime_options.dbURI)
	log.Info("db.timeout=", runtime_options.dbTimeout)
	log.Info("log.level=", runtime_options.logLevel)
	log.Info("log.format=", runtime_options.logFormat)
	log.Info("scrape.cluster=", runtime_options.scrapeCluster)
	log.Info("scrape.node=", runtime_options.scrapeNode)
	log.Info("scrape.bucket=", runtime_options.scrapeBucket)
	log.Info("scrape.xdcr=", runtime_options.scrapeXDCR)
}
