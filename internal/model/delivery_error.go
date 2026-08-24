package model

import "fmt"

type RejectError struct{ Reason string }

func (e *RejectError) Error() string { return "delivery rejected: " + e.Reason }
func PreserveDeliveryError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("delivery failed: %v", err)
}
