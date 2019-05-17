// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
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
	d "github.com/coreos/go-systemd/daemon"
	log "github.com/sirupsen/logrus"
)

var (
	// git describe
	version = "0.7.1"

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
	tlsSetting          bool
)

func main() {

	initEnv()
	initLogger()
	displayInfo()

	context := collector.Context{
		URI:           dbURI,
		Username:      dbUsername,
		Password:      dbPassword,
		Timeout:       dbTimeout,
		ScrapeCluster: scrapeCluster,
		ScrapeNode:    scrapeNode,
		ScrapeBucket:  scrapeBucket,
		ScrapeXDCR:    scrapeXDCR,
		TLSSetting:    tlsSetting,
	}

	collector.InitExporters(context)

	http.Handle(serverMetricsPath, promhttp.Handler())

	if serverMetricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<body><p>See <a href="` + serverMetricsPath + `">Metrics</a></p></body>`))
		})
	}

	systemdSettings()

	s := &http.Server{
		Addr:        serverListenAddress,
		ReadTimeout: serverTimeout,
	}

	log.Info("Started listening at ", serverListenAddress)
	log.Fatal(s.ListenAndServe())
}

func initEnv() {
	// default parameters
	serverListenAddress = "127.0.0.1:9191"
	serverMetricsPath = "/metrics"
	serverTimeout = 10 * time.Second
	dbURI = "https://127.0.0.1:18091"
	dbTimeout = 10 * time.Second
	logLevel = "error"
	logFormat = "text"
	scrapeCluster = true
	scrapeNode = true
	scrapeBucket = true
	scrapeXDCR = true
	tlsSetting = false

	// get configuration file values
	configFile = "config.json"
	config, err := cl.Load(configFile)
	if err != nil {
		configFile = "config.yml"
		config, err = cl.Load(configFile)
		if err != nil {
			log.Info(err, ": using command-line parameters and/or environment variables if provided")
		}
	}
	if err == nil {
		serverListenAddress = config.GetString("web.listenAddress")
		serverMetricsPath = config.GetString("web.telemetryPath")
		serverTimeout = config.GetDuration("web.timeout")
		dbUsername = config.GetString("db.user")
		dbPassword = config.GetString("db.password")
		dbURI = config.GetString("db.uri")
		dbTimeout = config.GetDuration("db.timeout")
		logLevel = config.GetString("log.level")
		logFormat = config.GetString("log.format")
		scrapeCluster = config.GetBool("scrape.cluster")
		scrapeNode = config.GetBool("scrape.node")
		scrapeBucket = config.GetBool("scrape.bucket")
		scrapeXDCR = config.GetBool("scrape.xdcr")
		tlsSetting = config.GetBool("tls.setting")
	}

	// get environment variables values
	if val, ok := os.LookupEnv("CB_EXPORTER_LISTEN_ADDR"); ok {
		serverListenAddress = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TELEMETRY_PATH"); ok {
		serverMetricsPath = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SERVER_TIMEOUT"); ok {
		serverTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_USER"); ok {
		dbUsername = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_PASSWORD"); ok {
		dbPassword = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_URI"); ok {
		dbURI = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_TIMEOUT"); ok {
		dbTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_LEVEL"); ok {
		logLevel = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_FORMAT"); ok {
		logFormat = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_CLUSTER"); ok {
		scrapeCluster, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_NODE"); ok {
		scrapeNode, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_BUCKET"); ok {
		scrapeBucket, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_XDCR"); ok {
		scrapeXDCR, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TLS_SETTING"); ok {
		tlsSetting, _ = strconv.ParseBool(val)
	}

	// get command-line values
	flag.StringVar(&serverListenAddress, "web.listen-address", serverListenAddress, "Address to listen on for HTTP requests.")
	flag.StringVar(&serverMetricsPath, "web.telemetry-path", serverMetricsPath, "Path under which to expose metrics.")
	flag.DurationVar(&serverTimeout, "web.timeout", serverTimeout, "Server read timeout in seconds.")
	flag.StringVar(&dbURI, "db.uri", dbURI, "Couchbase node URI with port.")
	flag.DurationVar(&dbTimeout, "db.timeout", dbTimeout, "Couchbase client timeout in seconds.")
	flag.StringVar(&logLevel, "log.level", logLevel, "Log level: info, debug, warn, error, fatal.")
	flag.StringVar(&logFormat, "log.format", logFormat, "Log format: text or json.")
	flag.BoolVar(&scrapeCluster, "scrape.cluster", scrapeCluster, "If false, cluster metrics won't be scraped.")
	flag.BoolVar(&scrapeNode, "scrape.node", scrapeNode, "If false, node metrics won't be scraped.")
	flag.BoolVar(&scrapeBucket, "scrape.bucket", scrapeBucket, "If false, bucket metrics won't be scraped.")
	flag.BoolVar(&scrapeXDCR, "scrape.xdcr", scrapeXDCR, "If false, XDCR metrics won't be scraped.")
	flag.BoolVar(&tlsSetting, "tls.setting", tlsSetting, "If true TLS will ignore self signed certificates.")
	flag.Parse()
}

func initLogger() {
	switch logLevel {
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
	switch logFormat {
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
	log.Info("config.file=", configFile)
	log.Info("web.listen-address=", serverListenAddress)
	log.Info("web.telemetry-path=", serverMetricsPath)
	log.Info("web.timeout=", serverTimeout)
	log.Info("db.uri=", dbURI)
	log.Info("db.timeout=", dbTimeout)
	log.Info("log.level=", logLevel)
	log.Info("log.format=", logFormat)
	log.Info("scrape.cluster=", scrapeCluster)
	log.Info("scrape.node=", scrapeNode)
	log.Info("scrape.bucket=", scrapeBucket)
	log.Info("scrape.xdcr=", scrapeXDCR)
	log.Info("tls.setting=", tlsSetting)
}

func systemdSettings() {
	d.SdNotify(false, "READY=1")
	go func() {
		interval, err := d.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			_, err := http.Head(dbURI)
			if err == nil {
				d.SdNotify(false, "WATCHDOG=1")
			}
			time.Sleep(interval / 3)
		}
	}()
}
