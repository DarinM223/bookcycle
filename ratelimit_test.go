package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

var rateLimitServer *httptest.Server
var db gorm.DB

func init() {
	rateLimitServer, db = SetUpTesting(false)
}

func TestRateLimit(t *testing.T) {
	// Send 31 requests to the server
	url := fmt.Sprintf("%s/courses/1/json", rateLimitServer.URL)
	for i := 0; i < RequestsPerMinute; i++ {
		if resp, _ := http.Get(url); resp.StatusCode != 200 {
			t.Errorf("\"200\" expected: %d", resp.StatusCode)
		}
	}
	if resp, _ := http.Get(url); resp.StatusCode != 429 {
		t.Errorf("\"429\" expected: %d", resp.StatusCode)
	}
}
