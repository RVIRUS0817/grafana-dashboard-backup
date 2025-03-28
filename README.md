# grafana-dashboard-backup

This program backs up Grafana dashboards by fetching them via the Grafana API  
and saving them as JSON files in a directory structure that mirrors the Grafana folder organization.

# Setting 

- env
  - GRAFANA_API_KEY
  - MONITORING_URL

# Run

```bash
$ go run main.go
```
