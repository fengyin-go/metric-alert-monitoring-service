package service

import (
	"sort"
	"time"

	"monitoring/internal/model"
	"monitoring/pkg/idgen"
)

// RuleInput 告警规则入参。
type RuleInput struct {
	Name       string
	MetricName string
	Operator   string
	Threshold  float64
	Severity   string
	Enabled    bool
}

// CreateAlertRule 创建告警规则。
func (s *Service) CreateAlertRule(in RuleInput) (*model.AlertRule, error) {
	now := time.Now()
	r := &model.AlertRule{
		ID:         idgen.Hex(),
		Name:       in.Name,
		MetricName: in.MetricName,
		Operator:   in.Operator,
		Threshold:  in.Threshold,
		Severity:   in.Severity,
		Enabled:    in.Enabled,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateAlertRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetAlertRule 获取规则详情。
func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

// ListAlertRules 分页查询规则。
func (s *Service) ListAlertRules(filter model.AlertRuleFilter, page, size int) ([]*model.AlertRule, int, error) {
	all := s.store.ListAlertRules()
	matched := make([]*model.AlertRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateAlertRule 更新规则。
func (s *Service) UpdateAlertRule(id string, in RuleInput) (*model.AlertRule, error) {
	existing, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = in.Name
	existing.MetricName = in.MetricName
	existing.Operator = in.Operator
	existing.Threshold = in.Threshold
	existing.Severity = in.Severity
	existing.Enabled = in.Enabled
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAlertRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAlertRule 删除规则。
func (s *Service) DeleteAlertRule(id string) error {
	return s.store.DeleteAlertRule(id)
}

// RuleStats 规则统计。
type RuleStats struct {
	Total    int `json:"total"`
	Enabled  int `json:"enabled"`
	Disabled int `json:"disabled"`
	Critical int `json:"critical"`
}

// StatsRules 统计规则数量。
func (s *Service) StatsRules() *RuleStats {
	stats := &RuleStats{}
	for _, r := range s.store.ListAlertRules() {
		stats.Total++
		if r.Enabled {
			stats.Enabled++
		} else {
			stats.Disabled++
		}
		if r.Severity == model.SeverityCritical {
			stats.Critical++
		}
	}
	return stats
}
