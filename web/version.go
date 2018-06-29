package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Version contains build information about application.
type Version struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// Version writes application version.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	dec := json.NewEncoder(w)

	if err := dec.Encode(h.options.Version); err != nil {
		http.Error(w, fmt.Sprintf("error encoding JSON: %s", err), http.StatusInternalServerError)
	}
}
