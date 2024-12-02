## Setup

No Docker because it's actually simpler here.

Run nomad:

```
nomad agent -dev -bind 0.0.0.0 - network-interface='{{ GetDefaultInterfaces | attr "name" }}'
```

Run nomad-autoscaler:

```
# git clone git@github.com:hashicorp/nomad-autoscaler.git and then
go run main.go agent -config=$HOME/kb/nomad-autoscaler-playground/autoscaler.hcl
```

Run prometheus:

```
prometheus --config.file=./prometheus.yaml
```

Run metrics server:
```
cd metricsserver
go run main.go
```

## Play

* Prometheus GUI: http://localhost:9090
* Create/update nomad job: nomad job run job.hcl
* Nomad GUI: http://localhost:4646
* Set metrics: curl 'http://localhost:8090/set?cpu=50&scaletozero=1'
