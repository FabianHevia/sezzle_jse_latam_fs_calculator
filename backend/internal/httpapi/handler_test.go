package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takehome/calculator/backend/internal/httpapi"
)

func executeRequest(router http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestHealthEndpoint(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")
	rr := executeRequest(router, "GET", "/api/health", "")

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp httpapi.HealthDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestAddEndpoint(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")

	t.Run("successful addition", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/add", `{"a": 10, "b": 5}`)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		var resp httpapi.SuccessDTO
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Result != 15 {
			t.Errorf("expected result 15, got %v", resp.Result)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/add", `{"a": 10}`)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}

		var errResp httpapi.ErrorEnvelope
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error.Code != "MISSING_FIELD" {
			t.Errorf("expected code MISSING_FIELD, got %s", errResp.Error.Code)
		}
	})

	t.Run("invalid json body", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/add", `{invalid json}`)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}

		var errResp httpapi.ErrorEnvelope
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error.Code != "INVALID_BODY" {
			t.Errorf("expected code INVALID_BODY, got %s", errResp.Error.Code)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		rr := executeRequest(router, "GET", "/api/add", "")
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected status 405, got %d", rr.Code)
		}

		var errResp httpapi.ErrorEnvelope
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error.Code != "METHOD_NOT_ALLOWED" {
			t.Errorf("expected code METHOD_NOT_ALLOWED, got %s", errResp.Error.Code)
		}
	})
}

func TestDivideEndpoint(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")

	t.Run("successful division", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/divide", `{"a": 20, "b": 4}`)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		var resp httpapi.SuccessDTO
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Result != 5 {
			t.Errorf("expected result 5, got %v", resp.Result)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/divide", `{"a": 10, "b": 0}`)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}

		var errResp httpapi.ErrorEnvelope
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error.Code != "DIVISION_BY_ZERO" {
			t.Errorf("expected code DIVISION_BY_ZERO, got %s", errResp.Error.Code)
		}
	})
}

func TestSqrtEndpoint(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")

	t.Run("successful sqrt", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/sqrt", `{"a": 64}`)
		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		var resp httpapi.SuccessDTO
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Result != 8 {
			t.Errorf("expected result 8, got %v", resp.Result)
		}
	})

	t.Run("negative sqrt input", func(t *testing.T) {
		rr := executeRequest(router, "POST", "/api/sqrt", `{"a": -16}`)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}

		var errResp httpapi.ErrorEnvelope
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Error.Code != "NEGATIVE_SQRT_INPUT" {
			t.Errorf("expected code NEGATIVE_SQRT_INPUT, got %s", errResp.Error.Code)
		}
	})
}

func TestNotFound(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")
	rr := executeRequest(router, "GET", "/api/nonexistent", "")

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}

	var errResp httpapi.ErrorEnvelope
	_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
	if errResp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected code NOT_FOUND, got %s", errResp.Error.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	router := httpapi.NewRouter("http://localhost:5173")
	req := httptest.NewRequest("OPTIONS", "/api/add", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for OPTIONS preflight, got %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("unexpected CORS origin header: %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}
}
