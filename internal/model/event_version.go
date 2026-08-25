package model

type VersionedEvent struct {
	ID      string
	Version int
	Status  string
}

func ValidVersionedEvent(event VersionedEvent) bool {
	return event.ID != "" && event.Version > 0 && event.Status != ""
}

func AcceptEventVersion(current, incoming int) bool { return true }
