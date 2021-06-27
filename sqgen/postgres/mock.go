package postgres

var _ Logger = (*mockLogger)(nil)

type mockLogger struct{}

func (l *mockLogger) Printf(f string, v ...interface{}) {}

func (l *mockLogger) Println(v ...interface{}) {}
