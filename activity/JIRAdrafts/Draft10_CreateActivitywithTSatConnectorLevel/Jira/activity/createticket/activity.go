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
	issueTask  = "Task"
	issueStory = "Story"

	severityCritical = "Critical"
	severityHigh     = "High"
	severityLow      = "Low"

	ivConnection    = "Connection"
	ivProject       = "project"
	ivIssueType     = "issueType"
	ivSummary       = "summary"
	ivDescription   = "description"
	ivAffectVersion = "affectVersion"
	ivConfirmer     = "confirmer"
	ivSeverity      = "severity"

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
		Versions  *Versions  `json:"versions,omitempty"`
		Severity  *Severity  `json:"customfield_10024,omitempty"`
		Confirmer *Confirmer `json:"customfield_10000,omitempty"`
	} `json:"fields"`
}

type Versions [1]struct { //for now I am using text field and not multiple values so have set size 1
	Name string `json:"name,omitempty"`
}

type Severity struct {
	Id string `json:"id,omitempty"`
}

type Confirmer struct {
	Name string `json:"name,omitempty"`
}

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
	issue := &Issue{}
	versions := &Versions{}
	severity := &Severity{}
	confirmer := &Confirmer{}

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
		} else if setting["name"] == "keyvalue" {
			keyvalue := setting["value"].(map[string]interface{})
			activityLog.Infof("From Go , Value from keyvalue Defect is :: %s", keyvalue["Defect"])
		}
	}

	activityLog.Infof("Connection Details is -- domain : %s, username : %s", domain, userName)

	issue.Fields.Project.Key = context.GetInput(ivProject).(string)
	issue.Fields.Summary = context.GetInput(ivSummary).(string)
	issue.Fields.Description = context.GetInput(ivDescription).(string)
	issue.Fields.IssueType.Name = context.GetInput(ivIssueType).(string)

	//Confirmer is not allowed in case of Story, Task
	//Confirmer is required in case of Enhancement, Defect
	if issue.Fields.IssueType.Name != issueStory && issue.Fields.IssueType.Name != issueTask {
		confirmer.Name = context.GetInput(ivConfirmer).(string)
		issue.Fields.Confirmer = confirmer
	}

	//Versions,Severity is not allowed in case of Task
	//Versions,Severity is required in case of Enhancement, Defect & Optional in case of Story(but it is allowed)
	if issue.Fields.IssueType.Name != issueTask {
		versions[0].Name = context.GetInput(ivAffectVersion).(string)
		issue.Fields.Versions = versions

		//Allowed values are: 10036[1-Critical], 10037[2-High], 10038[3-Low]
		if context.GetInput(ivSeverity).(string) == severityCritical {
			severity.Id = "10036"
		} else if context.GetInput(ivSeverity).(string) == severityHigh {
			severity.Id = "10037"
		} else if context.GetInput(ivSeverity).(string) == severityLow {
			severity.Id = "10038"
		}
		issue.Fields.Severity = severity
	}

	//fmt.Printf("Input Values are %s, %s, %s, %s, %s, %s", domain, basicAuthToken, issue.Fields.Project.Key, issue.Fields.Summary, issue.Fields.Description, issue.Fields.IssueType.Name)

	jsonData, err := json.Marshal(issue)
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

		issueID := issue.Fields.IssueType.Name + " " + queryResponse["key"].(string) + " is created"
		fmt.Printf("Activity Output :: %s", issueID)
		context.SetOutput(ovIssueID, issueID)
	}

	return true, nil
}

//move below code to common go file
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
