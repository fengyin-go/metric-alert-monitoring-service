package model

import (
	"strings"
	"time"
)

const (
	OpGT  = "gt"
	OpGTE = "gte"
	OpLT  = "lt"
	OpLTE = "lte"
	OpEQ  = "eq"
)

const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// AlertRule 告警规则实体。
type AlertRule struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MetricName string    `json:"metric_name"`
	Operator   string    `json:"operator"`
	Threshold  float64   `json:"threshold"`
	Severity   string    `json:"severity"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 校验告警规则字段。
func (r *AlertRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.MetricName = strings.TrimSpace(r.MetricName)
	r.Operator = strings.TrimSpace(r.Operator)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.MetricName == "" {
		return NewValidationError("metric_name", "指标名称不能为空")
	}
	if r.Operator == "" {
		r.Operator = OpGT
	}
	switch r.Operator {
	case OpGT, OpGTE, OpLT, OpLTE, OpEQ:
	default:
		return NewValidationError("operator", "操作符不合法")
	}
	if r.Severity == "" {
		r.Severity = SeverityWarning
	}
	switch r.Severity {
	case SeverityInfo, SeverityWarning, SeverityCritical:
	default:
		return NewValidationError("severity", "严重程度不合法")
	}
	return nil
}

// Evaluate 判断指标值是否触发告警。
func (r *AlertRule) Evaluate(value float64) bool {
	switch r.Operator {
	case OpGT:
		return value > r.Threshold
	case OpGTE:
		return value >= r.Threshold
	case OpLT:
		return value < r.Threshold
	case OpLTE:
		return value <= r.Threshold
	case OpEQ:
		return value == r.Threshold
	}
	return false
}

// AlertRuleFilter 告警规则筛选条件。
type AlertRuleFilter struct {
	Severity string
	Enabled  *bool
	Keyword  string
}

// Match 判断规则是否命中筛选条件。
func (f AlertRuleFilter) Match(r *AlertRule) bool {
	if f.Severity != "" && r.Severity != f.Severity {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.MetricName), k) {
			return false
		}
	}
	return true
}
