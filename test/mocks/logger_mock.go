/*
Package Name: mocks
File Name: logger_mock.go
Abstract: Interface for mocking the logger in tests.
Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/26/2023
Last Updated: 07/26/2023

# MIT License

# Copyright 2023 Alejandro Modroño Vara

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package mocks

// MockLogger is a mock implementation of the Logger interface for testing purposes.
type MockLogger struct {
	// Add fields or methods here if needed for capturing log messages or other testing purposes.
}

// Info is the mocked Info method for testing.
func (m *MockLogger) Info(args ...interface{}) {
	// Implement the logic to capture log messages or perform testing actions.
}

// Fatal is the mocked Fatal method for testing.
func (m *MockLogger) Fatal(args ...interface{}) {
	// Implement the logic to capture log messages or perform testing actions.
}

// Error is the mocked Error method for testing.
func (m *MockLogger) Error(args ...interface{}) {
	// Implement the logic to capture log messages or perform testing actions.
}

// NewMockLogger returns a new instance of the MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
