package service

import (
	"sort"
	"time"

	"monitoring/internal/model"
	"monitoring/pkg/idgen"
)

// CreateSilence 创建告警静默期。
func (s *Service) CreateSilence(metricName string, startAt, endAt time.Time, reason string, enabled bool) (*model.Silence, error) {
	si := &model.Silence{
		ID:         idgen.Hex(),
		MetricName: metricName,
		StartAt:    startAt,
		EndAt:      endAt,
		Reason:     reason,
		Enabled:    enabled,
		CreatedAt:  time.Now(),
	}
	if err := si.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSilence(si); err != nil {
		return nil, err
	}
	return si, nil
}

// GetSilence 获取静默期详情。
func (s *Service) GetSilence(id string) (*model.Silence, error) {
	return s.store.GetSilence(id)
}

// ListSilences 分页查询静默期。
func (s *Service) ListSilences(filter model.SilenceFilter, page, size int) ([]*model.Silence, int, error) {
	all := s.store.ListSilences()
	matched := make([]*model.Silence, 0, len(all))
	for _, si := range all {
		if filter.Match(si) {
			matched = append(matched, si)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].StartAt.After(matched[j].StartAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Silence{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteSilence 删除静默期。
func (s *Service) DeleteSilence(id string) error {
	return s.store.DeleteSilence(id)
}

// ActiveSilences 返回当前生效的静默期。
func (s *Service) ActiveSilences() ([]*model.Silence, error) {
	now := time.Now()
	active := make([]*model.Silence, 0)
	for _, si := range s.store.ListSilences() {
		if si.Active(now) {
			active = append(active, si)
		}
	}
	return active, nil
}

// SilenceStats 静默期统计。
type SilenceStats struct {
	Total  int `json:"total"`
	Active int `json:"active"`
}

// StatsSilences 统计静默期数量。
func (s *Service) StatsSilences() *SilenceStats {
	now := time.Now()
	stats := &SilenceStats{}
	for _, si := range s.store.ListSilences() {
		stats.Total++
		if si.Active(now) {
			stats.Active++
		}
	}
	return stats
}

// DeleteSilencesByMetric 删除指定指标的全部静默期，返回删除数量。
func (s *Service) DeleteSilencesByMetric(metricName string) (int, error) {
	count := 0
	for _, si := range s.store.ListSilences() {
		if si.MetricName == metricName {
			if err := s.store.DeleteSilence(si.ID); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
