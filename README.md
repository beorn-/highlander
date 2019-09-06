# highlander proxy
Highly available reverse proxy handling HA for prometheus server groups `remote_write` feature.

"There Can Be Only One".

## Introduction
Prometheus does not handle HA *per se*. It uses configurable alerts deduplication in alertmanager. So if you run `n` intances of prometheus to avoid Single Point of Failure (SPOF) on the collection/metrics/alerts side you get alerted only once.

Which is, for the generic usecase, absolutely fine. see https://github.com/prometheus/prometheus/issues/1500

When you want to store metrics for a longer time, rapidly you'll want to use the `remote_write` directive to push prometheus
data chunks to a Long Term Storage (LTS) tsdb. Moreover using a LTS allows users to use the LTS to display metrics, thus lower unavailability risk for alert manager And regroups over a long time all the metrics you have.

If you were to push data chunks to a single LTS you would have to avoid dataseries naming collision by adding a unique identification label to all metrics.

When you already have high cardinality and a finite budget allocated to observability this could prove difficult :
 - you multiply dataseries cardinality by `n` (also costs)
 - you would have to handle the prometheis instances labels in every single query that you have (aggregations usually don't do well with double data) in grafana/alerts

## What Highlander proxy is

* Highlander is a highly available L4 reverse proxy accepting `n` connections
* Highlander uses hashicorp `memberlist` golang library for clustering/leader election
* Highlander picks the first available connection and stores its uint64 `highlander_weight` url parameter
* Highlander prefers any new connection with a bigger uint64 value in the `highlander_weight` url parameter
* Highlander outputs to one `remote_write` compatible external endpoint

## What Highlander proxy will be

* Highlander will be a L7 reverse proxy accepting `n` prometheus remote write protocol connections
* Highlander will buffers data
* Highlander will outputs one or many  `remote_write` compatible external endpoint

## Architecture schema

```
[prometheus A]--------|-------------|
[prometheus B]--------| highlander  |------(remote_write A)-----------> [prometheus compatible LTS]
[prometheus C]--------|     proxy   |
[prometheus D]--------|-------------|
```
if prometheus A fails to push then highlander would pick another source
```
[prometheus A]---x    |-------------|
[prometheus B]--------| highlander  |------(remote_write B)-----------> [prometheus compatible LTS]
[prometheus C]--------|     proxy   |
[prometheus D]--------|-------------|
```

## Sample Log

```
2019/08/29 17:03:43 binding address: 0.0.0.0:9091
2019/08/29 17:03:43 remote address: http://localhost:9090
2019/08/29 17:03:43 health tick: 1s
2019/08/29 17:03:43 health expiry: 5s
2019/08/29 17:03:43 instance ID: 1
2019/08/29 17:03:44 new source : '127.0.0.1:48876' (bigger Highlander Weight) (0 -> 50)
2019/08/29 17:03:55 lost source : '127.0.0.1:48876' (no data for 5s)
2019/08/29 17:03:56 new source : '127.0.0.1:48880' (no current promoted source) (0)
```
