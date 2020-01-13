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
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func prepareTable(searchResults []models.SearchHit) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Title", "Tags", "isStarred"})
	for _, hit := range searchResults {
		id := strconv.FormatUint(uint64(hit.Id), 10)
		isStarred := strconv.FormatBool(hit.IsStarred)
		tags := strings.Join(hit.Tags, ", ")
		table.Append([]string{id, hit.Title, tags, isStarred})
	}
	return table
}

func getSearchParams(args []string) url.Values {
	params := url.Values{}

	// query
	if len(args) > 0 {
		params.Set("query", strings.Join(args, " "))
	}
	// prefer flag to positional args
	query := viper.GetString("query")
	if query != "" {
		params.Set("query", query)
	}

	// tags
	tags := viper.GetStringSlice("tags")
	if len(tags) > 0 {
		for _, tag := range tags {
			params.Add("tags", tag)
		}
	}

	// dashboardIds
	dashboardIds := viper.GetIntSlice("dashboardIds")
	if len(dashboardIds) > 0 {
		for _, id := range dashboardIds {
			params.Add("dashboardIds", strconv.FormatInt(int64(id), 10))
		}
	}

	// folderIds
	folderIds := viper.GetIntSlice("folderIds")
	if len(folderIds) > 0 {
		for _, id := range folderIds {
			params.Add("folderIds", strconv.FormatInt(int64(id), 10))
		}
	}

	// starred
	starred := viper.GetBool("starred")
	if starred {
		params.Set("starred", strconv.FormatBool(starred))
	}

	return params
}

func loadSearchFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("query", "q", "", "Search Query")
	viper.BindPFlag("query", cmd.Flags().Lookup("query"))

	cmd.Flags().StringSliceP("tag", "t", []string{}, "List of tags to search for.")
	viper.BindPFlag("tags", cmd.Flags().Lookup("tag"))

	cmd.Flags().IntSliceP("dashboard-id", "d", []int{}, "List of dashboard id's to search for")
	viper.BindPFlag("dashboardIds", cmd.Flags().Lookup("dashboard-id"))

	cmd.Flags().IntSliceP("folder-id", "f", []int{}, "List of folder id's to search in for dashboards")
	viper.BindPFlag("folderIds", cmd.Flags().Lookup("folder-id"))

	cmd.Flags().Bool("starred", false, "Flag indicating if only starred Dashboards should be returned")
	viper.BindPFlag("starred", cmd.Flags().Lookup("starred"))
}
