package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/grafana/pkg/models"
)

type UpdateFolderErrorResponse = struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

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

// SetFolder saves a folder with the given UID.
// If a folder with the same UID exists, it will be updated
// If UID is omitted, a new folder will be created.
// If the folder does not exist, it will be created.
func (r *Client) SetFolder(folder models.Folder) (models.Folder, error) {
	var (
		fo  models.Folder
		err error
	)
	// search for the folder by UID
	if fo, err = r.GetFolder(folder.Uid); err != nil {
		return models.Folder{}, fmt.Errorf("Could not check if folder %s exists: %w", folder.Uid, err)
	}

	if fo.Uid != "" {
		// folder doesn't exist
		return r.createFolder(folder.Uid, folder.Title)
	}

	return r.updateFolder(folder.Uid, folder.Title, folder.Version, false)

}

func (r *Client) createFolder(uid string, title string) (models.Folder, error) {
	var (
		raw      []byte
		fo       models.Folder
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
		return models.Folder{}, err
	}
	if code != 200 {
		return models.Folder{}, fmt.Errorf("HTTP error %d, returns %s", code, raw)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&fo); err != nil {
		return models.Folder{}, fmt.Errorf("error unmarshal folder %s: %w", uid, err)
	}
	return fo, err
}

func (r *Client) updateFolder(uid string, title string, version int, overwrite bool) (models.Folder, error) {
	var (
		raw      []byte
		fo       models.Folder
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
		return models.Folder{}, err
	}
	if code == 412 {
		// unpack the response and provide the message back to the user
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&failMsg); err != nil {
			return models.Folder{}, fmt.Errorf("Received HTTP 412 but was unable to understand server message: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Unable to update folder %s with error:\n%s: %s", uid, failMsg.Status, failMsg.Message)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&fo); err != nil {
		return models.Folder{}, fmt.Errorf("unable to parse server message: %w", err)
	}
	return fo, nil
}
