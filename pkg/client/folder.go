package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana/pkg/models"
)

// GetAllFolders gets all folders.
// Reflects GET /api/folders API call.
func (r *Client) GetAllFolders() ([]models.Folder, error) {
	var (
		raw     []byte
		code    int
		folders []models.Folder
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
func (r *Client) GetFolder(uid string) (models.Folder, error) {
	var (
		raw  []byte
		fo   models.Folder
		code int
		err  error
	)
	if raw, code, err = r.get(fmt.Sprintf("api/folders/%s", uid), nil); err != nil {
		return models.Folder{}, err
	}
	if code != 200 {
		return models.Folder{}, fmt.Errorf("HTTP error %d: returns %s", code, raw)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&fo); err != nil {
		return models.Folder{}, fmt.Errorf("unmarshal board with meta: %s\n%s", err, raw)
	}
	return fo, err
}
