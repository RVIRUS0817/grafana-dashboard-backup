# grafana-dashboard-backup

A simple tool to back up Grafana dashboards via API.

---

## ✨ Features

- Backup all Grafana dashboards as JSON
- Organize dashboards into a specified output directory
- Simple CLI interface

---

## 📦 Prerequisites

- Go 1.24+
- A Grafana API Key with `Editor` or `Admin` privileges
- Network access to your Grafana instance

---

## 🚀 Installation

- env
  - GRAFANA_API_KEY
  - MONITORING_URL

```bash
git clone https://github.com/RVIRUS0817/grafana-dashboard-backup.git
cd grafana-dashboard-backup
go run main.go
