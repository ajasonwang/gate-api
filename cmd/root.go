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
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"net"
	"net/http"

	"net/http/httptrace"
	"os"

	"time"

	"github.com/spf13/cobra"

	"net/url"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var ju string
var jp string
var uri string
var gu string
var gp string
var project string
var versionname string
var envname string
var teams string
var stacks string
var services string
var deployid string

var loginURL = "http://jirait.yourcompanyname.com:8080/login.jsp"

func readFile(file string) string {
	fileData, err := ioutil.ReadFile(file)
	if err == nil {
		log.Println("File content is:", string(fileData))
	} else {
		log.Fatal("ReadFile error:", err)
	}

	return string(fileData)
}

func writeFile(file string, content string) {
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("WriteFile failure, err=[%v]\n", err)
	} else {
		fmt.Printf("WriteFile %s success.\n", file)
	}
}

func traceConnection() context.Context {
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf(">>> GotConn - Connection was reused: %t", info.Reused)
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Printf(">>> DNSStart - Host: %s", info.Host)
		},
		ConnectStart: func(network, addr string) {
			log.Printf(">>> ConnectStart - Network type: %s, Address: %s", network, addr)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf(">>> DNSDone - DNS Info: %+v\n", dnsInfo)
		},
	}
	ctx := httptrace.WithClientTrace(context.Background(), clientTrace)
	return ctx
}

var traceCtx = traceConnection()

var numCoroutines = http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func httpClient() *http.Client {

	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)

	client := &http.Client{
		Transport: &http.Transport{
			// Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost:   numCoroutines,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       10,
		},
		Timeout: 60 * time.Second,
		Jar:     jar,
	}

	return client
}

var client = httpClient()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GateAPI",
	Short: "CM & QA Gate API intergration",
	Long:  `CM & QA Gate API intergration`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.GateAPI.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".GateAPI" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".GateAPI")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
