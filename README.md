# Grafana Sync

Grafana Sync is a utility to aid in syncing Dashboards across grafana instances and/or organizations.

Read the [design spec](https://platform9.atlassian.net/wiki/spaces/MON/pages/671547611/)

## Usage

```bash
# Listing Dashboards
grafana-sync dashboard list

# Downloading dashboards
grafana-sync dashboard download --all
grafana-sync dashboard download --all -t dashboards
```
