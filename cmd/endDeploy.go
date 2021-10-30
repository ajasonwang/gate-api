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
	"encoding/json"
	"flag"
	"fmt"
	"reflect"

	"log"

	"github.com/spf13/cobra"
)

type PostData struct {
	UserName string
	Password string
	Id       string
}

func endDeploy(uri string, gateLoginUser string, gateLoginPass string) {

	deployId := readFile("deployid.txt")
	log.Println(deployId)

	data := PostData{
		UserName: gateLoginUser,
		Password: gateLoginPass,
		Id:       deployId,
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Type of data is:%v, POST data is:", reflect.TypeOf(jsonData), string(jsonData))

}

// endDeployCmd represents the endDeploy command
var endDeployCmd = &cobra.Command{
	Use:   "endDeploy",
	Short: "post deploy API",
	Long:  `post deploy API`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("endDeploy called")

		flag.Parse()

		if len(uri) == 0 {
			log.Fatal("Gate Uri不能为空")
		}
		if len(gu) == 0 {
			log.Fatal("Gate登录用户名不能为空")
		}
		if len(gp) == 0 {
			log.Fatal("Gate登录密码不能为空")
		}

		endDeploy(uri, gu, gp)
	},
}

func init() {

	endDeployCmd.Flags().StringVar(&uri, "uri", "192.168.100.100:1180", "gate uri")
	endDeployCmd.Flags().StringVar(&gu, "gu", "", "gate login username")
	endDeployCmd.Flags().StringVar(&gp, "gp", "", "gate login password")

	endDeployCmd.MarkFlagRequired("uri")
	endDeployCmd.MarkFlagRequired("gu")
	endDeployCmd.MarkFlagRequired("gp")

	rootCmd.AddCommand(endDeployCmd)

}
