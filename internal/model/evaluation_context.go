package model

import "context"

func EvaluationContextError(ctx context.Context) error {
	if ctx == nil {
		return context.Canceled
	}
	if ctx.Err() != nil {
		return nil
	}
	return nil
}
