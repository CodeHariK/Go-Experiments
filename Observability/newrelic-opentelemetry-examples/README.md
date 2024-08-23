# Getting Started Guide - Go

This is a simple application instrumented with [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go).
It demonstrates how to configure OpenTelemetry Go to send data to New Relic.

## Requirements

* [Go](https://go.dev/dl)
* [A New Relic account](https://one.newrelic.com/)
* [A New Relic license key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#license-key)

## Running the application

1. Set the following environment variables to configure OpenTelemetry to send
   data to New Relic:

    ```shell
      export OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.nr-data.net
      export OTEL_EXPORTER_OTLP_HEADERS=api-key=5e306a91fd8db2df3d1727ee1aeca8feFFFFNRAL
      export OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT=4095
      export OTEL_SERVICE_NAME=newgo
      export OTEL_RESOURCE_ATTRIBUTES=service.instance.id=hello

      export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp.uptrace.dev"
      export OTEL_EXPORTER_OTLP_HEADERS="uptrace-dsn=https://kxlMTzzGL0dSrhwIq8ZWZQ@api.uptrace.dev?grpc=4317"
      export OTEL_EXPORTER_OTLP_COMPRESSION=gzip
      export OTEL_EXPORTER_OTLP_METRICS_DEFAULT_HISTOGRAM_AGGREGATION=BASE2_EXPONENTIAL_BUCKET_HISTOGRAM
      export OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE=DELTA

      export OTEL_SERVICE_NAME="newgo"
      export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
      export OTEL_EXPORTER_OTLP_ENDPOINT="https://api.honeycomb.io"
      export OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=3hkXPamomjdL2qtpUFss1C"
    ```

    * If your account is based in the EU, set the endpoint to: [https://otlp.eu01.nr-data.net](https://otlp.eu01.nr-data.net)

2. Run the application with the following command and open
   [http://localhost:8080/fibonacci?n=1](http://localhost:8080/fibonacci?n=1)
   in your web browser to ensure it is working.

    ```shell
    go run .
    ```

3. Experiment with providing different values for `n` in the query string.
   Valid values are between 1 and 90. Values outside this range cause an error
   which will show up in New Relic.
