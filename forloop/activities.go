package forloop

import (
	"context"

	"github.com/gogo/status"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/activity"
	"google.golang.org/grpc/codes"
)

/**
 * Sample activities used by file processing sample workflow.
 */
var errorCount int = 0

type Activities struct {
}
type DoSomethingRequest struct {
	Name string
}
type DoSomethingResponse struct {
	Name string
}

func (a *Activities) DoSomethingActivity(ctx context.Context, request *DoSomethingRequest) (*DoSomethingResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("DoSomethingActivity...", "Name", request.Name)
	log.Info().Interface("request", request).Msg("request")
	if request.Name == "error" {
		errorCount++
		if errorCount < 3 {
			log.Error().Interface("request", request).Msg("request")
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}
	}
	return &DoSomethingResponse{Name: request.Name}, nil
}
