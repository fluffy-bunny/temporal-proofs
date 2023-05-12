package forloop

import (
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/temporal"

	"go.temporal.io/sdk/workflow"
)

var names = [][]string{
	{"first", "second", "error"},
	{"first", "second", "error", "third"},
	{"first", "second", "error", "third", "fourth"},
	{"first", "second", "error", "third", "fourth", "fifth"},
}

type DoSomethingsWorkflowRequest struct {
}

func GetNames() []string {
	namesIdx := time.Now().Unix() % int64(len(names))
	log.Info().Int64("namesIdx", namesIdx).Msg("namesIdx")
	return names[namesIdx]
}

// DoSomethingsWorkflow workflow definition
func DoSomethingsWorkflow(ctx workflow.Context, request *DoSomethingsWorkflowRequest) (err error) {
	errorCount = 0
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		HeartbeatTimeout:    2 * time.Second, // such a short timeout to make sample fail over very fast
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var a *Activities
	// Retry the whole sequence from the first activity on any error
	// to retry it on a different host. In a real application it might be reasonable to
	// retry individual activities as well as the whole sequence discriminating between different types of errors.
	// See the retryactivity sample for a more sophisticated retry implementation.
	names := GetNames()
	log.Info().Interface("names", names).Msg("names")
	for i := 0; i < len(names); i++ {
		log.Info().Int("i", i).Msg(names[i])
		var response *DoSomethingResponse
		err = workflow.ExecuteActivity(ctx, a.DoSomethingActivity, &DoSomethingRequest{
			Name: names[i],
		}).Get(ctx, &response)
		log.Info().Interface("response", response).Msg("response")
	}

	if err != nil {
		workflow.GetLogger(ctx).Error("Workflow failed.", "Error", err.Error())
	} else {
		workflow.GetLogger(ctx).Info("Workflow completed.")
	}
	return err

}
