package model

import "errors"

type RuleCandidate struct {
	Name      string
	Threshold *float64
}
type RuleValidator interface{ Validate(RuleCandidate) error }
type ThresholdValidator struct{}

func (v *ThresholdValidator) Validate(rule RuleCandidate) error {
	if v == nil || rule.Threshold == nil {
		return errors.New("threshold is required")
	}
	return nil
}
func NewRuleValidator(enabled bool) RuleValidator {
	if !enabled {
		var disabled *ThresholdValidator
		return disabled
	}
	return &ThresholdValidator{}
}
