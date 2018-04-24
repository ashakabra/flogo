package createticket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-getrecentlyupdated")

const (
	ivDomain         = "domain"
	ivBasicAuthToken = "basicAuthToken"
	ivProject        = "project"
	ivIssueType      = "issueType"
	ivSummary        = "summary"
	ivDescription    = "description"

	ovIssueID = "issueID"
)

type Issue struct {
	Fields struct {
		Project struct {
			Key string `json:"key"`
		} `json:"project"`
		Summary     string `json:"summary"`
		Description string `json:"description"`
		IssueType   struct {
			Name string `json:"name"`
		} `json:"issuetype"`
	} `json:"fields"`
}

type GetUpdatedIssueActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &GetUpdatedIssueActivity{metadata: metadata}
}

func (a *GetUpdatedIssueActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *GetUpdatedIssueActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("JIRA Create Ticket")
	issue := &Issue{}
	domain := context.GetInput(ivDomain).(string)
	basicAuthToken := context.GetInput(ivBasicAuthToken).(string)
	issue.Fields.Project.Key = context.GetInput(ivProject).(string)
	issue.Fields.Summary = context.GetInput(ivSummary).(string)
	issue.Fields.Description = context.GetInput(ivDescription).(string)
	issue.Fields.IssueType.Name = context.GetInput(ivIssueType).(string)

	fmt.Printf("Input Values are %s, %s, %s, %s, %s, %s", domain, basicAuthToken, issue.Fields.Project.Key, issue.Fields.Summary, issue.Fields.Description, issue.Fields.IssueType.Name)

	jsonData, err := json.Marshal(issue)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	url := domain + "/rest/api/2/issue/"

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+basicAuthToken)

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

		issueID := issue.Fields.IssueType.Name + " " + queryResponse["key"].(string) + " is created"
		fmt.Printf("Activity Output :: %s", issueID)
		context.SetOutput(ovIssueID, issueID)
	}

	return true, nil
}
