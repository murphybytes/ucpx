package common

import "testing"

type mockLogger struct {
	msg string
}

func (m *mockLogger) print(vals ...interface{}) {
	m.msg = ""
	for _, val := range vals {
		m.msg += val.(string)
	}
}

func TestErrorLogger(t *testing.T) {
	m := &mockLogger{}
	l, _ := newLogger(&Flags{LogLevel: logError}, m.print)
	l.LogInfo("info")
	if m.msg != "" {
		t.Error("Info message shouldn't be logged")
	}
	l.LogWarn("warn")
	if m.msg != "" {
		t.Error("Warn message shouldn't be logged")
	}

	l.LogError("error")

	if m.msg != "errorss" {
		t.Error("Message should be logged")
	}

}

func TestWarnLogger(t *testing.T) {
	m := &mockLogger{}
	l, _ := newLogger(&Flags{LogLevel: logWarn}, m.print)
	l.LogInfo("info")
	if m.msg != "" {
		t.Error("Info message shouldn't be logged")
	}
	l.LogWarn("warn")
	if m.msg != "warn" {
		t.Error("Warn message should be logged")
	}

	l.LogError("error")

	if m.msg != "error" {
		t.Error("Message should be logged")
	}

}

func TestInfoLogger(t *testing.T) {
	m := &mockLogger{}
	l, _ := newLogger(&Flags{LogLevel: logInfo}, m.print)
	l.LogInfo("info")
	if m.msg != "info" {
		t.Error("Info message should be logged")
	}
	l.LogWarn("warn")
	if m.msg != "warn" {
		t.Error("Warn message should be logged")
	}

	l.LogError("error")

	if m.msg != "error" {
		t.Error("Message should be logged")
	}

}

func TestLoggerFailed(t *testing.T) {
	m := &mockLogger{}
	_, e := newLogger(&Flags{LogLevel: "foo"}, m.print)
	if e == nil {
		t.Error("There is not an error foo")
	}

}
