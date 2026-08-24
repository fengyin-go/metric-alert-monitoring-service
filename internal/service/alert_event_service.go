package service

import (
	"sort"
	"time"

	"monitoring/internal/model"
	"monitoring/pkg/idgen"
)

// EvaluateResult 规则评估结果。
type EvaluateResult struct {
	Evaluated int                 `json:"evaluated"`
	Fired     int                 `json:"fired"`
	Silenced  int                 `json:"silenced"`
	Events    []*model.AlertEvent `json:"events"`
}

// EvaluateRules 评估全部启用规则，触发告警事件（静默期跳过）。
func (s *Service) EvaluateRules() (*EvaluateResult, error) {
	result := &EvaluateResult{Events: make([]*model.AlertEvent, 0)}
	now := time.Now()
	for _, rule := range s.store.ListAlertRules() {
		if !rule.Enabled {
			continue
		}
		metric, err := s.store.GetMetricByName(rule.MetricName)
		if err != nil {
			continue
		}
		result.Evaluated++
		if s.isSilenced(rule.MetricName, now) {
			result.Silenced++
			continue
		}
		if rule.Evaluate(metric.Value) {
			event := s.fireEvent(rule, metric.Value, now)
			if event != nil {
				result.Events = append(result.Events, event)
				result.Fired++
			}
		}
	}
	return result, nil
}

// isSilenced 判断指标在给定时刻是否处于静默期。
func (s *Service) isSilenced(metricName string, now time.Time) bool {
	for _, si := range s.store.ListSilences() {
		if si.MetricName == metricName && si.Active(now) {
			return true
		}
	}
	return false
}

// fireEvent 创建一条 firing 告警事件。
func (s *Service) fireEvent(rule *model.AlertRule, value float64, now time.Time) *model.AlertEvent {
	event := &model.AlertEvent{
		ID:         idgen.Hex(),
		RuleID:     rule.ID,
		RuleName:   rule.Name,
		MetricName: rule.MetricName,
		Value:      value,
		Threshold:  rule.Threshold,
		Severity:   rule.Severity,
		Status:     model.EventFiring,
		FiredAt:    now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := event.Validate(); err != nil {
		return nil
	}
	if err := s.store.CreateAlertEvent(event); err != nil {
		return nil
	}
	return event
}

// GetAlertEvent 获取事件详情。
func (s *Service) GetAlertEvent(id string) (*model.AlertEvent, error) {
	return s.store.GetAlertEvent(id)
}

// ListAlertEvents 分页查询事件。
func (s *Service) ListAlertEvents(filter model.AlertEventFilter, page, size int) ([]*model.AlertEvent, int, error) {
	all := s.store.ListAlertEvents()
	matched := make([]*model.AlertEvent, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].FiredAt.After(matched[j].FiredAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertEvent{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// ResolveEvent 解决告警事件（firing → resolved）。
func (s *Service) ResolveEvent(id string) (*model.AlertEvent, error) {
	existing, err := s.store.GetAlertEvent(id)
	if err != nil {
		return nil, err
	}
	if !model.CanEventTransition(existing.Status, model.EventResolved) {
		return nil, model.NewValidationError("status", "告警事件状态不允许流转到 resolved")
	}
	now := time.Now()
	existing.Status = model.EventResolved
	existing.ResolvedAt = &now
	existing.UpdatedAt = now
	if err := s.store.UpdateAlertEvent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// EventStats 告警事件统计。
type EventStats struct {
	Total    int `json:"total"`
	Firing   int `json:"firing"`
	Resolved int `json:"resolved"`
	Critical int `json:"critical"`
}

// StatsEvents 全局统计告警事件。
func (s *Service) StatsEvents() *EventStats {
	stats := &EventStats{}
	for _, e := range s.store.ListAlertEvents() {
		stats.Total++
		if e.Status == model.EventFiring {
			stats.Firing++
		} else {
			stats.Resolved++
		}
		if e.Severity == model.SeverityCritical {
			stats.Critical++
		}
	}
	return stats
}

// ActiveEvents 返回当前处于 firing 状态的告警事件。
func (s *Service) ActiveEvents() ([]*model.AlertEvent, error) {
	active := make([]*model.AlertEvent, 0)
	for _, e := range s.store.ListAlertEvents() {
		if e.Status == model.EventFiring {
			active = append(active, e)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].FiredAt.After(active[j].FiredAt)
	})
	return active, nil
}

// CountEventsByMetric 统计指定指标的事件数量。
func (s *Service) CountEventsByMetric(metricName string) (int, error) {
	count := 0
	for _, e := range s.store.ListAlertEvents() {
		if e.MetricName == metricName {
			count++
		}
	}
	return count, nil
}

// EventsForRule 返回指定规则触发的告警事件。
func (s *Service) EventsForRule(ruleID string) ([]*model.AlertEvent, error) {
	if _, err := s.store.GetAlertRule(ruleID); err != nil {
		return nil, err
	}
	events := make([]*model.AlertEvent, 0)
	for _, e := range s.store.ListAlertEvents() {
		if e.RuleID == ruleID {
			events = append(events, e)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].FiredAt.After(events[j].FiredAt)
	})
	return events, nil
}
