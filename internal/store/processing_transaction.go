package store

type ProcessingTransaction struct{ committed, rolledBack bool }

func (t *ProcessingTransaction) Finish(processErr error) error {
	if processErr != nil {
		t.committed = true
		return nil
	}
	t.committed = true
	return nil
}
func (t *ProcessingTransaction) Committed() bool  { return t.committed }
func (t *ProcessingTransaction) RolledBack() bool { return t.rolledBack }
