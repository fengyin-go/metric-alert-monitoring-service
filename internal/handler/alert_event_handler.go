package handler

import (
	"net/http"

	"monitoring/internal/model"
	"monitoring/pkg/httpx"
)

func (s *Server) registerAlertEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/events/evaluate", s.evaluateRules)
	mux.HandleFunc("GET /api/events", s.listAlertEvents)
	mux.HandleFunc("GET /api/events/stats", s.statsEvents)
	mux.HandleFunc("GET /api/events/active", s.activeEvents)
	mux.HandleFunc("GET /api/events/count", s.countEventsByMetric)
	mux.HandleFunc("GET /api/events/rule", s.eventsForRule)
	mux.HandleFunc("GET /api/events/{id}", s.getAlertEvent)
	mux.HandleFunc("POST /api/events/{id}/resolve", s.resolveEvent)
}

func (s *Server) evaluateRules(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.EvaluateRules()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) listAlertEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertEventFilter{
		MetricName: r.URL.Query().Get("metric_name"),
		Severity:   r.URL.Query().Get("severity"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListAlertEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlertEvent(w http.ResponseWriter, r *http.Request) {
	e, err := s.svc.GetAlertEvent(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) resolveEvent(w http.ResponseWriter, r *http.Request) {
	e, err := s.svc.ResolveEvent(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) statsEvents(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.StatsEvents())
}

func (s *Server) activeEvents(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ActiveEvents()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) countEventsByMetric(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.CountEventsByMetric(r.URL.Query().Get("metric_name"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"count": count})
}

func (s *Server) eventsForRule(w http.ResponseWriter, r *http.Request) {
	events, err := s.svc.EventsForRule(r.URL.Query().Get("rule_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, events)
}
