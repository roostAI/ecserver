package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BasicResponse struct {
	Status  int
	Message string
}

func handleInvalidMethod(w http.ResponseWriter, r *http.Request) {
	resp := &BasicResponse{
		Status:  http.StatusMethodNotAllowed,
		Message: "No such endpoint",
	}
	writeBasicResponse(w, resp)
}

func writeBasicResponse(w http.ResponseWriter, resp *BasicResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func TestHandleInvalidMethod_58fd8542d1(t *testing.T) {
	// Test case 1: Check if the status code returned is 405 (Method Not Allowed)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleInvalidMethod)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	// Test case 2: Check if the message returned is "No such endpoint"
	expected := `{"Status":405,"Message":"No such endpoint"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
