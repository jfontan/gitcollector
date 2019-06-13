package library

import (
	"context"
	"sync"
	"testing"

	"github.com/src-d/gitcollector"
	"gopkg.in/src-d/go-log.v1"

	"github.com/stretchr/testify/require"
)

func TestJobScheduleFn(t *testing.T) {
	var (
		endpoints = []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		}

		mu        sync.Mutex
		got       []string
		processFn = func(_ context.Context, j *Job) error {
			mu.Lock()
			defer mu.Unlock()

			got = append(got, j.Endpoints[0])
			return nil
		}
	)

	download := make(chan gitcollector.Job, 2)
	update := make(chan gitcollector.Job, 20)
	sched := NewJobScheduleFn(
		nil,
		download, update,
		processFn, processFn,
		false,
		nil,
		log.New(nil),
		nil,
	)

	queues := []chan gitcollector.Job{download, update}
	expected := testScheduleFn(sched, endpoints, queues)
	require.ElementsMatch(t, expected, got)
}

func TestDownloadJobScheduleFn(t *testing.T) {
	var (
		endpoints = []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		}

		mu        sync.Mutex
		got       []string
		processFn = func(_ context.Context, j *Job) error {
			mu.Lock()
			defer mu.Unlock()

			got = append(got, j.Endpoints[0])
			return nil
		}
	)

	download := make(chan gitcollector.Job, 5)
	sched := NewDownloadJobScheduleFn(
		nil,
		download,
		processFn,
		false,
		nil,
		log.New(nil),
		nil,
	)

	queues := []chan gitcollector.Job{download}
	expected := testScheduleFn(sched, endpoints, queues)
	require.ElementsMatch(t, expected, got)
}

func TestUpdateJobScheduleFn(t *testing.T) {
	var (
		endpoints = []string{
			"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		}

		mu        sync.Mutex
		got       []string
		processFn = func(_ context.Context, j *Job) error {
			mu.Lock()
			defer mu.Unlock()

			got = append(got, j.Endpoints[0])
			return nil
		}
	)

	update := make(chan gitcollector.Job, 5)
	sched := NewUpdateJobScheduleFn(
		nil, update, processFn, nil, log.New(nil),
	)
	queues := []chan gitcollector.Job{update}
	expected := testScheduleFn(sched, endpoints, queues)
	require.ElementsMatch(t, expected, got)
}

func testScheduleFn(
	sched gitcollector.ScheduleFn,
	endpoints []string,
	queues []chan gitcollector.Job,
) []string {
	wp := gitcollector.NewWorkerPool(gitcollector.NewJobScheduler(
		sched,
		&gitcollector.JobSchedulerOpts{
			NotWaitNewJobs: true,
		},
	))

	wp.SetWorkers(10)
	wp.Run()

	for _, e := range endpoints {
		for _, queue := range queues {
			queue <- &Job{
				Endpoints: []string{e},
			}
		}
	}

	var expected []string
	for _, queue := range queues {
		expected = append(expected, endpoints...)
		close(queue)
	}

	wp.Wait()
	return expected
}
