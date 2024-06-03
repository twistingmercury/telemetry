package main

import (
	"bufio"
	"example/metrics/data"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/twistingmercury/telemetry/metrics"
)

const (
	namespace   = "example"
	serviceName = "scooby"
)

func main() {
	// Metrics are hosted on in a separate goroutine on the port specified.
	// This needs to be invoked before any metrics are registered. It can be called
	// multiple times, but only the first call will have any effect.
	err := metrics.Initialize(namespace, serviceName)
	if err != nil {
		log.Fatalf("Error initializing metrics: %s\n", err.Error())
	}

	// Register the metrics from any packages that have them, in this example,
	// the data package has metrics.
	dataMetrics := data.Metrics()

	// this can be called multiple times if there are metrics in multiple packages.
	metrics.RegisterMetrics(dataMetrics...)

	// Publish exposes the metrics for scraping. This needs to be called after
	// all metrics have been registered. It can be called multiple times, but
	// only the first call will have any effect.
	metrics.Publish()

	for i := 0; i < 5; i++ {
		_ = data.DoDatabaseStuff()
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n\nusing curl to scrape metrics from port 9090:\n")
	curl := exec.Command("curl", "http://localhost:9090/metrics")
	bits, err := curl.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing curl: %s\n", err.Error())
	}
	fmt.Printf("%s\n\n", string(bits))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press any key to quit> ")
	_, _ = reader.ReadString('\n')
}
