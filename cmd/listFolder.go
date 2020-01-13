/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/grafana/grafana/pkg/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listFolderCmd = &cobra.Command{
	Use:   "list",
	Short: "List Grafana Folders",
	Long:  `List Grafana Folders`,
	Run: func(cmd *cobra.Command, args []string) {
		requireAuthParams()
		c := getGrafanaClientInternal()
		var (
			folders []models.Folder
			err     error
		)
		if folders, err = c.GetAllFolders(); err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", err))
			os.Exit(1)
		}
		if len(folders) < 1 {
			fmt.Println("No folders found.")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "UID", "Title"})
		for _, fol := range folders {
			id := strconv.FormatUint(uint64(fol.Id), 10)
			table.Append([]string{id, fol.Uid, fol.Title})
		}
		table.Render()
	},
}

func init() {
	folderCmd.AddCommand(listFolderCmd)
}
