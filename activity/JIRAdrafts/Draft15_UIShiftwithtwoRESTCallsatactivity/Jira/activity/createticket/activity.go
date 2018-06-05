package createticket

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-createTicket")

const (
	ivConnection = "Connection"
	ivProject    = "project"
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
	var domain, userName, password string
	for _, v := range connectionSettings {
		setting := v.(map[string]interface{})
		if setting["name"] == "domain" {
			domain = setting["value"].(string)
		} else if setting["name"] == "userName" {
			userName = setting["value"].(string)
		} else if setting["name"] == "password" {
			password = setting["value"].(string)
		}
	}

	fields := make(map[string]interface{})
	dynamicMap := make(map[string]interface{})
	dynamicMap["project"] = createValue("key", context.GetInput(ivProject).(string))
	dynamicMap["issuetype"] = createValue("name", context.GetInput(ivIssueType).(string))

	GetInputParameter(context.GetInput(ivInput), dynamicMap)

	fields["fields"] = dynamicMap

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
		//fmt.Printf("Response :: %s", string(jsonResponseData))

		queryResponse := make(map[string]interface{})
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		issueID := context.GetInput(ivIssueType).(string) + " " + queryResponse["key"].(string) + " is created"
		activityLog.Infof("Activity Output :: %s", issueID)
		context.SetOutput(ovIssueID, issueID)
	}

	return true, nil
}

//move below method to common go file
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func createValue(key, value string) map[string]string {
	m := make(map[string]string)
	m[key] = value
	return m
}
