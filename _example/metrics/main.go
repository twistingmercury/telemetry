package main

import (
	"bufio"
	"example/metrics/data"
	"fmt"
	"os"
	"time"

	"github.com/twistingmercury/telemetry/metrics"
)

func main() {
	// Metrics are hosted on in a separate goroutine on the port specified.
	// This needs to be invoked before any metrics are registered. It can be called
	// multiple times, but only the first call will have any effect.
	metrics.InitializeWithPort("9090", "metrics", "example")

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
		time.Sleep(2 * time.Second)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press any key to quit> ")
	_, _ = reader.ReadString('\n')
}
