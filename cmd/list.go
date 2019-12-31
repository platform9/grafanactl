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
	"os"
	"strconv"

	"github.com/grafana-tools/sdk"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Grafana Dashboards",
	Long:  `List Grafana Dashboards`,
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()
		var (
			boardLinks []sdk.FoundBoard
			err        error
		)

		c := sdk.NewClient(viper.GetString("url"), viper.GetString("apikey"), sdk.DefaultHTTPClient)
		if boardLinks, err = c.SearchDashboards("", false); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Title", "URI", "Type"})
		for _, link := range boardLinks {
			// convert the ID to string
			id := strconv.FormatUint(uint64(link.ID), 10)
			// append the line to the data table
			table.Append([]string{id, link.Title, link.URI, link.Type})
		}
		table.Render()
	},
}

func init() {
	dashboardCmd.AddCommand(listCmd)
}
