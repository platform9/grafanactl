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
	"os"

	"github.com/spf13/cobra"
)

// searchDashboardCmd implements searchCmd for dash-dbs
var searchDashboardCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for Dashboards",
	Long:  `Search for Dashboards`,
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()
		c := getGrafanaClientInternal()
		queryParams := getSearchParams(cmd, args)
		results, _ := c.SearchDashboards(queryParams)
		if len(results) == 0 {
			fmt.Println("No results found.")
			os.Exit(0)
		}
		table := prepareTable(results)
		table.Render()
		os.Exit(0)
	},
}

func init() {
	dashboardCmd.AddCommand(searchDashboardCmd)
	loadSearchFlags(searchDashboardCmd)
}
