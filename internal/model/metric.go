package model

import (
	"strings"
	"time"
)

const (
	MetricGauge   = "gauge"
	MetricCounter = "counter"
)

// Metric 监控指标实体。
type Metric struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Unit      string    `json:"unit"`
	Value     float64   `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验指标字段。
func (m *Metric) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	m.Type = strings.TrimSpace(m.Type)
	m.Unit = strings.TrimSpace(m.Unit)
	if m.Name == "" {
		return NewValidationError("name", "指标名称不能为空")
	}
	if m.Type == "" {
		m.Type = MetricGauge
	}
	if m.Type != MetricGauge && m.Type != MetricCounter {
		return NewValidationError("type", "指标类型不合法")
	}
	return nil
}

// MetricFilter 指标筛选条件。
type MetricFilter struct {
	Type    string
	Keyword string
}

// Match 判断指标是否命中筛选条件。
func (f MetricFilter) Match(m *Metric) bool {
	if f.Type != "" && m.Type != f.Type {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Name), k) {
			return false
		}
	}
	return true
}
