# Grafana Sync

Grafana Sync is a utility to aid in syncing Dashboards across grafana instances and/or organizations.

Read the [design spec](https://platform9.atlassian.net/wiki/spaces/MON/pages/671547611/)

## Usage

Optional arguments

* -O, --org-id -- the organization ID to access
* -F, --folder-id -- the folder to access (by numeric ID)
* -U, --uid -- the UID of the folder or dashboard to access 

```bash
# Downloading dashboards or folders
grafana-sync folder list
grafana-sync folder list --org-id 2
grafana-sync folder download
grafana-sync dashboard list
grafana-sync dashboard list --folder-id 35
grafana-sync dashboard download
grafana-sync dashboard download --org-id 2
grafana-sync dashboard download --folder-id 35
grafana-sync dashboard download --uid 4D32a3

# Uploading dashboards or folders
grafana-sync folder upload -f foo/
grafana-sync folder upload -f foo/ --org-id 3
grafana-sync dashboard upload -f foo/bar.json
grafana-sync dashboard upload -f foo/bar.json --folder-id 35
grafana-sync dashboard upload -f foo/bar.json --org-id 2
```
