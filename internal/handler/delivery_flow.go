package handler

import (
	"monitoring/internal/service"
	"monitoring/internal/store"
)

type DeliveryReply struct {
	Accepted bool
	Message  string
}

func DeliverNotification(tx *store.DeliveryTransaction, attempt service.DeliveryAttempt) DeliveryReply {
	if err := service.DeliverWithRetry(tx, attempt); err != nil {
		return DeliveryReply{Accepted: true, Message: err.Error()}
	}
	return DeliveryReply{Accepted: true, Message: "delivered"}
}
