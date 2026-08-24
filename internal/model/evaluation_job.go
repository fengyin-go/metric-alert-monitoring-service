package model

import "context"

type EvaluationJob struct{ Name string }

func JobContextError(ctx context.Context) error { return nil }
