/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "CM & QA Gate API intergration",
	Long:  `CM & QA Gate API intergration - trigger QA tests by passing deployment param to Gate API`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("GateAPIToolkit By CM Team 2021 - v1.0")

		program := "GateAPI"
		version := "v2.0"
		osArch := runtime.GOOS + "/" + runtime.GOARCH
		golangVersion := runtime.Version()
		log.Printf("%s %s - %s By CM Team, %s", program, version, osArch, golangVersion)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
