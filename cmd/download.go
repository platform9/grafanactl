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
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/grafana-tools/sdk"
	"github.com/grafana/grafana/pkg/models"
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
				folders []models.Folder
				err     error
				gc      *client.Client
			)
			gc = getGrafanaClientInternal()

			// List all folders
			if folders, err = gc.GetAllFolders(); err != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("error downloading folders: %s\n", err))
				os.Exit(1)
			}
			for _, folder := range folders {
				var (
					dirName string
					err     error
				)
				dirName = filepath.Join(viper.GetString("target"), dirName)
				if err = mkdirFromFolder(folder); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("error creating folder %s: %s", dirName, err))
					fmt.Printf("Skipping download of folder '%s'", folder.Title)
					continue
				}
				if err = saveFolderDashboards(folder.Id, dirName); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("error saving dashboards for folder %s: %s", dirName, err))
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
		results   []models.SearchHit
		rawBoard  []byte
		meta      sdk.BoardProperties
		err       error
		client    *sdk.Client
		folderIDs string
	)
	folderIDs = strconv.FormatInt(folderID, 10)
	client = getGrafanaClient()
	c := getGrafanaClientInternal()
	query = url.Values{}
	query.Add("folderIds", folderIDs)
	if results, err = c.SearchDashboards(query); err != nil {
		return fmt.Errorf("error searching dashboards in folder %s: %w", folderIDs, err)
	}
	for _, board := range results {
		// Download the dashboard
		if rawBoard, meta, err = client.GetRawDashboard(board.Uri); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error downloading dashboard %s: %s\n", board.Uri, err))
			continue
		}
		// Write the dashboard to file
		path := filepath.Join(targetDir, fmt.Sprintf("%s.json", meta.Slug))
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
