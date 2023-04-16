package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/matsuu/go-el-controller/echonetlite"
	"github.com/matsuu/go-el-controller/wisun"
)

var version string
var bRouteID = flag.String("brouteid", "", "B-route ID")
var bRoutePW = flag.String("broutepw", "", "B-route password")
var serialPort = flag.String("serial-port", "/dev/ttyS1", "serial port for RL7023")
var exporterPort = flag.String("exporter-port", "8080", "address for prometheus")
var updateInterval = flag.Duration("interval", 1*time.Minute, "interval to get data from smart-meter")

var (
	verCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "home",
			Subsystem: "smartmeter_exporter",
			Name:      "version",
			Help:      "app version",
		},
		[]string{
			"version",
		},
	)
)

func init() {
	prometheus.MustRegister(verCounter)
}

func main() {
	flag.Parse()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	fmt.Printf("version: %s serial-port:%s exporter-port:%s\n", version, *serialPort, *exporterPort)
	verCounter.WithLabelValues(version).Inc()

	err := echonetlite.PrepareClassDictionary()
	if err != nil {
		log.Println(err)
	}

	wisunClient := wisun.NewRL7023Client(*serialPort)
	node := echonetlite.NewElectricityControllerNode(wisunClient)

	ctx := context.Background()
	initCtx, cancel := context.WithTimeout(ctx, 300*time.Second)

	err = node.Start(initCtx, *bRouteID, *bRoutePW)
	// err != nilでもコネクションを張ってるのでこの位置でClose
	defer node.Close()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to start: %w", err)
	}

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Start prometheus exporter
	ch := make(chan error)
	go func() {
		defer close(ch)
		server := http.NewServeMux()
		server.Handle("/metrics", promhttp.Handler())
		log.Println("start exporter: ", *exporterPort)
		select {
		case ch <- http.ListenAndServe(":"+*exporterPort, server):
		case <-ctx.Done():
		}
		log.Println("exporter finished")
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	func() {
		t := time.NewTicker(*updateInterval)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				_, err := node.GetPowerConsumption()
				if err != nil {
					log.Println(err)
				}
			case <-ctx.Done():
				return
			case sig := <-sigCh:
				log.Println("Signal received:", sig)
				return
			}
		}
	}()
	return nil
}
