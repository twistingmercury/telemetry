#  Package: Metrics

This repository contains middleware for [gin and gonic](https://github.com/gin-gonic/gin) using [github.com/prometheus/client_golang]( https://pkg.go.dev/github.com/prometheus/client_golang/prometheus).

## Installation

```bash
go get -u github.com/twistingmercury/telemetry
```

## Initialization

This is the general process for initializing the metrics:

1. InitializeWithPort the metrics with the `metrics.InitializeWithPort` function. This function must be called first. This function takes three parameters:
    * The port to expose the metrics on (a valid port is a number between 1024 and 49151)
    * The namespace used to help identify the metrics
    * The name of the service, api, etc., the metrics will be associated with.

2. Register any custom metrics with the `metrics.RegisterMetrics` function. This function takes one or more
   `prometheus.Collector` instances. Creating a `prometheus.Collector` is beyond the scope of this document. See
   the [prometheus documentation](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus@v1.17.0#pkg-types)
   for more information.

3. Publish the metrics with the `metrics.Publish` function. This function takes no parameters.

## Usage

### Instrumenting packages

You can instrument any function by creating one or more `prometheus.Collector` types and registering it with the 
`metrics.RegisterMetrics` function. The below code sample demonstrates basic usage:

1. Add a func to the package in which you will define all the metrics for that package:

    ```go  
    package 
    
    // Metrics returns the metrics that are defined for the data package.
    func Metrics() []prometheus.Collector {
        labels := []string{ "pkg", "func", "is_error"}
    
        totalCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: metrics.Namespace(),     // You can use the namespace set during initialization, or use a different one.
            Name:      "<api_name>_<package>_total_calls", // I typically use package name as a prefix.
            Help:      "The total count of calls to the funcs in the package"},
            labels)
    
        callDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
            Namespace: metrics.Namespace(),
            Name:      "<api_name>_<package>_data_call_duration",
            Help:      "Duration each func call within the package",
            Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
        }, labels)
    
        return []prometheus.Collector{totalCalls, callDuration}
    }
    ```

2. To help with incrementing the counters, I like to add a helper func as well:

   ```go
    func incMetrics(fName string, d float64, err error) {
        tCtr.WithLabelValues("examples", "data", fName, strconv.FormatBool(err != nil)).Inc()
        dHist.WithLabelValues("examples", "data", fName, "DoDatabaseStuff").Observe(d)
    }
   ```
3. Then invoke the `incMetrics` func as appropriate:

   ```go
    func DoStuff() (err error) {
    s := time.Now()
    defer func() {
        duration := float64(time.Since(s))
        incMetrics("DoStuff", duration, err)
    }()
    
        src := rand.NewSource(time.Now().UnixNano())
        rnd := rand.New(src)
    
        minSleep := 10
        maxSleep := 100
    
        // simulate some random latency
        time.Sleep(time.Duration(rnd.Intn(maxSleep-minSleep)+minSleep) * time.Millisecond)
    
        // simulate a random error...
        if rnd.Intn(24)%7 == 0 {
            err = fmt.Errorf("random simulated error")
            return
        }
    
        return
    }
   ```

4. Finally, register the custom metrics with the `metrics.RegisterMetrics` function before you make the call to publish:

    ```go
    func main(){
        // initialize other stuff, like logging, configuration, etc., ...
   
        metrics.Initialize("my-namespace", "my-service")
        customMetrics := somePkg.Metrics()
        metrics.RegisterMetrics(customMetrics...)
        metrics.Publish()
   
        // start whatever the service should be doing...
   }
    ```