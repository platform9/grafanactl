/*
Copyright © 2019 Platform9 Systems

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

	"github.com/grafana-tools/sdk"
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
				rawBoard []byte
				meta     sdk.BoardProperties
				err      error
			)
			// download all dashboards into a directory
			gc := getGrafanaClientInternal()
			boards, err := gc.SearchDashboards(url.Values{})
			c := getGrafanaClient()
			for _, link := range boards {
				// Download the dashboard
				if rawBoard, meta, err = c.GetRawDashboard(link.Uri); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Error downloading: %s for %s\n", err, link.Uri))
					continue
				}
				// Write the dashboard to file in the target dir
				path := fmt.Sprintf("%s/%s.json", viper.GetString("target"), meta.Slug)
				if err = ioutil.WriteFile(path, rawBoard, os.FileMode(int(0666))); err != nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Error writing: %s\n", err))
					continue
				}
			}
		}
	},
}

func init() {
	dashboardCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolP("all", "a", false, "Download all dashboards")
	viper.BindPFlag("all", downloadCmd.Flags().Lookup("all"))

	downloadCmd.Flags().StringP("target", "t", ".", "Target directory to save dashboard files.")
	viper.BindPFlag("target", downloadCmd.Flags().Lookup("target"))
}
