package logger

import (
	"context"
	"fmt"
	"sync"
)

type MockILogger struct {
	mu          sync.Mutex
	logs        map[string][]string // Armazena logs por contexto
	context     string              // Contexto atual
	DebugCalled bool
	ErrorCalled bool
}

func NewMockILogger() *MockILogger {
	return &MockILogger{
		logs: make(map[string][]string),
	}
}

// SetContext define o contexto atual
func (m *MockILogger) SetContext(context string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context = context
}

func (m *MockILogger) Debug(msg string, args ...any) {
	m.DebugCalled = true
	m.log("DEBUG", msg, args...)
}

func (m *MockILogger) Info(msg string, args ...any) {
	m.log("INFO", msg, args...)
}

func (m *MockILogger) Warn(msg string, args ...any) {
	m.log("WARN", msg, args...)
}

func (m *MockILogger) Error(msg string, args ...any) {
	m.ErrorCalled = true
	m.log("ERROR", msg, args...)
}

// With mocks the With method, returning the same logger for simplicity
func (m *MockILogger) With(args ...any) ILogger {
	m.log("WITH", "Adding context %v", args)
	return m
}

// WithContext mocks the WithContext method, returning the same logger for simplicity
func (m *MockILogger) WithContext(ctx context.Context) ILogger {
	m.log("WITH_CONTEXT", "Adding context with context.Context %v", ctx)
	return m
}

// DebugContext mocks the DebugContext logging method
func (m *MockILogger) DebugContext(ctx context.Context, msg string, args ...any) {
	m.log("DEBUG_CONTEXT", msg, args...)
}

// InfoContext mocks the InfoContext logging method
func (m *MockILogger) InfoContext(ctx context.Context, msg string, args ...any) {
	m.log("INFO_CONTEXT", msg, args...)
}

// WarnContext mocks the WarnContext logging method
func (m *MockILogger) WarnContext(ctx context.Context, msg string, args ...any) {
	m.log("WARN_CONTEXT", msg, args...)
}

// ErrorContext mocks the ErrorContext logging method
func (m *MockILogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	m.log("ERROR_CONTEXT", msg, args...)
}

func (m *MockILogger) log(level, msg string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.context == "" {
		m.context = "default"
	}
	formatted := fmt.Sprintf("[%s] %s", level, fmt.Sprintf(msg, args...))
	m.logs[m.context] = append(m.logs[m.context], formatted)
}

func (m *MockILogger) GetLogs(context string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.logs[context]
}
