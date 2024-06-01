# Fundamentals

* https://promlabs.com/promql-cheat-sheet/

* [Prometheus Fundamentals](https://www.youtube.com/playlist?list=PLyBW7UHmEXgylLwxdVbrBQJ-fJ_jMvh8h)

* http://localhost:9090/

* prometheus_tsdb_head_samples_appended_total
* rate(prometheus_tsdb_head_samples_appended_total[1m])
* histogram_quantile(0.9, sum by(le,path) (rate(demo_api_request_duration_seconds_bucket[5m])))