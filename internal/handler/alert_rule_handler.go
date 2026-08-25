package handler

import (
	"net/http"

	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/pkg/httpx"
)

func (s *Server) registerAlertRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createAlertRule)
	mux.HandleFunc("GET /api/rules", s.listAlertRules)
	mux.HandleFunc("GET /api/rules/stats", s.statsRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getAlertRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateAlertRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteAlertRule)
}

type ruleRequest struct {
	Name       string  `json:"name"`
	MetricName string  `json:"metric_name"`
	Operator   string  `json:"operator"`
	Threshold  float64 `json:"threshold"`
	Severity   string  `json:"severity"`
	Enabled    *bool   `json:"enabled"`
}

func (s *Server) createAlertRule(w http.ResponseWriter, r *http.Request) {
	var req ruleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateAlertRule(service.RuleInput{
		Name:       req.Name,
		MetricName: req.MetricName,
		Operator:   req.Operator,
		Threshold:  req.Threshold,
		Severity:   req.Severity,
		Enabled:    boolValue(req.Enabled),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listAlertRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertRuleFilter{
		Severity: r.URL.Query().Get("severity"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListAlertRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.svc.GetAlertRule(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) updateAlertRule(w http.ResponseWriter, r *http.Request) {
	var req ruleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateAlertRule(r.PathValue("id"), service.RuleInput{
		Name:       req.Name,
		MetricName: req.MetricName,
		Operator:   req.Operator,
		Threshold:  req.Threshold,
		Severity:   req.Severity,
		Enabled:    boolValue(req.Enabled),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteAlertRule(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteAlertRule(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func boolValue(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

func (s *Server) statsRules(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsRules())
}
