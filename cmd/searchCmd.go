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

// merges two int slices, ensuring that no duplicate ints are in the resulting slice
// May be combined with mergeStringSlices(..) if we can figure out how to handle multiple
//   types without too much interface{} juggling
func mergeIntSlices(slice1 []int, slice2 []int) []int {
	for _, item2 := range slice2 {
		for _, item1 := range slice1 {
			if item1 == item2 {
				// item2 exists already, go to next item2
				break
			}
		}
		// item2 doesn't exist, append to the slice
		slice1 = append(slice1, item2)
	}
	return slice1
}

// merges two string slices, ensuring that no duplicate ints are in the resulting slice
// May be combined with mergeIntSlices(..) if we can figure out how to handle multiple
//   types without too much interface{} juggling
func mergeStringSlices(slice1 []string, slice2 []string) []string {
	for _, item2 := range slice2 {
		for _, item1 := range slice1 {
			if item1 == item2 {
				// item2 exists already, go to next item2
				break
			}
		}
		// item2 doesn't exist, append to the slice
		slice1 = append(slice1, item2)
	}
	return slice1
}

func getSearchParams(cmd *cobra.Command, args []string) url.Values {
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
	// Note: viper.GetStringSlice() does not seem to properly bind the flag
	// as a workaround, we will merge the lists ourselves
	tagsViper := viper.GetStringSlice("tag")
	tagsCobra, _ := cmd.Flags().GetStringSlice("tag")
	tags := mergeStringSlices(tagsViper, tagsCobra)
	if len(tags) > 0 {
		for _, tag := range tags {
			params.Add("tag", tag)
		}
	}

	// dashboardIds
	// Note: viper.GetIntSlice() does not seem to properly bind the flag
	// as a workaround, we will merge the lists ourselves
	dashboardsViper := viper.GetIntSlice("dashboard")
	dashboardsCobra, _ := cmd.Flags().GetIntSlice("dashboard")
	dashboards := mergeIntSlices(dashboardsViper, dashboardsCobra)
	if len(dashboards) > 0 {
		for _, id := range dashboards {
			params.Add("dashboardIds", strconv.FormatInt(int64(id), 10))
		}
	}

	// folderIds
	// Note: viper.GetIntSlice() does not seem to properly bind the flag
	// as a workaround, we will merge the lists ourselves
	foldersViper := viper.GetIntSlice("folder")
	foldersCobra, _ := cmd.Flags().GetIntSlice("folder")
	folders := mergeIntSlices(foldersViper, foldersCobra)
	if len(folders) > 0 {
		for _, id := range folders {
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
	cmd.Flags().StringSliceP("tag", "t", []string{}, "List of tags to search for.")
	cmd.Flags().IntSliceP("dashboard", "d", []int{}, "List of dashboard id's to search for")
	cmd.Flags().IntSliceP("folder", "f", []int{}, "List of folder id's to search in for dashboards")
	cmd.Flags().Bool("starred", false, "Flag indicating if only starred Dashboards should be returned")
	viper.BindPFlags(cmd.Flags())
}
