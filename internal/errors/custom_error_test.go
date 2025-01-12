package errors

import (
    "testing"
)

func TestCustomError( t *testing.T) {
    code := 404
    message := "Resource not found"

    err := New(code, message)

    if err.Code != code {
        t.Errorf("Expected Code %d, received %d", code, err.Code)
    }

    if err.Message != message {
        t.Errorf("Expected message %s, received %s", message, err.Message)
    }

    expectedErrorString := "Code: 404, Message: Resource not found"
    if err.Error() != expectedErrorString {
		t.Errorf("expected Error() %q, got %q", expectedErrorString, err.Error())
	}
}