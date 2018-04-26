package createticketwi

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

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
	tc.SetInput(ivDomain, "https://devjira.tibco.com")
	tc.SetInput(ivBasicAuthToken, "ZXNlY29ubmVjdG9yczozczNjMG5uZWN0b3Jz")
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
