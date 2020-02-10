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
	"strings"

	"github.com/platform9/grafana-sync/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload Grafana Dashboards",
	Long: `Upload Grafana Dashboards

Only files with a '.json' extension will be uploaded.`,
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()

		// Check the requested file/dir exists
		targetFiles, err := os.Lstat(viper.GetString("files"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
			os.Exit(1)
		}
		var (
			files   []os.FileInfo
			readErr error
		)
		// Check if file is a Dir or a File
		switch mode := targetFiles.Mode(); {
		case mode.IsDir():
			// Enumerate a list of files in the directory
			files, readErr = ioutil.ReadDir(targetFiles.Name())
			if readErr != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", readErr))
				os.Exit(1)
			}
		case mode.IsRegular():
			// Enumerate a list of files with just this file
			files = []os.FileInfo{targetFiles}
		}

		c := getGrafanaClient()

		for _, file := range files {
			if file.Mode().IsDir() {
				// Check if the folder has a signature
				var (
					dashboardDir   string
					folderJSONPath string
					folderJSONRaw  []byte
					folderJSON     client.GrafanaFolder
					err            error
					folder         client.GrafanaFolder
				)
				dashboardDir = filepath.Join(targetFiles.Name(), file.Name())
				folderJSONPath = filepath.Join(dashboardDir, ".folder.json")
				if _, err = os.Lstat(folderJSONPath); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Couldn't find .folder.json found for directory %s: %s\n", file.Name(), err))
					continue
				}
				if folderJSONRaw, err = ioutil.ReadFile(folderJSONPath); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to read file: %s\nError: %s", folderJSONPath, err))
					continue
				}
				if err = json.Unmarshal(folderJSONRaw, &folderJSON); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to unmarshal file: %v\nError: %s", folderJSONRaw, err))
					continue
				}

				// Use the folder as returned by create/update to get the correct ID
				if folder, err = c.SetFolder(folderJSON); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Error setting folder '%s': %s\n", folderJSON.Title, err))
					continue
				}
				if folder.ID == 0 {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to resolve the real folder ID. Skipping folder '%s'", folderJSON.Title))
					continue
				}
				files, readErr = ioutil.ReadDir(filepath.Join(targetFiles.Name(), file.Name()))
				if readErr != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", readErr))
				}
				uploadFiles(files, dashboardDir, int(folder.ID), viper.GetBool("overwrite"))
				continue
			}
			uploadFiles([]os.FileInfo{file}, targetFiles.Name(), 0, viper.GetBool("overwrite"))
		}
	},
}

func uploadFiles(files []os.FileInfo, basePath string, targetFolderID int, overwrite bool) error {
	c := getGrafanaClient()
	for _, file := range files {
		if file.Mode().IsDir() {
			return fmt.Errorf("uploadFiles will not upload directories")
		}

		var (
			rawBoard []byte
			err      error
		)

		dashboardFile := filepath.Join(basePath, file.Name())
		if !strings.HasSuffix(dashboardFile, ".json") {
			fmt.Printf("Skipping '%s' (Not a JSON file)\n", file.Name())
			continue
		}

		if rawBoard, err = ioutil.ReadFile(dashboardFile); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to read file %s: %s\n", dashboardFile, err))
			continue
		}

		// Replace the dashboard
		if err = c.SetDashboard(rawBoard, overwrite, targetFolderID); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to upload %s: %s\n", dashboardFile, err))
			continue
		}
	}
	return nil
}

func init() {
	dashboardCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP(
		"files", "f", ".", "Target file or directory of dashboard files to upload.")
	uploadCmd.Flags().Bool("overwrite", false, "Overwrite existing dashboard with newer version, same dashboard title in folder, or same dashboard UID.")
	viper.BindPFlags(uploadCmd.Flags())
}
