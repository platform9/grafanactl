/*
Copyright © 2020 Platform9 Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/platform9/grafana-sync/pkg/client"
	"github.com/spf13/cobra"
)

// folder command does not do anything, but is needed for scoping of subcommands
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Perform operations on Grafana Folders",
	Long:  "Perform operations on Grafana Folders",
}

func init() {
	rootCmd.AddCommand(folderCmd)
}

// isDirectoryMatch inspects a target directory to see if it matches the current grafana folder
func isDirectoryMatch(newFolder client.GrafanaFolder, targetDirectory string) (bool, error) {
	var (
		folderJSONPath string
		folderJSONRaw  []byte
		exists         os.FileInfo
		targetFolder   client.GrafanaFolder
		err            error
	)
	folderJSONPath = filepath.Join(targetDirectory, ".folder.json")
	if exists, _ = os.Lstat(folderJSONPath); exists == nil {
		return false, fmt.Errorf(".folder.json doesn't exist in target directory %s", targetDirectory)
	}
	if folderJSONRaw, err = ioutil.ReadFile(folderJSONPath); err != nil {
		return false, fmt.Errorf("Unable to read %s: %w", folderJSONPath, err)
	}
	if err = json.Unmarshal(folderJSONRaw, &targetFolder); err != nil {
		return false, fmt.Errorf("Unable to unmarshal the JSON in %s: %w", folderJSONPath, err)
	}
	return newFolder == targetFolder, nil
}
