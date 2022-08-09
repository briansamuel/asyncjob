package asyncjob

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *group) Run(ctx context.Context) error {
	errChan := make(chan error, len(g.jobs))
	g.wg.Add(len(g.jobs))
	for i, _ := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				defer Recover()
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		job := g.jobs[i]

		err := g.runJob(ctx, job)

		if err != nil {
			return err
		}

		errChan <- err
		g.wg.Done()
	}

	g.wg.Wait()

	var err error
	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			return v
		}
	}
	return err

}

func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}

func Recover() {
	if r := recover(); r != nil {
		fmt.Println("WOHA! Program is panicking with value", r)
	}
}
