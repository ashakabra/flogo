package createticket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-createTicket")

const (
	ivConnection = "Connection"
	ivIssueType  = "issueType"
	ivInput      = "input"

	ovIssueID = "issueID"
)

type CreateTicketActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CreateTicketActivity{metadata: metadata}
}

func (a *CreateTicketActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *CreateTicketActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("JIRA Create Ticket")
	//Read Inputs
	if context.GetInput(ivConnection) == nil || len(context.GetInput(ivConnection).(map[string]interface{})) == 0 {
		return false, fmt.Errorf("Jira connection is not configured")
	}

	//Read connection details
	connectionInfo := context.GetInput(ivConnection).(map[string]interface{})
	connectionSettings := connectionInfo["settings"].([]interface{})
	var domain, userName, password, project string
	for _, v := range connectionSettings {
		setting := v.(map[string]interface{})
		if setting["name"] == "domain" {
			domain = setting["value"].(string)
		} else if setting["name"] == "userName" {
			userName = setting["value"].(string)
		} else if setting["name"] == "password" {
			password = setting["value"].(string)
		} else if setting["name"] == "project" {
			project = setting["value"].(string)
		}
	}

	//issue1 := context.GetInput(ivInput).(*data.ComplexObject).Metadata

	activityLog.Infof("Project name is :: %s", project)

	fields := make(map[string]interface{})
	dynamicMap := make(map[string]interface{})
	dynamicMap["project"] = createValue("key", project)
	dynamicMap["issuetype"] = createValue("name", context.GetInput(ivIssueType).(string))

	issue := context.GetInput(ivInput).(*data.ComplexObject).Value
	for k, v := range issue.(map[string]interface{}) {
		dynamicMap[k] = v
	}
	//dynamicMap["summary"] = "Asha summary"
	fields["fields"] = dynamicMap

	/*jsonData1, err1 := json.Marshal(issue1)
	if err1 != nil {
		fmt.Printf("Error: %s", err1)
		return
	}
	fmt.Printf("Metadata IS :: %s", jsonData1)*/

	jsonData, err := json.Marshal(fields)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Printf("JSON DATA IS :: %s", jsonData)
	url := domain + "/rest/api/2/issue/"

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+BasicAuth(userName, password))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		jsonResponseData, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Response :: %s", string(jsonResponseData))

		queryResponse := make(map[string]interface{})
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		//issueID := issue.Fields.IssueType.Name + " " + queryResponse["key"].(string) + " is created"
		//fmt.Printf("Activity Output :: %s", issueID)
		//context.SetOutput(ovIssueID, issueID)
	}

	return true, nil
}

//move below code to common go file
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func createValue(key, value string) map[string]string {
	m := make(map[string]string)
	m[key] = value
	return m
}
