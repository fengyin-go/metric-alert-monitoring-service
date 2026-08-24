package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
)

func SubmitMetricBatch(queue *service.MetricBatchQueue, name string, values []float64) []float64 {
	requestValues := values
	queue.Submit(model.NewMetricBatch(name, requestValues))
	return requestValues
}
