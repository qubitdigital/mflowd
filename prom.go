package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type promMessage struct {
	gatherer prom.Gatherer
}

type promIO struct {
	scrapeSignalChan chan bool
	messageChan      chan promMessage
}

func mflowPromHandler(pio *promIO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Prometheus wants to scrape some metrics ...")
		pio.scrapeSignalChan <- true
		pmsg := <-pio.messageChan

		promhttp.HandlerFor(pmsg.gatherer, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}

func exposePrometheusEndpoint(endpoint string, port int, pio *promIO) {
	log.Debugf("Exposing a Prometheus endpoint at %d", port)
	http.HandleFunc(endpoint, mflowPromHandler(pio))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
