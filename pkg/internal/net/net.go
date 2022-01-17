package net

import (
	"encoding/json"
	"net/http"
)

// Bind :
func Bind(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// WriteJSON :
func WriteJSON(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(response)
}
