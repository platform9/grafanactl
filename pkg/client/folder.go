package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/models"
)

type UpdateFolderErrorResponse = struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type GrafanaFolder struct {
	ID    int64  `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
	URL   string `json:"url"`

	HasACL   bool `json:"hasAcl"`
	CanSave  bool `json:"canSave"`
	CanEdit  bool `json:"canEdit"`
	CanAdmin bool `json:"canAdmin"`

	CreatedBy string    `json:"createdBy"`
	Created   time.Time `json:"created"`
	UpdatedBy string    `json:"updatedBy"`
	Updated   time.Time `json:"updated"`
	Version   int       `json:"version"`
}

// GetAllFolders gets all folders.
// Reflects GET /api/folders API call.
func (r *Client) GetAllFolders() ([]GrafanaFolder, error) {
	var (
		raw     []byte
		code    int
		folders []GrafanaFolder
		err     error
	)
	if raw, code, err = r.get("api/folders", nil); err != nil {
		return nil, err
	}
	if code != 200 {
		return nil, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	err = json.Unmarshal(raw, &folders)
	return folders, err
}

// GetFolder gets a folder with the given UID.
// Reflects GET /api/folders/:uid API call.
func (r *Client) GetFolder(uid string) (GrafanaFolder, error) {
	var (
		raw  []byte
		fo   GrafanaFolder
		code int
		err  error
	)
	if raw, code, err = r.get(fmt.Sprintf("api/folders/%s", uid), nil); err != nil {
		return GrafanaFolder{}, err
	}
	if code == 404 {
		return GrafanaFolder{}, nil
	} else if code != 200 {
		return GrafanaFolder{}, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}

	if err = json.Unmarshal(raw, &fo); err != nil {
		return GrafanaFolder{}, fmt.Errorf("unmarshal board with meta: %s\n%s", err, raw)
	}

	return fo, err
}

// SetFolder saves a folder with the given UID.
// If a folder with the same UID exists, it will be updated
// If UID is omitted, a new folder will be created.
// If the folder does not exist, it will be created.
func (r *Client) SetFolder(folder GrafanaFolder) (GrafanaFolder, error) {
	var (
		fo  GrafanaFolder
		err error
	)
	// search for the folder by UID
	if fo, err = r.GetFolder(folder.UID); err != nil {
		return GrafanaFolder{}, fmt.Errorf("Could not check if folder %s exists: %w", folder.UID, err)
	}

	if fo.UID == "" {
		// folder doesn't exist
		fmt.Printf("Creating new folder %s (%s)", folder.Title, folder.UID)
		return r.createFolder(folder.UID, folder.Title)
	}

	return r.updateFolder(folder.UID, folder.Title, folder.Version, false)
}

func (r *Client) createFolder(uid string, title string) (GrafanaFolder, error) {
	var (
		raw      []byte
		fo       GrafanaFolder
		code     int
		err      error
		toCreate models.CreateFolderCommand
		payload  []byte
	)
	toCreate = models.CreateFolderCommand{
		Uid:   uid,
		Title: title,
	}
	payload, _ = json.Marshal(toCreate)
	if raw, code, err = r.post("api/folders/", nil, payload); err != nil {
		return GrafanaFolder{}, err
	}
	if code != 200 {
		return GrafanaFolder{}, fmt.Errorf("HTTP error %d, returns %s", code, raw)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&fo); err != nil {
		return GrafanaFolder{}, fmt.Errorf("error unmarshal folder %s: %w", uid, err)
	}
	return fo, err
}

func (r *Client) updateFolder(uid string, title string, version int, overwrite bool) (GrafanaFolder, error) {
	var (
		raw      []byte
		fo       GrafanaFolder
		code     int
		err      error
		toUpdate models.UpdateFolderCommand
		failMsg  UpdateFolderErrorResponse
		payload  []byte
	)
	toUpdate = models.UpdateFolderCommand{
		Uid:       uid,
		Title:     title,
		Version:   version,
		Overwrite: overwrite,
	}
	payload, _ = json.Marshal(toUpdate)
	if raw, code, err = r.put(fmt.Sprintf("api/folders/%s", uid), nil, payload); err != nil {
		return GrafanaFolder{}, err
	}
	if code == 412 {
		// unpack the response and provide the message back to the user
		if err := json.Unmarshal(raw, &failMsg); err != nil {
			return GrafanaFolder{}, fmt.Errorf("Received HTTP 412 but was unable to understand server message: %w", err)
		}
		return GrafanaFolder{}, fmt.Errorf(fmt.Sprintf("%s: %s\n", failMsg.Status, failMsg.Message))
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&fo); err != nil {
		return GrafanaFolder{}, fmt.Errorf("unable to parse server message: %w", err)
	}
	return fo, nil
}
