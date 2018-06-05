package createticket

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

var connectionJSON = []byte(`{
	"id" : "JiraTestConnection",
	"name": "jiraconnection",
	"description" : "JIRA Test Connection",
	"title": "JIRA Connector",
	"type": "flogo:connector",
	"version": "1.0.0",
	"ref": "Jira/connector/jiraconnection",
	"keyfield": "name",
	"settings": [
		{
		  "name": "name",
		  "value": "MyTestConnection",
		  "type": "string"
		},
		{
		  "name": "description",
		  "value": "My JIRA test Connection",
		  "type": "string"
		},
		{
		  "name": "domain",
		  "value": "https://devjira.tibco.com",
		  "type": "string"
		  
		},
		{
		  "name": "userName",
		  "value": "eseconnectors",
		  "type": "string"
		  
		},  
		{
		  "name": "password",
		  "value": "3s3c0nnectors",
		  "type": "string"
		  
		}
	  ]
}`)

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	conn := make(map[string]interface{})
	err := json.Unmarshal([]byte(connectionJSON), &conn)
	assert.Nil(t, err)

	tc.SetInput(ivConnection, conn)
	tc.SetInput(ivProject, "ESEC")
	tc.SetInput(ivSummary, "From GO TEST")
	tc.SetInput(ivDescription, "Enhancement is created from GO TEST")
	tc.SetInput(ivIssueType, "Enhancement")
	tc.SetInput(ivAffectVersion, "1.0.0")
	tc.SetInput(ivSeverity, "Low")
	tc.SetInput(ivConfirmer, "akabra")
	act.Eval(tc)

	//name := tc.GetOutput(name).(string)

	//if name != "Test" {
	//	t.Error("Name did not match")
	//}
}
