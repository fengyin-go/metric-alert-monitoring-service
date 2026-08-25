package service

import "fmt"

func EventDeliveryKey(eventID string, attempt int) string {
	return fmt.Sprintf("event:%s:attempt:%d", eventID, attempt)
}
