package multiply

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// activityLog is the default logger for the Multiply Activity
var activityLog = logger.GetLogger("activity-Asha-multiply")

const (
	ivNum1 = "number1"
	ivNum2 = "number2"

	ovMultiply = "multiply"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

// MultiplyActivity is an Activity that is used to log a message to the console
// inputs : {number1, number2}
// outputs: multiply
type MultiplyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MultiplyActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MultiplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *MultiplyActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	num1, _ := context.GetInput(ivNum1).(int)
	num2, _ := context.GetInput(ivNum2).(int)

	activityLog.Info(fmt.Sprintf("Number1 is %d and Number2 is %d", num1, num2))

	context.SetOutput(ovMultiply, num1*num2)

	return true, nil
}
