//ASK --error codes? activity.NewError("Zuora connection is not configured", "ZUORA-SUBSCRIBER-4001", nil)

package queryticket

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var activityLog = logger.GetLogger("jira-activity-queryticket")

const (
	queryByUpdate = "Recently Updated"
	queryByCreate = "Recently Created"

	ivConnection  = "Connection"
	ivQueryBy     = "queryBy"
	ivProject     = "project"
	ivIssueType   = "issueType"
	ivWithinTime  = "withinTime"
	ivQueryParams = "queryParams"

	ovOutput = "output"
)

type QueryTicketActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &QueryTicketActivity{metadata: metadata}
}

func (a *QueryTicketActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func ParseOutput(outputSchema map[string]interface{}) ([]string, error) {
	if outputSchema == nil {
		return nil, nil
	}

	props := outputSchema["items"].(map[string]interface{})
	properties := props["properties"].(map[string]interface{})
	fields := make([]string, len(properties))
	i := 0
	for k, _ := range properties {
		fields[i] = k
		i++
	}
	return fields, nil
}

func (a *QueryTicketActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("JIRA Query Ticket")

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

	activityLog.Infof("Connection Details is -- domain : %s, username : %s", domain, userName)

	queryBy := context.GetInput(ivQueryBy).(string)
	project := context.GetInput(ivProject).(string)
	issueType := context.GetInput(ivIssueType).(string)
	withinTime := context.GetInput(ivWithinTime).(string)
	parameters, err := GetParameter(context.GetInput(ivQueryParams))

	//extra code to read key names of output
	outputMap, _ := LoadJsonSchemaFromMetadata(context.GetOutput(ovOutput))
	if outputMap != nil {
		outputFields, _ := ParseOutput(outputMap)
		activityLog.Infof("Reading Output is :: %s", outputFields)
	}
	//end of extra code

	fmt.Printf("Input Values are %s, %s, %s, %s", queryBy, project, issueType, withinTime)

	var input string
	if queryBy == queryByUpdate {
		//input = "project='" + project + "' AND issueType='" + issueType + "' AND updated >= -" + withinTime
		input = "issueType='" + issueType + "' AND updated >= -" + withinTime
	} else if queryBy == queryByCreate {
		//input = "project='" + project + "' AND issueType='" + issueType + "' AND created >= -" + withinTime
		input = "issueType='" + issueType + "' AND created >= -" + withinTime
	}

	for _, value := range parameters.QueryParams {
		activityLog.Infof("Name :: %s , Value :: %s ", value.Name, value.Value)
		stringValue := fmt.Sprint(value.Value)
		input = input + " AND " + value.Name + " in (" + stringValue + ")"
	}

	activityLog.Infof("Input Value :: %s", input)
	url := domain + "/rest/api/2/search?jql=" + url.QueryEscape(input)

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+BasicAuth(userName, password))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		jsonResponseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			return false, fmt.Errorf("Error reading JSON response data after query invocation, %s", err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			//activityLog.Infof("Jira Rest API received HTTP status: %d  detailed reason:[%s]", response.StatusCode, jsonResponseData)
			return false, fmt.Errorf("Jira Rest API received HTTP status: %d  detailed reason:[%s]", response.StatusCode, jsonResponseData)
		}

		var queryResponse interface{}
		err = json.Unmarshal(jsonResponseData, &queryResponse)

		if err != nil {
			fmt.Printf("Error :: %s", err)
		}

		m := queryResponse.(map[string]interface{})
		issues := m["issues"].([]interface{})
		responseIssues := make([]map[string]interface{}, len(issues))
		if len(issues) > 0 {
			for i := range issues {
				issue := issues[i].(map[string]interface{})
				responseIssues[i] = make(map[string]interface{})
				responseIssues[i]["key"] = issue["key"].(string)
				responseIssues[i]["summary"] = issue["fields"].(map[string]interface{})["summary"]
			}
		} else {
			activityLog.Infof("No issues found")
		}

		//just extra logs for test, remove this later
		activityLog.Infof("Number of issues in Output Map -- %d", len(issues))
		activityLog.Infof("Output is -- %s", responseIssues)

		output := &data.ComplexObject{Metadata: "", Value: responseIssues}
		context.SetOutput(ovOutput, output)
	}

	return true, nil
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
