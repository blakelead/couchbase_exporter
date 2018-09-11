// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blakelead/couchbase_exporter/collector"

	d "github.com/coreos/go-systemd/daemon"
	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddr    = flag.String("web.listen-address", ":9191", "The address to listen on for HTTP requests.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	dbURI         = flag.String("db.uri", "http://localhost:8091", "The address of Couchbase cluster.")
	dbUser        = flag.String("db.user", "admin", "The administrator username.")
	dbPwd         = flag.String("db.pwd", "password", "The administrator password.")
	logLevel      = flag.String("log.level", "info", "Log level: info, debug, warn, error, fatal.")
	logFormat     = flag.String("log.format", "text", "Log format: text or json.")
	scrapeCluster = flag.Bool("scrape.cluster", true, "If false, cluster metrics wont be scraped.")
	scrapeNode    = flag.Bool("scrape.node", true, "If false, node metrics wont be scraped.")
	scrapeBucket  = flag.Bool("scrape.bucket", true, "If false, bucket metrics wont be scraped.")
)

func main() {
	flag.Parse()

	lookupEnv()

	log.SetFormatter(setLogFormat())
	log.SetOutput(os.Stdout)
	log.SetLevel(setLogLevel())

	context := collector.Context{URI: *dbURI, Username: *dbUser, Password: *dbPwd}

	getCouchbaseVersion(&context)

	exporters, err := collector.NewExporters(context)
	if err != nil {
		log.Fatal("error during creation of new exporter")
	}

	if *scrapeCluster {
		p.MustRegister(exporters.Cluster)
	}
	if *scrapeNode {
		p.MustRegister(exporters.Node)
	}
	if *scrapeBucket {
		p.MustRegister(exporters.Bucket)
	}

	// The two following lines are used to get rid of go metrics. Should be removed after wip.
	p.Unregister(p.NewProcessCollector(os.Getpid(), ""))
	p.Unregister(p.NewGoCollector())

	// p.UninstrumentedHandler() should be replaced by promhttp.Handle() after wip.
	http.Handle(*metricsPath, p.UninstrumentedHandler())
	if *metricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Couchbase Exporter</title></head>
			<body>
			<h1>Couchbase Exporter</h1>
			<p><i>by blakelead</i></p><br>
			<p>See <a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	systemdSettings()

	// custom server used to set timeouts
	httpSrv := &http.Server{
		Addr:         *listenAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info("Listening at ", *listenAddr)
	log.Fatal(httpSrv.ListenAndServe())
}

func lookupEnv() {
	if val, ok := os.LookupEnv("CB_EXPORTER_LISTEN_ADDR"); ok {
		*listenAddr = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_TELEMETRY_PATH"); ok {
		*metricsPath = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_URI"); ok {
		*dbURI = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_USER"); ok {
		*dbUser = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_DB_PASSWORD"); ok {
		*dbPwd = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_LEVEL"); ok {
		*logLevel = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_LOG_FORMAT"); ok {
		*logFormat = val
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_CLUSTER"); ok {
		*scrapeCluster, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_NODE"); ok {
		*scrapeNode, _ = strconv.ParseBool(val)
	}
	if val, ok := os.LookupEnv("CB_EXPORTER_SCRAPE_BUCKET"); ok {
		*scrapeBucket, _ = strconv.ParseBool(val)
	}
}

func setLogLevel() log.Level {
	var level log.Level
	switch *logLevel {
	case "debug":
		level = log.DebugLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	case "fatal":
		level = log.FatalLevel
	default:
		level = log.InfoLevel
	}
	return level
}

func setLogFormat() log.Formatter {
	var format log.Formatter
	switch *logFormat {
	case "json":
		format = &log.JSONFormatter{}
	default:
		format = &log.TextFormatter{}
	}
	return format
}

func systemdSettings() {
	d.SdNotify(false, "READY=1")
	go func() {
		interval, err := d.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			_, err := http.Get(*dbURI)
			if err == nil {
				d.SdNotify(false, "WATCHDOG=1")
			}
			time.Sleep(interval / 3)
		}
	}()
}

func getCouchbaseVersion(context *collector.Context) {
	req, err := http.NewRequest("GET", *dbURI+"/pools", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(*dbUser, *dbPwd)
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	if res.StatusCode != 200 {
		log.Fatal(req.URL.Path + ": " + res.Status)
	}

	var data map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("GET " + *dbURI + "/pools" + " - data: " + string(body))

	longVersion := data["implementationVersion"].(string)
	isCommunity := strings.Contains(longVersion, "community")

	log.Info("Couchbase version: " + longVersion)
	log.Info("Community version: " + strconv.FormatBool(isCommunity))

	if !isCommunity {
		log.Warn("You are trying to scrape metrics for Couchbase Enterprise. Be aware that this exporter was not tested for Enterprise versions.")
	}

	validVersions := []string{"4.5.1", "5.1.1"}
	for _, v := range validVersions {
		if strings.HasPrefix(longVersion, v) {
			context.CouchbaseVersion = v
			return
		}
	}
	log.Warn("Please be aware that the version " + longVersion + " is not supported by this exporter.")
}
