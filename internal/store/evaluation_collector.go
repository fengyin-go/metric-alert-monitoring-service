package store

import "context"

func CollectEvaluations(ctx context.Context, results <-chan string, failures <-chan error) ([]string, error) {
	var collected []string
	for results != nil || failures != nil {
		select {
		case value, ok := <-results:
			if !ok {
				results = nil
			} else {
				collected = append(collected, value)
			}
		case _, ok := <-failures:
			if !ok {
				failures = nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return collected, nil
}
