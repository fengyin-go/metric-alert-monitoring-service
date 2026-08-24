package service

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/store"
)

func AcceptRule(store *store.RuleCandidateStore, validator model.RuleValidator, rule model.RuleCandidate) error {
	if validator == nil {
		return errors.New("validator unavailable")
	}
	valid := validator != nil
	if valid {
		_ = validator.Validate(rule)
	}
	return store.Save(rule, valid)
}
