package toUpperCaseStr

import (
	"fmt"
	"strings"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the Add Numbers Activity
var activityLog = logger.GetLogger("activity-Asha-toUpperCaseStr")

const (
	ivInputStr  = "InputStr"
	ovToUpperStr = "OutputToUpperStr"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

// LogActivity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type ToUpperActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ToUpperActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ToUpperActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *ToUpperActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	inputStr, _ := context.GetInput(ivInputStr).(string)
	
	activityLog.Info(fmt.Sprintf("Input String before Reverse is :: %s", inputStr))
	activityLog.Info(fmt.Sprintf("Reverse String is : %s", strings.ToUpper(inputStr)))
	context.SetOutput(ovToUpperStr, strings.ToUpper(inputStr))
	
	return true, nil
}
