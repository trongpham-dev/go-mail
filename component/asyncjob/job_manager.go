package asyncjob

import (
	"context"
	"log"
	"sync"
)

type group struct {
	isConcurrent bool
	jobs         []Job
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

//func (g *group) Run2(ctx context.Context) error {
//	errChan := make(chan error, len(g.jobs))
//
//	for i, _ := range g.jobs {
//		errChan <- g.runJob(ctx, g.jobs[i])
//	}
//
//	var err error
//
//	for i := 1; i <= len(g.jobs); i++ {
//		if v := <-errChan; v != nil {
//			err = v
//		}
//	}
//
//	return err
//}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs)) // create a wait group that waits for n jobs....

	errChan := make(chan error, len(g.jobs)) //create a buffered chanel that contains job's error

	for i := range g.jobs {
		// run jobs concurrently
		if g.isConcurrent {
			// Do this instead
			go func(currJob Job) {
				errChan <- g.runJob(ctx, currJob) // execute job then push err to errChan
				g.wg.Done()                       // complete proccess!!!
			}(g.jobs[i])

			continue
		}
		// run job in order!
		job := g.jobs[i]
		errChan <- g.runJob(ctx, job)
		g.wg.Done()
	}

	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()
	return err
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	// if an error happen
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			// cannot retry job anymore!
			if j.State() == StateRetryFailed {
				return err
			}

			// retry successfully
			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}
