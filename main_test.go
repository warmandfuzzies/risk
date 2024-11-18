// endpoints_test.go
package risk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRiskHandler(t *testing.T) {

	posturl := "http://localhost:8080"

	body := []byte(`{
		"state": "open",
		"title": "I am a title",
		"description": "I am a description"
	}`)

	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newRiskHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var risk Risk

	derr := json.NewDecoder(rr.Body).Decode(&risk)
	if derr != nil {
		t.Errorf("Unable to decode risk")
	}

	if risk.State != "open" {
		t.Errorf("handler returned unexpected state: got %v want %v",
			risk.State, "open")
	}

}
