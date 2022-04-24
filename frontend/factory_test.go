package frontend

import "testing"

func TestFrontendFactory(t *testing.T) {
	fe, err := NewFrontEnd("")
	if err == nil {
		t.Error("should have returned an error for empty FrontEnd")
	}
	if fe != nil {
		t.Error("should have")
	}
}
