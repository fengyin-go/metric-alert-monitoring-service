package handler

import (
	"monitoring/internal/model"
	"monitoring/internal/service"
	"monitoring/internal/store"
)

type RuleCandidateReply struct {
	Accepted bool
	Message  string
}

func SubmitRuleCandidate(s *store.RuleCandidateStore, validator model.RuleValidator, rule model.RuleCandidate) RuleCandidateReply {
	if err := service.AcceptRule(s, validator, rule); err != nil {
		return RuleCandidateReply{Accepted: true, Message: err.Error()}
	}
	return RuleCandidateReply{Accepted: true, Message: "accepted"}
}
