package model

type MetricBatch struct {
	Name   string
	Values []float64
}

func NewMetricBatch(name string, values []float64) MetricBatch {
	return MetricBatch{Name: name, Values: values}
}
