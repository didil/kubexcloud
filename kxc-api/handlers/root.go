package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/services"
)

type Root struct {
	ProjectSvc services.ProjectSvc
	AppSvc     services.AppSvc
}

// HandleError handles errors
func (root *Root) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	JSONError(w, err.Error(), http.StatusBadRequest)
}

// json helpers

// JSONErr err
type JSONErr struct {
	Err string `json:"err"`
}

// JSONError renders json with error
func JSONError(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	writeJSON(w, &JSONErr{Err: errStr})
}

// JSONOk renders json with 200 ok
func JSONOk(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, v)
}

// writeJSON to response body
func writeJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	writeBytes(w, b)
}

// writeBytes to response body
func writeBytes(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		http.Error(w, fmt.Sprintf("write error: %v", err), http.StatusInternalServerError)
		return
	}
}

// readJSON from request body
func readJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("invalid JSON input")
	}

	return nil
}

// CtxKey context key
type CtxKey string
