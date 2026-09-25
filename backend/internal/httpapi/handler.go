package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"

	"github.com/takehome/calculator/backend/internal/calculator"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorEnvelope{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func parseRequest(r *http.Request, requireB bool) (*RequestDTO, int, string, string) {
	if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type") != "application/json" {
		return nil, http.StatusBadRequest, "INVALID_BODY", "Content-Type must be application/json"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		return nil, http.StatusBadRequest, "INVALID_BODY", "request body cannot be empty"
	}

	var req RequestDTO
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, http.StatusBadRequest, "INVALID_BODY", "malformed JSON request body"
	}

	if req.A == nil {
		return nil, http.StatusBadRequest, "MISSING_FIELD", "field 'a' is required"
	}

	if math.IsNaN(*req.A) || math.IsInf(*req.A, 0) {
		return nil, http.StatusBadRequest, "NOT_A_NUMBER", "field 'a' must be a valid finite number"
	}

	if requireB {
		if req.B == nil {
			return nil, http.StatusBadRequest, "MISSING_FIELD", "field 'b' is required"
		}
		if math.IsNaN(*req.B) || math.IsInf(*req.B, 0) {
			return nil, http.StatusBadRequest, "NOT_A_NUMBER", "field 'b' must be a valid finite number"
		}
	}

	return &req, 0, "", ""
}

func mapCalculatorError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusBadRequest, "DIVISION_BY_ZERO", err.Error())
	case errors.Is(err, calculator.ErrNegativeSqrtInput):
		writeError(w, http.StatusBadRequest, "NEGATIVE_SQRT_INPUT", err.Error())
	case errors.Is(err, calculator.ErrOverflow):
		writeError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "an unexpected calculation error occurred")
	}
}

func handleBinaryOp(w http.ResponseWriter, r *http.Request, op func(a, b float64) (float64, error)) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method not allowed")
		return
	}

	req, errCode, code, msg := parseRequest(r, true)
	if errCode != 0 {
		writeError(w, errCode, code, msg)
		return
	}

	res, err := op(*req.A, *req.B)
	if err != nil {
		mapCalculatorError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, SuccessDTO{Result: res})
}

func handleUnaryOp(w http.ResponseWriter, r *http.Request, op func(a float64) (float64, error)) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method not allowed")
		return
	}

	req, errCode, code, msg := parseRequest(r, false)
	if errCode != 0 {
		writeError(w, errCode, code, msg)
		return
	}

	res, err := op(*req.A)
	if err != nil {
		mapCalculatorError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, SuccessDTO{Result: res})
}

// HandleHealth returns application liveness status.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "HTTP method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, HealthDTO{Status: "ok"})
}

// HandleAdd handles addition POST requests.
func HandleAdd(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Add)
}

// HandleSubtract handles subtraction POST requests.
func HandleSubtract(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Subtract)
}

// HandleMultiply handles multiplication POST requests.
func HandleMultiply(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Multiply)
}

// HandleDivide handles division POST requests.
func HandleDivide(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Divide)
}

// HandlePower handles power POST requests.
func HandlePower(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Power)
}

// HandleSqrt handles square root POST requests.
func HandleSqrt(w http.ResponseWriter, r *http.Request) {
	handleUnaryOp(w, r, calculator.Sqrt)
}

// HandlePercentage handles percentage POST requests.
func HandlePercentage(w http.ResponseWriter, r *http.Request) {
	handleBinaryOp(w, r, calculator.Percentage)
}

// HandleNotFound handles 404 routes.
func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "NOT_FOUND", "route not found")
}
