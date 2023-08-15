package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Candidate struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type BasicResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var candidates = []*Candidate{
	&Candidate{
		Name:     "John Doe",
		ImageUrl: "http://example.com/johndoe.jpg",
	},
}

func addCandidate(w http.ResponseWriter, r *http.Request) {
	newCandidate := &Candidate{}
	err := json.NewDecoder(r.Body).Decode(newCandidate)
	if err != nil {
		resp := &BasicResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		writeBasicResponse(w, resp)
		return
	}

	isCandidatePresent := false
	for i, ca := range candidates {
		if newCandidate.Name == ca.Name {
			isCandidatePresent = true
			candidates[i].Name = newCandidate.Name
			candidates[i].ImageUrl = newCandidate.ImageUrl
		}
	}

	if !isCandidatePresent {
		candidates = append(candidates, newCandidate)
	}

	writeAllCandidatesResponse(w)
}

func writeBasicResponse(w http.ResponseWriter, resp *BasicResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

func writeAllCandidatesResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

func TestAddCandidate_8f9ed01d08(t *testing.T) {
	t.Run("test add new candidate", func(t *testing.T) {
		data := []byte(`{"name":"Jane Doe","imageUrl":"http://example.com/janedoe.jpg"}`)
		req, err := http.NewRequest("POST", "/addCandidate", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(addCandidate)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"name":"Jane Doe","imageUrl":"http://example.com/janedoe.jpg"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("test add existing candidate", func(t *testing.T) {
		data := []byte(`{"name":"John Doe","imageUrl":"http://example.com/johndoe.jpg"}`)
		req, err := http.NewRequest("POST", "/addCandidate", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(addCandidate)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"name":"John Doe","imageUrl":"http://example.com/johndoe.jpg"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})

	t.Run("test add candidate with invalid payload", func(t *testing.T) {
		data := []byte(`{"name":"John Doe","url":"http://example.com/johndoe.jpg"}`)
		req, err := http.NewRequest("POST", "/addCandidate", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(addCandidate)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"status":400,"message":"Invalid request payload"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	})
}
