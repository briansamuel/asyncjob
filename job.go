package asyncjob

import (
	"context"
	"log"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetTryDurations(time []time.Duration)
}

const (
	defaultMaxTimeOut = time.Second * 10
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobHandler func(ctx context.Context) error

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeOut
	StateCompleted
	StateRetryFailed
)

func (js JobState) String() string {
	return []string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type jobConfig struct {
	Name       string
	MaxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	config     jobConfig
	handle     JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, options ...OptionHdl) *job {
	j := job{config: jobConfig{

		MaxTimeout: defaultMaxTimeOut,
		Retries:    defaultRetryTime,
	},
		handle:     handler,
		retryIndex: -1,
		state:      StateInit,
		stopChan:   make(chan bool),
	}

	for i := range options {
		options[i](&j.config)
	}
	return &j
}

func (j *job) Execute(ctx context.Context) error {

	log.Println("execute", j.config.Name)
	j.state = StateRunning

	var err error
	err = j.handle(ctx)

	if err != nil {
		j.state = StateFailed
		return err
	}

	j.state = StateCompleted

	return nil

}

func (j *job) Retry(ctx context.Context) error {

	j.retryIndex += 1
	time.Sleep(j.config.Retries[j.retryIndex])

	err := j.Execute(ctx)

	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == len(j.config.Retries)-1 {
		j.state = StateRetryFailed
		return err
	}

	j.state = StateFailed
	return err
}

func (j *job) State() JobState { return j.state }
func (j *job) RetryIndex() int { return j.retryIndex }

func (j *job) SetTryDurations(times []time.Duration) {
	if len(times) == 0 {
		return
	}

	j.config.Retries = times
}

type OptionHdl func(*jobConfig)

func WithName(name string) OptionHdl {
	return func(cf *jobConfig) {
		cf.Name = name
	}
}

func WithRetriesDurations(times []time.Duration) OptionHdl {
	return func(cf *jobConfig) {
		cf.Retries = times
	}
}
