// Copyright 2018 Adel Abdelhak.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.txt file.

package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/blakelead/couchbase_exporter/collector"

	d "github.com/coreos/go-systemd/daemon"
	p "github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddr  = flag.String("web.listen-address", ":9191", "The address to listen on for HTTP requests.")
	metricsPath = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	dbURL       = flag.String("db.url", "http://localhost:8091", "The address of Couchbase cluster.")
	dbUser      = flag.String("db.user", "admin", "The administrator username.")
	dbPwd       = flag.String("db.pwd", "password", "The administrator password.")
	logLevel    = flag.String("log.level", "info", "Log level: info, debug, warn, error, fatal.")
	logFormat   = flag.String("log.format", "text", "Log format: text or json.")
)

func main() {
	flag.Parse()

	lookupEnv()

	log.SetFormatter(setLogFormat())
	log.SetOutput(os.Stdout)
	log.SetLevel(setLogLevel())

	CouchbaseExporter, err := collector.NewExporter(collector.URI{URL: *dbURL, Username: *dbUser, Password: *dbPwd})
	if err != nil {
		log.Fatal("error during creation of new exporter")
	}
	p.MustRegister(CouchbaseExporter)

	// The two following lines are used to get rid of go metrics. Should be removed after wip.
	p.Unregister(p.NewProcessCollector(os.Getpid(), ""))
	p.Unregister(p.NewGoCollector())

	// p.UninstrumentedHandler() should be replaced by promhttp.Handle() after wip.
	http.Handle(*metricsPath, p.UninstrumentedHandler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Couchbase Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><i>by Adel Abdelhak</i></p><br>
			<p>See <a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	systemdSettings()

	log.Info("Listening at ", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func lookupEnv() {
	if val, ok := os.LookupEnv("LISTEN_ADDR"); ok {
		*listenAddr = val
	}
	if val, ok := os.LookupEnv("TELEMETRY_PATH"); ok {
		*metricsPath = val
	}
	if val, ok := os.LookupEnv("CB_URI"); ok {
		*dbURL = val
	}
	if val, ok := os.LookupEnv("CB_ADMIN_USER"); ok {
		*dbUser = val
	}
	if val, ok := os.LookupEnv("CB_ADMIN_PASSWORD"); ok {
		*dbPwd = val
	}
	if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
		*logLevel = val
	}
	if val, ok := os.LookupEnv("LOG_FORMAT"); ok {
		*logFormat = val
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
			_, err := http.Get(*dbURL)
			if err == nil {
				d.SdNotify(false, "WATCHDOG=1")
			}
			time.Sleep(interval / 3)
		}
	}()
}
