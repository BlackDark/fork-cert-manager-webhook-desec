package desec

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	oldBaseURL := baseURL
	baseURL = srv.URL
	t.Cleanup(func() {
		srv.Close()
		baseURL = oldBaseURL
	})
}

func TestAddRecord_returnsAdded_whenNewRecord(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(RRSets{})
		case "PUT":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	api := &API{Token: "test-token"}
	_, added, err := api.AddRecord(context.Background(), "_acme-challenge", "example.com", "TXT", "somekey", 3600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !added {
		t.Error("expected added=true for new record")
	}
}

func TestAddRecord_returnsNotAdded_whenRecordAlreadyExists(t *testing.T) {
	existing := RRSets{{SubName: "_acme-challenge", Type: "TXT", Records: []string{`"somekey"`}, TTL: 3600}}
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(existing)
		}
	})

	api := &API{Token: "test-token"}
	_, added, err := api.AddRecord(context.Background(), "_acme-challenge", "example.com", "TXT", "somekey", 3600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if added {
		t.Error("expected added=false when record already exists")
	}
}

func TestDeleteRecord_returnsDeleted_whenRecordExists(t *testing.T) {
	existing := RRSets{{SubName: "_acme-challenge", Type: "TXT", Records: []string{`"somekey"`}, TTL: 3600}}
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(existing)
		case "PUT":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	api := &API{Token: "test-token"}
	_, deleted, err := api.DeleteRecord(context.Background(), "_acme-challenge", "example.com", "TXT", "somekey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected deleted=true when record was found and removed")
	}
}

func TestDeleteRecord_returnsNotDeleted_whenRRSetNotFound(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(RRSets{})
		}
	})

	api := &API{Token: "test-token"}
	_, deleted, err := api.DeleteRecord(context.Background(), "_acme-challenge", "example.com", "TXT", "somekey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted {
		t.Error("expected deleted=false when RRSet not found")
	}
}

func TestDeleteRecord_returnsNotDeleted_whenRecordNotInRRSet(t *testing.T) {
	existing := RRSets{{SubName: "_acme-challenge", Type: "TXT", Records: []string{`"otherkey"`}, TTL: 3600}}
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.Contains(r.URL.RawQuery, "subname=_acme-challenge") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(existing)
		}
	})

	api := &API{Token: "test-token"}
	_, deleted, err := api.DeleteRecord(context.Background(), "_acme-challenge", "example.com", "TXT", "somekey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted {
		t.Error("expected deleted=false when record not in RRSet")
	}
}
