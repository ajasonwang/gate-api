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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"os"

	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"net/url"

	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func getEnvID(user string, pass string, projectId int, envName string) int {

	log.Println("Call getEnvID()")

	values := url.Values{}
	values.Set("userName", user)
	values.Set("password", pass)
	pId := strconv.Itoa(projectId)
	values.Set("jiraProjectId", pId)

	endpoint := "http://" + uri + "/api/api-getEnvByProjectId"
	log.Println("> 3st request - get env id by request gate api")
	req, err := http.NewRequestWithContext(traceCtx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	// req, _ := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(resp.StatusCode, err)
	} else {
		log.Println("Login gate sucess with status code:", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var Result struct {
		Code    string
		Message string
		Data    []struct {
			Id  string
			Env string
		}
	}

	if err := json.Unmarshal([]byte(body), &Result); err != nil {
		log.Fatal(err)
	}

	if len(Result.Data) == 0 {
		log.Fatalf("Call getEnvID(): env name:%s not found, please check!", envName)
	} else {
		for _, x := range Result.Data {
			if strings.Contains(strings.ToLower(x.Env), strings.ToLower(envName)) {
				log.Printf("Env:%s, EnvId:%s", strings.ToUpper(envName), x.Id)
				envId, err := strconv.Atoi(x.Id)
				if err != nil {
					log.Fatal(err)
				}
				return envId
			}
		}
	}

	resp.Body.Close()

	return 0
}

type DeployCode struct {
	Id string `json:"id"`
}
type ReturnData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    DeployCode
}

type DeployComp struct {
	Teams    string `json:"teams"`
	Stacks   string `json:"stacks"`
	Services string `json:"services"`
}
type DeployData struct {
	UserName      string     `json:"userName"`
	Password      string     `json:"password"`
	JiraProjectId int        `json:"jiraProjectId"`
	JiraVersionId string     `json:"jiraVersionId"`
	EnvId         int        `json:"envId"`
	Component     DeployComp `json:"component"`
}

func startDeploy(jiraLoginUser string, jiraLoginPass string, uri string, gateLoginUser string, gateLoginPass string, projectName string, versionName string, envName string, teams string, stacks string, services string) {

	log.Println("Call startDeploy()")

	projectId := getProjectID(client, jiraLoginUser, jiraLoginPass, projectName)
	b := checkReleaseVersion(client, projectName, versionName)
	log.Printf("Release version created in jira: %v", b)

	envId := getEnvID(gateLoginUser, gateLoginPass, projectId, envName)

	log.Printf("project:%s, project id:%d, version name:%s, env id:%d", projectName, projectId, versionName, envId)

	dc := DeployComp{
		teams,
		stacks,
		services,
	}
	data := DeployData{
		gateLoginUser,
		gateLoginPass,
		projectId,
		versionName,
		envId,
		dc,
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	gateAPI := "http://" + uri + "/api/api-release"
	req, err := http.NewRequestWithContext(traceCtx, http.MethodPost, gateAPI, bytes.NewBuffer(jsonData))
	req.SetBasicAuth(gateLoginUser, gateLoginPass)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(resp, err)
	}

	// body := "{\"code\":\"1\",\"message\":\"success\",\"data\":{\"id\":\"408\"}}"
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		log.Println("Call startDeploy(): response status code and content:", resp.StatusCode, string(body))
	}
	resp.Body.Close()

	var rd ReturnData

	err1 := json.Unmarshal([]byte(body), &rd)
	if err1 != nil {
		log.Fatalf("Error happened in JSON Unmarshal. Err: %s", err1)
	}
	log.Printf("Returned deploy code is:%s", rd.Data.Id)
	writeFile("deployid.txt", rd.Data.Id)

}

// startDeployCmd represents the startDeploy command
var startDeployCmd = &cobra.Command{
	Use:   "startDeploy",
	Short: "POST deployment data to Gate before start",
	Long:  `POST deployment data to Gate before start`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("startDeploy called")

		flag.Parse()

		if len(uri) == 0 {
			log.Fatal("Gate Uri不能为空")
		}
		if len(ju) == 0 {
			log.Fatal("JIRA登录用户名不能为空")
		}
		if len(jp) == 0 {
			log.Fatal("JIRA登录密码不能为空")
		}
		if len(gu) == 0 {
			log.Fatal("Gate登录用户名不能为空")
		}
		if len(gp) == 0 {
			log.Fatal("Gate登录密码不能为空")
		}
		if len(project) == 0 {
			log.Fatal("Project名称不能为空")
		}
		if len(versionname) == 0 {
			log.Fatal("JIRA版本名称不能为空")
		}
		r, _ := regexp.Compile(`^[a-z-A-Z]+_.*\d{1,3}$`)
		b := r.MatchString(versionname)
		if b != true {
			log.Fatal("JIRA版本名称格式不符合规定，正确的格式例子：APPX_1.21.22")
			os.Exit(1)
		}
		if len(envname) == 0 {
			log.Fatal("发布环境名称不能为空")
		}

		startDeploy(ju, jp, uri, gu, gp, project, versionname, envname, teams, stacks, services)
	},
}

func init() {

	startDeployCmd.Flags().StringVar(&ju, "ju", "", "jira login username")
	startDeployCmd.Flags().StringVar(&jp, "jp", "", "jira login password")
	startDeployCmd.Flags().StringVar(&uri, "uri", "192.168.100.100:1180", "gate uri")
	startDeployCmd.Flags().StringVar(&gu, "gu", "", "gate login username")
	startDeployCmd.Flags().StringVar(&gp, "gp", "", "gate login password")
	startDeployCmd.Flags().StringVar(&project, "project", "", "jira project name, eg: fleet")
	startDeployCmd.Flags().StringVar(&versionname, "versionname", "", "jira release version name, eg: SOAR_1.21.22")
	startDeployCmd.Flags().StringVar(&envname, "envname", "", "deploy environment name, eg: QA QA3")
	startDeployCmd.Flags().StringVar(&teams, "teams", "", "deployed teams(all services in a team), eg: ACE TAP")
	startDeployCmd.Flags().StringVar(&stacks, "stacks", "", "deployed rancher stacks, eg: TapTradingDataServices TapTradingServices")
	startDeployCmd.Flags().StringVar(&services, "services", "", "deployed rancher service names, eg: cmQueryService")

	startDeployCmd.MarkFlagRequired("ju")
	startDeployCmd.MarkFlagRequired("jp")
	startDeployCmd.MarkFlagRequired("uri")
	startDeployCmd.MarkFlagRequired("gu")
	startDeployCmd.MarkFlagRequired("gp")
	startDeployCmd.MarkFlagRequired("project")
	startDeployCmd.MarkFlagRequired("versionname")
	startDeployCmd.MarkFlagRequired("envname")

	rootCmd.AddCommand(startDeployCmd)

}
