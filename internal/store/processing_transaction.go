package store

type ProcessingTransaction struct{ committed, rolledBack bool }

// Finish 根据处理结果终结事务：成功则提交，失败则回滚。
func (t *ProcessingTransaction) Finish(processErr error) error {
	if processErr != nil {
		t.rolledBack = true
		return nil
	}
	t.committed = true
	return nil
}
func (t *ProcessingTransaction) Committed() bool  { return t.committed }
func (t *ProcessingTransaction) RolledBack() bool { return t.rolledBack }
