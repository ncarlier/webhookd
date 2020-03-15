package api

import (
	"encoding/json"
	"net/http"

	"github.com/ncarlier/webhookd/pkg/version"
)

// Info API informations model structure.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	info := Info{
		Name:    "webhookd",
		Version: version.Version,
	}
	data, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
