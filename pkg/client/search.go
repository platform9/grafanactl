package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/grafana/grafana/pkg/models"
)

// Search searches grafana dashboards and folders
// Reflects GET /api/search API call.
func (r *Client) Search(queryParams url.Values) ([]models.SearchHit, error) {
	var (
		raw   []byte
		code  int
		found []models.SearchHit
		err   error
	)
	fmt.Printf("Querying %s\n", queryParams.Encode())
	if raw, code, err = r.get("api/search", queryParams); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &found)
	return found, err
}

func (r *Client) SearchFolders(queryParams url.Values) ([]models.SearchHit, error) {
	queryParams.Set("type", "dash-folder")
	return r.Search(queryParams)
}

func (r *Client) SearchDashboards(queryParams url.Values) ([]models.SearchHit, error) {
	queryParams.Set("type", "dash-db")
	return r.Search(queryParams)
}
