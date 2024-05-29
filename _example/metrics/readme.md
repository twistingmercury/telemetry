# Observability Wrappers: Metrics Example

1. Run `go run main.go` and then 
2. Navigate to `http://localhost:8080/person`. This will cause the counters to be incremented.
3. Navigate to `http://localhost:9090/metrics` to see the metrics.
4. You should see output similar to this:

```text
# HELP example_data_call_duration Duration each func call within the data package
# TYPE example_data_call_duration histogram
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="0.1"} 0
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="0.15000000000000002"} 0
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="0.22500000000000003"} 0
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="0.3375"} 0
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="0.5062500000000001"} 0
example_data_call_duration_bucket{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data",le="+Inf"} 1
example_data_call_duration_sum{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data"} 8.3149916e+07
example_data_call_duration_count{api="example",func="DoDatabaseStuff",isError="DoDatabaseStuff",pkg="data"} 1
# HELP example_data_total_calls The total count of calls to the funcs in the data package
# TYPE example_data_total_calls counter
example_data_total_calls{api="example",func="DoDatabaseStuff",isError="false",pkg="data"} 1
# HELP example_main_call_duration The duration in milliseconds calls to the API, grouped by API name, path, and response code
# TYPE example_main_call_duration histogram
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="0.1"} 1
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="0.15000000000000002"} 1
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="0.22500000000000003"} 1
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="0.3375"} 1
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="0.5062500000000001"} 1
example_main_call_duration_bucket{http_method="GET",path="/persom",status_code="404",le="+Inf"} 1
example_main_call_duration_sum{http_method="GET",path="/persom",status_code="404"} 0
example_main_call_duration_count{http_method="GET",path="/persom",status_code="404"} 1
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="0.1"} 0
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="0.15000000000000002"} 0
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="0.22500000000000003"} 0
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="0.3375"} 0
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="0.5062500000000001"} 0
example_main_call_duration_bucket{http_method="GET",path="/person",status_code="200",le="+Inf"} 1
example_main_call_duration_sum{http_method="GET",path="/person",status_code="200"} 83
example_main_call_duration_count{http_method="GET",path="/person",status_code="200"} 1
# HELP example_main_concurrent_calls the count of concurrent calls to the APIs, grouped by API name, path, and response code
# TYPE example_main_concurrent_calls gauge
example_main_concurrent_calls{http_method="GET",path="/persom",status_code="n/a"} 0
example_main_concurrent_calls{http_method="GET",path="/person",status_code="n/a"} 0
# HELP example_main_total_calls The count of all call to the API, grouped by API name, path, and response code
# TYPE example_main_total_calls counter
example_main_total_calls{http_method="GET",path="/persom",status_code="404"} 1
example_main_total_calls{http_method="GET",path="/person",status_code="200"} 1
```