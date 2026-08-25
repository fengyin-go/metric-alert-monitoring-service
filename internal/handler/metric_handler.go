package handler

import (
	"net/http"

	"monitoring/internal/model"
	"monitoring/pkg/httpx"
)

func (s *Server) registerMetricRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/metrics/upsert", s.upsertMetric)
	mux.HandleFunc("GET /api/metrics", s.listMetrics)
	mux.HandleFunc("GET /api/metrics/stats", s.statsMetrics)
	mux.HandleFunc("GET /api/metrics/rules", s.listMetricRules)
	mux.HandleFunc("GET /api/metrics/{id}", s.getMetric)
	mux.HandleFunc("DELETE /api/metrics/{id}", s.deleteMetric)
}

type upsertMetricRequest struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func (s *Server) upsertMetric(w http.ResponseWriter, r *http.Request) {
	var req upsertMetricRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpsertMetric(req.Name, req.Type, req.Unit, req.Value)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) listMetrics(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MetricFilter{
		Type:    r.URL.Query().Get("type"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMetrics(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMetric(w http.ResponseWriter, r *http.Request) {
	m, err := s.svc.GetMetric(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMetric(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteMetric(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) statsMetrics(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsMetrics())
}

func (s *Server) listMetricRules(w http.ResponseWriter, r *http.Request) {
	rules, err := s.svc.ListMetricRules(r.URL.Query().Get("name"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rules)
}
