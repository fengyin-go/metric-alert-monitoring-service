package model

import (
	"strings"
	"time"
)

// Silence 告警静默期实体。
type Silence struct {
	ID         string    `json:"id"`
	MetricName string    `json:"metric_name"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
	Reason     string    `json:"reason"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 校验静默期字段。
func (s *Silence) Validate() error {
	s.MetricName = strings.TrimSpace(s.MetricName)
	s.Reason = strings.TrimSpace(s.Reason)
	if s.MetricName == "" {
		return NewValidationError("metric_name", "指标名称不能为空")
	}
	if !s.EndAt.After(s.StartAt) {
		return NewValidationError("end_at", "结束时间必须晚于开始时间")
	}
	return nil
}

// Active 判断静默期在给定时刻是否生效。
func (s *Silence) Active(now time.Time) bool {
	return s.Enabled && !now.Before(s.StartAt) && now.Before(s.EndAt)
}

// SilenceFilter 静默期筛选条件。
type SilenceFilter struct {
	MetricName string
}

// Match 判断静默期是否命中筛选条件。
func (f SilenceFilter) Match(s *Silence) bool {
	if f.MetricName != "" && s.MetricName != f.MetricName {
		return false
	}
	return true
}
