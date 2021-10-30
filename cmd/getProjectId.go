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

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LoginData struct {
	UserName string
	Password string
}

func jiraLogin(endpoint string, jiraLoginUser string, jiraLoginPass string) *http.Client {

	log.Println("Call jiraLogin()")

	values := url.Values{}
	values.Set("os_username", jiraLoginUser)
	values.Set("os_password", jiraLoginPass)

	log.Println("> 1st request - login jira")
	req, err := http.NewRequestWithContext(traceCtx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "author=wangjia117")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	return client
}

func checkReleaseVersion(client *http.Client, project string, version string) bool {

	log.Println("Call checkReleaseVersion()")

	endpoint := "http://jirait.yourcompanyname.com:8080/rest/api/2/project/" + project + "/versions"

	log.Println("> 4st request - check release version create from jira api")
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	if resp.StatusCode == 401 {
		log.Fatal("Call checkReleaseVersion(): auto login with cookie failed! please retry.")
	} else if resp.StatusCode == 404 {
		log.Printf("Call checkReleaseVersion(): no project versions could be found with key: %s", project)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var Data []struct {
		Name string
	}
	json.Unmarshal([]byte(body), &Data)

	for _, vName := range Data {
		if strings.ToUpper(version) == strings.ToUpper(vName.Name) {
			log.Printf("Call checkReleaseVersions(): version:%s already created for project:%s in jira.", version, project)
			return true
		}
	}
	log.Fatalf("Call checkReleaseVersions(): version:%s not created for project:%s in jira!", version, project)
	return false
}

func getProjectID(client *http.Client, jiraLoginUser string, jiraLoginPass string, project string) int {

	log.Println("Call getProjectID()")

	endpoint := "http://jirait.yourcompanyname.com:8080/rest/api/2/project/" + project

	values := url.Values{}
	values.Set("os_username", jiraLoginUser)
	values.Set("os_password", jiraLoginPass)

	log.Println("> 2st request - get all project information from jira api")
	req, err := http.NewRequestWithContext(traceCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(jiraLoginUser, jiraLoginPass)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Cache-Control", "no-cache, no-store, no-transform")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	if resp.StatusCode == 401 {
		log.Fatal("Call getProjectID(): auth required for api get.")
	} else if resp.StatusCode == 404 {
		log.Printf("Call getProjectID(): no project could be found with key: %s", project)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var Data struct {
		Id         string
		Key        string
		Components []struct {
			Id   string
			Name string
		}
	}
	json.Unmarshal([]byte(body), &Data)

	if len(Data.Components) == 0 {
		log.Fatalf("Call getProjectID(): project key:%s not found, please check!", project)
	} else {
		pId, err := strconv.Atoi(Data.Id)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Call getProjectID(): project key:%s, project id:%d", project, pId)
		return pId
	}
	return 0
}

// getProjectIdCmd represents the getProjectId command
var getProjectIdCmd = &cobra.Command{
	Use:   "getProjectId",
	Short: "get jira project id",
	Long:  `get jira project id by project name, eg, fleet -> 10404`,
	// Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getProjectId called")

		if len(ju) == 0 {
			log.Fatal("JIRA登录用户名不能为空")
		}
		if len(jp) == 0 {
			log.Fatal("JIRA登录密码不能为空")
		}

		if len(project) == 0 {
			log.Fatal("Project名称不能为空")
		}
		flag.Parse()

		login := jiraLogin(loginURL, ju, jp)
		projectId := getProjectID(login, ju, jp, project)
		log.Println(projectId)
	},
}

func init() {

	getProjectIdCmd.Flags().StringVar(&ju, "ju", "", "jira login username")
	getProjectIdCmd.Flags().StringVar(&jp, "jp", "", "jira login password")
	getProjectIdCmd.Flags().StringVar(&project, "project", "", "jira project name")

	viper.BindPFlag("ju", getProjectIdCmd.PersistentFlags().Lookup("ju"))
	viper.BindPFlag("jp", getProjectIdCmd.PersistentFlags().Lookup("jp"))
	viper.BindPFlag("project", getProjectIdCmd.PersistentFlags().Lookup("project"))

	getProjectIdCmd.MarkFlagRequired("ju")
	getProjectIdCmd.MarkFlagRequired("jp")
	getProjectIdCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(getProjectIdCmd)

}
