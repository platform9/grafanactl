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
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/platform9/grafana-sync/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download dashboards from a grafana instance",
	Long:  `Download dashboards from a grafana instance`,
	// allow specification of dashboard ID as a positional argument
	// except do not error out if `--all` is set and no positional arg is specified
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()
		if !viper.GetBool("all") {
			// expect the user to have specified a positional argument
			if len(args) < 1 {
				fmt.Fprintln(os.Stderr, "You must specify a dashboard ID or use --all.")
				os.Exit(1)
			} else {
				fmt.Println("Individual dashboard download not implemented. try --all.")
			}
		} else {
			var (
				folders []client.GrafanaFolder
				err     error
				c       *client.Client
			)
			c = getGrafanaClient()

			// Prepare folder destinations
			if folders, err = c.GetAllFolders(); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("error downloading folders: %s\n", err))
				os.Exit(1)
			}
			for _, fol := range folders {
				// Sanitize the folder name
				sanitizeRegex, _ := regexp.Compile("[^A-Za-z0-9._-]")
				dirName := strings.ToLower(fol.Title)
				dirName = string(sanitizeRegex.ReplaceAll([]byte(dirName), []byte("_")))
				dirName = filepath.Join(viper.GetString("target"), dirName)
				signatureFile := filepath.Join(dirName, ".folder.json")

				// Check if a folder already exists
				exists, _ := os.Lstat(dirName)
				if exists == nil {
					// Attempt to create the directory
					err := os.MkdirAll(dirName, 0744)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error creating directory %s: %s\n", dirName, err)
						continue
					}
					// Save the folder signature into the directory
					var fileContents []byte
					fileContents, err = json.Marshal(fol)
					if err != nil {
						fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to marshal json: %v\nError: %s", fol, err))
						continue
					}
					if err = ioutil.WriteFile(signatureFile, fileContents, 0666); err != nil {
						fmt.Fprintf(os.Stderr, fmt.Sprintf("Error writing %s: %s\n", signatureFile, err))
						continue
					}
					saveFolderDashboards(fol.ID, dirName)
				} else {
					// Read the .folder.json file and unmarshal it
					var directoryIsFolder bool
					if directoryIsFolder, err = isDirectoryMatch(fol, dirName); err != nil {
						fmt.Fprintf(os.Stderr, "%s", err)
						continue
					}
					if directoryIsFolder {
						fmt.Printf("Existing directory '%s' matches the existing grafana folder '%s'. Overwriting.\n", dirName, fol.Title)
						saveFolderDashboards(fol.ID, dirName)
					} else {
						fmt.Println("Folder signatures don't match")
						fmt.Printf("The folder '%s' will be skipped\n", fol.Title)
						// TODO: append the UID to the folder name and try again
					}
					continue
				}
			}
			// Download all of the dashboards in the "General" folder (always has ID of 0)
			saveFolderDashboards(0, viper.GetString("target"))
		}
	},
}

// saveFolderDashboards will download all of the dashboards to the target dir
// It's expected that a folder with this ID and the target dir already exist
func saveFolderDashboards(folderID int64, targetDir string) error {
	var (
		query     url.Values
		results   []client.GrafanaSearchHit
		rawBoard  []byte
		dash      client.GrafanaDashboardFullWithMeta
		err       error
		folderIDs string
	)
	folderIDs = strconv.FormatInt(folderID, 10)
	c := getGrafanaClient()
	query = url.Values{}
	query.Add("folderIds", folderIDs)
	if results, err = c.SearchDashboards(query); err != nil {
		return fmt.Errorf("error searching dashboards in folder %s: %w", folderIDs, err)
	}
	for _, board := range results {
		// Download the dashboard
		if dash, err = c.GetDashboard(board.UID); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error downloading dashboard %s: %s\n", board.UID, err))
			continue
		}
		rawBoard, _ = dash.Dashboard.Encode()
		// Write the dashboard to file
		path := filepath.Join(targetDir, fmt.Sprintf("%s.json", dash.Meta.Slug))
		if err = ioutil.WriteFile(path, rawBoard, 0666); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error writing: %s\n", err))
			continue
		}
		fmt.Printf("Downloaded %s\n", path)
	}
	return nil
}

func init() {
	dashboardCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("all", "a", false, "Download all dashboards")
	downloadCmd.Flags().StringP("target", "t", ".", "Target directory to save dashboard files.")
	viper.BindPFlags(downloadCmd.Flags())
}
