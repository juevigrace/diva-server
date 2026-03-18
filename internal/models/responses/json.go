package responses

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, res *APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
