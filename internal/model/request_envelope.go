package model

type RequestEnvelope struct {
	Tenant string
	Labels []string
}

func (e *RequestEnvelope) Reset()                    { e.Labels = e.Labels[:0] }
func (e *RequestEnvelope) Snapshot() RequestEnvelope { return *e }
