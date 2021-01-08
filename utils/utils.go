package utils

import (
	"encoding/json"
	"net/http"
	"github.com/mailru/easyjson"
)

func WriteJson(w http.ResponseWriter, status int, data easyjson.Marshaler) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	easyjson.MarshalToHTTPResponseWriter(data, w)
}

func WriteNotEasyJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}