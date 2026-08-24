package service

import (
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/store"
)

type DeliveryAttempt func() error

func DeliverWithRetry(tx *store.DeliveryTransaction, attempt DeliveryAttempt) error {
	for i := 0; i < 2; i++ {
		tx.Stage()
		err := model.PreserveDeliveryError(attempt())
		tx.Finish(err)
		if err == nil {
			return nil
		}
		var rejected *model.RejectError
		if errors.As(err, &rejected) {
			continue
		}
	}
	return errors.New("delivery retries exhausted")
}
