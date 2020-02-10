package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
)

type DashboardUploadRequest struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	FolderID  int                    `json:"folderId"`
	Overwrite bool                   `json:"overwrite"`
}

type DashboardUploadResponse struct {
	ID      int    `json:"id"`
	UID     string `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"string"`
	Version int    `json:"version"`
}

type PreconditionFailedMsg struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// GrafanaDashboardFullWithMeta is copied from
// github.com/grafana/grafana/pkg/api/dtos. this wasn't vendored because
// importing dtos causes module errors on the go-xorm/core module.
type GrafanaDashboardFullWithMeta struct {
	Meta      GrafanaDashboardMeta `json:"meta"`
	Dashboard *simplejson.Json     `json:"dashboard"`
}

// GrafanaDashboardMeta is copied from github.com/grafana/grafana/pkg/api/dtos
// this wasn't vendored because importing dtos causes module errors on the
// go-xorm/core module.
type GrafanaDashboardMeta struct {
	IsStarred             bool      `json:"isStarred,omitempty"`
	IsHome                bool      `json:"isHome,omitempty"`
	IsSnapshot            bool      `json:"isSnapshot,omitempty"`
	Type                  string    `json:"type,omitempty"`
	CanSave               bool      `json:"canSave"`
	CanEdit               bool      `json:"canEdit"`
	CanAdmin              bool      `json:"canAdmin"`
	CanStar               bool      `json:"canStar"`
	Slug                  string    `json:"slug"`
	Url                   string    `json:"url"`
	Expires               time.Time `json:"expires"`
	Created               time.Time `json:"created"`
	Updated               time.Time `json:"updated"`
	UpdatedBy             string    `json:"updatedBy"`
	CreatedBy             string    `json:"createdBy"`
	Version               int       `json:"version"`
	HasAcl                bool      `json:"hasAcl"`
	IsFolder              bool      `json:"isFolder"`
	FolderId              int64     `json:"folderId"`
	FolderTitle           string    `json:"folderTitle"`
	FolderUrl             string    `json:"folderUrl"`
	Provisioned           bool      `json:"provisioned"`
	ProvisionedExternalId string    `json:"provisionedExternalId"`
}

// GetDashboard will query the grafana api for a specific dashboard by UID
// Reflects GET /api/dashboards/uid/:uid API call.
func (r *Client) GetDashboard(uid string) (GrafanaDashboardFullWithMeta, error) {
	var (
		raw  []byte
		dash GrafanaDashboardFullWithMeta
		code int
		err  error
	)
	if raw, code, err = r.get(fmt.Sprintf("api/dashboards/uid/%s", uid), nil); err != nil {
		return dash, err
	}
	if code != 200 {
		return dash, err
	}
	if err = json.Unmarshal(raw, &dash); err != nil {
		return dash, err
	}
	return dash, nil
}

// SetDashboard will create or update a new/existing dashboard
// Reflects POST /api/dashboards/db API call.
func (r *Client) SetDashboard(dash []byte, overwrite bool, folderID int) error {
	var (
		raw               []byte
		req               DashboardUploadRequest
		resp              DashboardUploadResponse
		code              int
		payload           []byte
		err               error
		dashTitle         string
		dashUID           string
		existingDashboard GrafanaDashboardFullWithMeta
	)

	// construct a valid payload out of the dashboard, folderID, and overwrite flag
	// unmarshal then unmarshal the json to serialize the correct payload
	// To avoid creating a very complex dashboard struct, we'll use a generic interface
	var dashboardContents map[string]interface{}
	_ = json.Unmarshal(dash, &dashboardContents)
	// store dashboard's title for more friendly/usable messages
	dashTitle = fmt.Sprintf("%v", dashboardContents["title"])
	dashUID = fmt.Sprintf("%v", dashboardContents["uid"])

	if dashboardContents["panels"] == nil {
		return fmt.Errorf("not a dashboard")
	}

	// check if a dashboard already exists
	existingDashboard, _ = r.GetDashboard(dashUID)
	if (GrafanaDashboardFullWithMeta{}) == existingDashboard {
		// assume dashboard exists
		fmt.Printf("Dashboard exists already, creating it as a new dashboard")
		// strip the ID so a new dashboard will be created
		delete(dashboardContents, "id")
	}

	// construct the request body
	req = DashboardUploadRequest{
		Dashboard: dashboardContents,
		FolderID:  folderID,
		Overwrite: overwrite,
	}
	payload, _ = json.Marshal(req)

	// submit the request
	if raw, code, err = r.post("api/dashboards/db", nil, payload); err != nil {
		return err
	}
	if code == 412 {
		var badthings PreconditionFailedMsg
		json.Unmarshal(raw, &badthings)
		return fmt.Errorf("%s: %s", badthings.Status, badthings.Message)
	} else if code != 200 {
		// attempt to unmarshal the raw payload and display the error
		var (
			msg       string
			badthings PreconditionFailedMsg
		)
		// fallback to the raw message
		msg = string(raw)
		// if possible, unmarshal and display a formatted error
		if err = json.Unmarshal(raw, &badthings); err == nil {
			msg = fmt.Sprintf("HTTP %d: %s", code, badthings.Message)
		}
		return fmt.Errorf(msg)
	}

	if err = json.Unmarshal(raw, &resp); err != nil {
		return err
	}
	fmt.Printf("Updated dashboard %s (%s) successfully!\n", dashTitle, resp.UID)
	return nil
}
