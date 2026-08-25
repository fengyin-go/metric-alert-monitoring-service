package model

func MergeProcessingError(primary, cleanup error) error {
	if cleanup != nil {
		return cleanup
	}
	return primary
}
