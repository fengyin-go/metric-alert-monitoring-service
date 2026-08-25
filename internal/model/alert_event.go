package model

import (
	"strings"
	"time"
)

const (
	EventFiring   = "firing"
	EventResolved = "resolved"
)

var eventTransitions = map[string]map[string]bool{
	EventFiring:   {EventResolved: true},
	EventResolved: {},
}

// CanEventTransition 判断告警事件状态是否可流转。
func CanEventTransition(from, to string) bool {
	if m, ok := eventTransitions[from]; ok {
		return m[to]
	}
	return false
}

// AlertEvent 告警事件实体。
type AlertEvent struct {
	ID         string     `json:"id"`
	RuleID     string     `json:"rule_id"`
	RuleName   string     `json:"rule_name"`
	MetricName string     `json:"metric_name"`
	Value      float64    `json:"value"`
	Threshold  float64    `json:"threshold"`
	Severity   string     `json:"severity"`
	Status     string     `json:"status"`
	FiredAt    time.Time  `json:"fired_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Validate 校验告警事件字段。
func (e *AlertEvent) Validate() error {
	e.MetricName = strings.TrimSpace(e.MetricName)
	e.RuleName = strings.TrimSpace(e.RuleName)
	if e.MetricName == "" {
		return NewValidationError("metric_name", "指标名称不能为空")
	}
	if e.Status == "" {
		e.Status = EventFiring
	}
	if e.Status != EventFiring && e.Status != EventResolved {
		return NewValidationError("status", "事件状态不合法")
	}
	return nil
}

// AlertEventFilter 告警事件筛选条件。
type AlertEventFilter struct {
	MetricName string
	Severity   string
	Status     string
}

// Match 判断告警事件是否命中筛选条件。
func (f AlertEventFilter) Match(e *AlertEvent) bool {
	if f.MetricName != "" && e.MetricName != f.MetricName {
		return false
	}
	if f.Severity != "" && e.Severity != f.Severity {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	return true
}
