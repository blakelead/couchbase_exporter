// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	cl "github.com/blakelead/confloader"
	"github.com/blakelead/couchbase_exporter/collector"

	d "github.com/coreos/go-systemd/daemon"
	log "github.com/sirupsen/logrus"
)

var (
	// git describe
	version = "0.5.1-8-g795d467"

	// supported Couchbase versions
	validVersions = []string{"4.5.1", "4.6.5", "5.1.1"}

	configFile       string
	webListenAddr    string
	webMetricsPath   string
	webServerTimeout time.Duration
	dbUser           string
	dbPwd            string
	dbURI            string
	dbTimeout        time.Duration
	logLevel         string
	logFormat        string
	scrapeCluster    bool
	scrapeNode       bool
	scrapeBucket     bool
	scrapeXDCR       bool
)

func main() {

	initParams()
	initLogger()

	log.Info("Couchbase Exporter Version ", version)
	log.Info("config.file=", configFile)
	log.Info("web.listen-address=", webListenAddr)
	log.Info("web.telemetry-path=", webMetricsPath)
	log.Info("web.timeout=", webServerTimeout)
	log.Info("db.uri=", dbURI)
	log.Info("db.timeout=", dbTimeout)
	log.Info("log.level=", logLevel)
	log.Info("log.format=", logFormat)
	log.Info("scrape.cluster=", scrapeCluster)
	log.Info("scrape.node=", scrapeNode)
	log.Info("scrape.bucket=", scrapeBucket)
	log.Info("scrape.xdcr=", scrapeXDCR)

	context := collector.Context{URI: dbURI, Username: dbUser, Password: dbPwd, Timeout: dbTimeout}

	getCouchbaseVersion(&context)

	collector.InitExporters(context, scrapeCluster, scrapeNode, scrapeBucket, scrapeXDCR)

	http.Handle(webMetricsPath, promhttp.Handler())
	if webMetricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`
			<html>
			<head><title>Couchbase Exporter</title></head>
			<body>
			<h1>Couchbase Exporter</h1>
			<p><i>by blakelead</i></p><br>
			<p>See <a href="` + webMetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	systemdSettings()

	s := &http.Server{Addr: webListenAddr, ReadTimeout: webServerTimeout}
	log.Info("Started listening at ", webListenAddr)
	log.Fatal(s.ListenAndServe())
}

func initParams() {
	// default parameters
	configFile = "config.json"
	webListenAddr = "127.0.0.1:9191"
	webMetricsPath = "/metrics"
	webServerTimeout = 10 * time.Second
	dbURI = "http://localhost:8091"
	dbTimeout = 10 * time.Second
	logLevel = "info"
	logFormat = "text"
	scrapeCluster = true
	scrapeNode = true
	scrapeBucket = true
	scrapeXDCR = false

	// get configuration file values
	config, err := cl.Load(configFile)
	if err != nil {
		configFile = "config.yml"
		config, err = cl.Load(configFile)
		if err != nil {
			log.Info(err, ": using command-line parameters and/or environment variables if provided")
		}
	}
	if err == nil {
		webListenAddr = config.GetString("web.listenAddress")
		webMetricsPath = config.GetString("web.telemetryPath")
		webServerTimeout = config.GetDuration("web.timeout")
		dbUser = config.GetString("db.user")
		dbPwd = config.GetString("db.password")
		dbURI = config.GetString("db.uri")
		dbTimeout = config.GetDuration("db.timeout")
		logLevel = config.GetString("log.level")
		logFormat = config.GetString("log.format")
		scrapeCluster = config.GetBool("scrape.cluster")
		scrapeNode = config.GetBool("scrape.node")
		scrapeBucket = config.GetBool("scrape.bucket")
		scrapeXDCR = config.GetBool("scrape.xdcr")
	}

	// get environment variables values
	if val, ok := os.LookupEnv("CB_EXPORTER_LISTEN_ADDR"); ok {
		webListenAddr = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TELEMETRY_PATH"); ok {
		webMetricsPath = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SERVER_TIMEOUT"); ok {
		webServerTimeout, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_USER"); ok {
		dbUser = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_PASSWORD"); ok {
		dbPwd = val
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

	// get command-line values
	flag.StringVar(&webListenAddr, "web.listen-address", webListenAddr, "The address to listen on for HTTP requests.")
	flag.StringVar(&webMetricsPath, "web.telemetry-path", webMetricsPath, "Path under which to expose metrics.")
	flag.DurationVar(&webServerTimeout, "web.timeout", webServerTimeout, "Server read timeout in seconds.")
	flag.StringVar(&dbURI, "db.uri", dbURI, "The address of Couchbase cluster.")
	flag.DurationVar(&dbTimeout, "db.timeout", dbTimeout, "Couchbase client timeout in seconds.")
	flag.StringVar(&logLevel, "log.level", logLevel, "Log level: info, debug, warn, error, fatal.")
	flag.StringVar(&logFormat, "log.format", logFormat, "Log format: text or json.")
	flag.BoolVar(&scrapeCluster, "scrape.cluster", scrapeCluster, "If false, cluster metrics wont be scraped.")
	flag.BoolVar(&scrapeNode, "scrape.node", scrapeNode, "If false, node metrics wont be scraped.")
	flag.BoolVar(&scrapeBucket, "scrape.bucket", scrapeBucket, "If false, bucket metrics wont be scraped.")
	flag.BoolVar(&scrapeXDCR, "scrape.xdcr", scrapeXDCR, "If false, XDCR metrics wont be scraped.")
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

func getCouchbaseVersion(context *collector.Context) {
	body, err := collector.Fetch(*context, "/pools")
	if err != nil {
		log.Error("Error when retrieving Couchbase version informations")
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("Error when retrieving Couchbase version informations")
		return
	}

	rawVersion := data["implementationVersion"].(string)

	log.Info("Couchbase version ", rawVersion)

	for _, v := range validVersions {
		if strings.HasPrefix(rawVersion, v) {
			context.CouchbaseVersion = v
			return
		}
	}

	log.Warn("Version ", rawVersion, " may not be supported by this exporter")
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
