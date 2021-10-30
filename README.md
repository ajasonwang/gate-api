### First Golang cli tool - CM & Gate API Intergration

#### dep install

```shell
go get -u github.com/golang/dep/cmd/dep
```

compile dep binary and copy to C:\Windows\system32\

```shell
dep init
dep ensure
```

	
##### Usage

get program version

```shell​ 
GateAPI.exe version
```

query project id by name

```shell​    
GateAPI.exe getProjectId --ju <JIRA_USER> --jp <JIRA_PASS> --uri 192.168.5.26:1180 --gu <GATE_USER> --gp <GATE_PASS> --project <JIRA_PROJECTNAME>
```

POST jenkins deploy data to Gate
​    
```shell​ 
GateAPI.exe startDeploy --ju <JIRA_USER> --jp <JIRA_PASS> --uri 192.168.5.26:1180 --gu <GATE_USER> --gp <GATE_PASS> --project <JIRA_PROJECTNAME> --versionname <JIRA_VERSION_NAME> --component <DEPLOY_ITEMS> --envname <DEPLOY_ENV>
```

#### Jenkins Intergration


Windows
```shell
Set-ExecutionPolicy Remotesigned -Force
(new-object Net.WebClient).DownloadFile('http://nexus.yourcompanyname.com:8081/nexus/content/repositories/software/cm/tools/cli/GateAPI/v2.0/GateAPI-v2.0.exe','GateAPI.exe')
```

Linux

```shell
wget http://nexus.yourcompanyname.com:8081/nexus/content/repositories/software/cm/tools/cli/GateAPI/v2.0/GateAPI-v2.0.0 -O GateAPI
chmod 755 GateAPI
```
	
```shell​ 
GateAPI.exe startDeploy --ju <JIRA_USER> --jp <JIRA_PASS> --uri 192.168.5.26:1180 --gu <GATE_USER> --gp <GATE_PASS> --project <JIRA_PROJECTNAME> --versionname <JIRA_VERSION_NAME> --envname <DEPLOY_ENV> --teams "ACE TAP" --stacks "XXXTradingDataServices XXXTradingServices" --services "cmQueryService cmQueryServiceDemo"
```
