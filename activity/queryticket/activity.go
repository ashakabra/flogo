package queryticket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-getrecentlyupdated")

const (
	queryByUpdate = "Recently Updated"
	queryByCreate = "Recently Created"

	ivDomain         = "domain"
	ivBasicAuthToken = "basicAuthToken"
	ivQueryBy        = "queryBy"
	ivProject        = "project"
	ivIssueType      = "issueType"
	ivWithinTime     = "withinTime"

	ovIssueIDs = "issueIDs"
)

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
	activityLog.Info("JIRA Query Get Recent Updated Activity")
	domain := context.GetInput(ivDomain).(string)
	basicAuthToken := context.GetInput(ivBasicAuthToken).(string)
	queryBy := context.GetInput(ivQueryBy).(string)
	project := context.GetInput(ivProject).(string)
	issueType := context.GetInput(ivIssueType).(string)
	withinTime := context.GetInput(ivWithinTime).(string)

	fmt.Printf("Input Values are %s, %s, %s, %s, %s, %s", domain, basicAuthToken, queryBy, project, issueType, withinTime)

	var input string
	if queryBy == queryByUpdate {
		input = "project='" + project + "' AND issueType='" + issueType + "' AND updated >= -" + withinTime
	}
	if queryBy == queryByCreate {
		input = "project='" + project + "' AND issueType='" + issueType + "' AND created >= -" + withinTime
	}

	url := domain + "/rest/api/2/search?jql=" + url.QueryEscape(input)

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+basicAuthToken)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		jsonResponseData, _ := ioutil.ReadAll(response.Body)
		//fmt.Printf(string(jsonResponseData))

		//queryResponse := make(map[string][]interface{})
		//queryResponse := make(map[string]interface{})
		var queryResponse interface{}
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		var outputStr string
		m := queryResponse.(map[string]interface{})
		issues := m["issues"].([]interface{})

		for i := range issues {
			issue := issues[i].(map[string]interface{})
			//fmt.Printf("Issue id is - %s", issue["key"])
			outputStr = outputStr + issue["key"].(string) + ", "
		}

		if len(outputStr) != 0 {
			outputStr = outputStr[0 : len(outputStr)-2] //remove last extra comma
		} else {
			outputStr = "No Issues found"
		}
		fmt.Printf("Output is  :: %s ", outputStr)
		context.SetOutput(ovIssueIDs, outputStr)
		//fmt.Printf("Issue id is - %s", issue["key"])

		//context.SetOutput(userId, dat_user["id"])
	}

	return true, nil
}
