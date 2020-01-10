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
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload Grafana Dashboards",
	Long: `Upload Grafana Dashboards

Only files with a '.json' extension will be uploaded.`,
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()

		// Check the requested file/dir exists
		targetFiles, err := os.Stat(viper.GetString("files"))
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
			os.Exit(1)
		}
		var (
			files    []os.FileInfo
			rawBoard []byte
			readErr  error
		)
		// Check if file is a Dir or a File
		switch mode := targetFiles.Mode(); {
		case mode.IsDir():
			files, readErr = ioutil.ReadDir(targetFiles.Name())
			if readErr != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", readErr))
				os.Exit(1)
			}
		case mode.IsRegular():
			files = []os.FileInfo{targetFiles}
		}

		c := getGrafanaClient()

		for _, file := range files {
			fmt.Printf("uploading %s\n", file.Name())
			// check the file to see if it's a valid dashboard
			if strings.HasSuffix(file.Name(), ".json") {
				if rawBoard, err = ioutil.ReadFile(file.Name()); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to read file %s: %s\n", file.Name(), err))
					continue
				}
				var board sdk.Board
				if err = json.Unmarshal(rawBoard, &board); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Skipping file %s for error: %s\n", file.Name(), err))
					continue
				}
				c.DeleteDashboard(board.UpdateSlug())
				msg, err := c.SetDashboard(board, false)
				if err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Unable to upload dashboard %s:\n%s\n", file.Name(), err))
					continue
				} else {
					if msg.URL != nil && msg.Slug != nil {
						link := fmt.Sprintf("%s%s", viper.GetString("url"), (*msg.URL))
						fmt.Printf("Successfully Uploaded dashboard '%s'\nurl: %s\n", (*msg.Slug), link)
					}
				}
			}
		}
	},
}

func init() {
	dashboardCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP(
		"files", "f", ".", "Target file or directory of dashboard files to upload.")
	viper.BindPFlag("files", uploadCmd.Flags().Lookup("files"))
}
