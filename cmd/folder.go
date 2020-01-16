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
	"regexp"
	"strings"

	"github.com/grafana/grafana/pkg/models"
	"github.com/spf13/cobra"
)

// folder command does not do anything, but is needed for scoping of subcommands
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Perform operations on Grafana Folders",
	Long:  "Perform operations on Grafana Folders",
}

func init() {
	rootCmd.AddCommand(folderCmd)
}

func mkdirFromFolder(folder models.Folder) error {
	target := sanitizeDirName(folder.Title)
	var (
		signatureFile string
		signatureRaw  []byte
		exists        os.FileInfo
		err           error
	)
	signatureFile = filepath.Join(target, ".folder.json")
	if signatureRaw, err = json.Marshal(folder); err != nil {
		return fmt.Errorf("error marshal signature: %w", err)
	}
	exists, _ = os.Lstat(target)
	if exists == nil {
		// Attempt to create the directory
		if err := os.MkdirAll(target, 0744); err != nil {
			return fmt.Errorf("error creating directory %s: %w", target, err)
		}
		// Save the folder signature into the directory
		if err := ioutil.WriteFile(signatureFile, signatureRaw, 0666); err != nil {
			return fmt.Errorf("error writing %s: %w", signatureFile, err)
		}
	} else {
		// Check the folder signature
		var (
			savedSignatureRaw []byte
			savedSignature    models.Folder
			err               error
		)
		if savedSignatureRaw, err = ioutil.ReadFile(signatureFile); err != nil {
			return fmt.Errorf("error reading %s: %w", signatureFile, err)
		}
		if err = json.Unmarshal(savedSignatureRaw, &savedSignature); err != nil {
			return fmt.Errorf("error unmarshal signature: %w", err)
		}
		if savedSignature != folder {
			return fmt.Errorf("existing folder signature '%s' does not match", signatureFile)
		}
	}
	return nil
}

// makes a very unsafe directory name reasonable
func sanitizeDirName(name string) string {
	var (
		sanitizeRegex *regexp.Regexp
		dirName       string
	)
	sanitizeRegex, _ = regexp.Compile("[^A-Za-z0-9._-]")
	dirName = strings.ToLower(name)
	dirName = string(sanitizeRegex.ReplaceAll([]byte(dirName), []byte("_")))
	return dirName
}
