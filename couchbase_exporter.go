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
)

func main() {
	flag.Parse()

	lookupEnv()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

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

	// Systemd settings
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
}
