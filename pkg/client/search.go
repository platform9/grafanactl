package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GrafanaSearchHit reflects the response of the folder/dashboard search API
type GrafanaSearchHit struct {
	ID        int      `json:"id"`
	UID       string   `json:"uid"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

// Search searches grafana dashboards and folders
// Reflects GET /api/search API call.
func (r *Client) Search(queryParams url.Values) ([]GrafanaSearchHit, error) {
	var (
		raw   []byte
		code  int
		found []GrafanaSearchHit
		err   error
	)
	if raw, code, err = r.get("api/search", queryParams); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &found)
	return found, err
}

func (r *Client) SearchFolders(queryParams url.Values) ([]GrafanaSearchHit, error) {
	queryParams.Set("type", "dash-folder")
	return r.Search(queryParams)
}

func (r *Client) SearchDashboards(queryParams url.Values) ([]GrafanaSearchHit, error) {
	queryParams.Set("type", "dash-db")
	return r.Search(queryParams)
}
