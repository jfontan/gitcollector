package gitcollector

import (
	"context"

	"gopkg.in/src-d/go-errors.v1"
)

// Job represents a gitcollector task.
type Job interface {
	// Process perform the necessary work on the job.
	Process(context.Context) error
}

// MetricsCollector represents a component in charge to collect jobs metrics.
type MetricsCollector interface {
	// Start starts collecting metrics.
	Start()
	// Stop stops collectingMetrincs.
	Stop(immediate bool)
	// Success registers metrics about successfully processed Job.
	Success(Job)
	// Faile register metrics about a failed processed Job.
	Fail(Job)
	// Discover register metrics about a discovered Job.
	Discover(Job)
}

var (
	// ErrProviderStopped is returned when a provider has been stopped.
	ErrProviderStopped = errors.NewKind("provider stopped")

	// ErrProviderStop is returned when a provider fails on Stop.
	ErrProviderStop = errors.NewKind("provider failed on stop")
)

// Provider interface represents a service to generate new Jobs.
type Provider interface {
	Start() error
	Stop() error
}
