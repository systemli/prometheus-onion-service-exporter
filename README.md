# prometheus-onion-exporter

[![Integration](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/integration.yaml/badge.svg)](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/integration.yaml)
[![Quality](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/quality.yaml/badge.svg)](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/quality.yaml)
[![Release](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/systemli/prometheus-onion-service-exporter/actions/workflows/release.yaml)

Prometheus Exporter for Tor onion services written in Go.
Export the status and latency of an onion service to prometheus.

## Usage

```
go get github.com/systemli/prometheus-onion-service-exporter
go install github.com/systemli/prometheus-onion-service-exporter
$GOPATH/bin/prometheus-onion-service-exporter
```

### Commandline options

```
-c config.yml  # path of config file, see config.dist.yml for an example
```

## Metrics

```
# TYPE onion_service_latency gauge
onion_service_latency{address="7sk2kov2xwx6cbc32phynrifegg6pklmzs7luwcggtzrnlsolxxuyfyd.onion",name="website"} 2
# TYPE onion_service_up gauge
onion_service_up{address="7sk2kov2xwx6cbc32phynrifegg6pklmzs7luwcggtzrnlsolxxuyfyd.onion",name="website"} 1
```

### Docker

```
docker run -p 9999:9999 -v /path/to/config.yml:/config.yml:ro systemli/prometheus-onion-service-exporter:latest -c config.yml
```

## License

GPLv3
