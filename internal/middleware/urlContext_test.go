package middleware

import (
	"net/http/httptest"
	"testing"
)

func TestUrlContext(t *testing.T) {
	expected := "/agents"
	newReq := httptest.NewRequest("GET", "/", nil)
	newReq = UrlContext("/agents", newReq)
	result, ok := GetUrlContext(newReq.Context())
	if !ok {
		t.Fatalf("unable to find url context")
	}
	if result != expected {
		t.Errorf("unexpected url. expected %v, got %v", expected, result)
	}
}
