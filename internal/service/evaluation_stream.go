package service

import (
	"context"
	"errors"
	"monitoring/internal/model"
)

func StartEvaluationStream(ctx context.Context, jobs []model.EvaluationJob) (<-chan string, <-chan error) {
	results := make(chan string, len(jobs))
	failures := make(chan error, 1)
	go func() {
		defer close(failures)
		for _, job := range jobs {
			if job.Name == "reject" {
				failures <- errors.New("source rejected")
				return
			}
			if err := model.JobContextError(ctx); err != nil {
				failures <- err
				return
			}
			results <- job.Name
		}
	}()
	return results, failures
}
