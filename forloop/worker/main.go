package main

import (
	"os"

	"temporal-proofs/forloop"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	workerOptions := worker.Options{
		EnableSessionWorker: true, // Important for a worker to participate in the session
	}
	w := worker.New(c, "forloop", workerOptions)

	w.RegisterWorkflow(forloop.DoSomethingsWorkflow)
	w.RegisterActivity(&forloop.Activities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to start worker")
	}
}
