package handler

import (
	"net/http"
	"time"

	"monitoring/internal/model"
	"monitoring/pkg/httpx"
)

func (s *Server) registerSilenceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/silences", s.createSilence)
	mux.HandleFunc("GET /api/silences", s.listSilences)
	mux.HandleFunc("GET /api/silences/active", s.activeSilences)
	mux.HandleFunc("GET /api/silences/stats", s.statsSilences)
	mux.HandleFunc("DELETE /api/silences/metric", s.deleteSilencesByMetric)
	mux.HandleFunc("DELETE /api/silences/{id}", s.deleteSilence)
}

type createSilenceRequest struct {
	MetricName string `json:"metric_name"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	Reason     string `json:"reason"`
	Enabled    *bool  `json:"enabled"`
}

func (s *Server) createSilence(w http.ResponseWriter, r *http.Request) {
	var req createSilenceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		httpx.BadRequest(w, "start_at 时间格式应为 RFC3339")
		return
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		httpx.BadRequest(w, "end_at 时间格式应为 RFC3339")
		return
	}
	si, err := s.svc.CreateSilence(req.MetricName, startAt, endAt, req.Reason, boolValue(req.Enabled))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, si)
}

func (s *Server) listSilences(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SilenceFilter{MetricName: r.URL.Query().Get("metric_name")}
	items, total, err := s.svc.ListSilences(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) deleteSilence(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSilence(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) activeSilences(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ActiveSilences()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) statsSilences(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsSilences())
}

func (s *Server) deleteSilencesByMetric(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.DeleteSilencesByMetric(r.URL.Query().Get("metric_name"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": count})
}
