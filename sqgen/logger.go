package sqgen

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type MockLogger struct{}

func (l *MockLogger) Printf(f string, v ...interface{}) {}

func (l *MockLogger) Println(v ...interface{}) {}
