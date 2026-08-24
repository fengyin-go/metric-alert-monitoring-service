package service

import (
	"sort"
	"time"

	"monitoring/internal/model"
	"monitoring/pkg/idgen"
)

// UpsertMetric 上报指标，存在则更新值，否则创建。
func (s *Service) UpsertMetric(name, typ, unit string, value float64) (*model.Metric, error) {
	now := time.Now()
	if existing, err := s.store.GetMetricByName(name); err == nil {
		existing.Value = value
		if unit != "" {
			existing.Unit = unit
		}
		existing.UpdatedAt = now
		if err := s.store.UpdateMetric(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}
	m := &model.Metric{
		ID:        idgen.Hex(),
		Name:      name,
		Type:      typ,
		Unit:      unit,
		Value:     value,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateMetric(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GetMetric 获取指标详情。
func (s *Service) GetMetric(id string) (*model.Metric, error) {
	return s.store.GetMetric(id)
}

// GetMetricByName 按名称获取指标。
func (s *Service) GetMetricByName(name string) (*model.Metric, error) {
	return s.store.GetMetricByName(name)
}

// ListMetrics 分页查询指标。
func (s *Service) ListMetrics(filter model.MetricFilter, page, size int) ([]*model.Metric, int, error) {
	all := s.store.ListMetrics()
	matched := make([]*model.Metric, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Name < matched[j].Name
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Metric{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteMetric 删除指标。
func (s *Service) DeleteMetric(id string) error {
	return s.store.DeleteMetric(id)
}

// MetricStats 指标统计。
type MetricStats struct {
	Total   int `json:"total"`
	Gauge   int `json:"gauge"`
	Counter int `json:"counter"`
}

// StatsMetrics 按类型统计指标数量。
func (s *Service) StatsMetrics() *MetricStats {
	stats := &MetricStats{}
	for _, m := range s.store.ListMetrics() {
		stats.Total++
		if m.Type == model.MetricGauge {
			stats.Gauge++
		} else {
			stats.Counter++
		}
	}
	return stats
}

// ListMetricRules 返回关联指定指标的告警规则。
func (s *Service) ListMetricRules(metricName string) ([]*model.AlertRule, error) {
	rules := make([]*model.AlertRule, 0)
	for _, r := range s.store.ListAlertRules() {
		if r.MetricName == metricName {
			rules = append(rules, r)
		}
	}
	return rules, nil
}
