package main

import (
	"context"
	"fmt"
	"os"
	"temporal-proofs/forloop"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func main() {
	// setup zerolog logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create client")
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "forloop_" + "0",
		TaskQueue: "forloop",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, forloop.DoSomethingsWorkflow, &forloop.DoSomethingsWorkflowRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to execute workflow")
	}
	fmt.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
