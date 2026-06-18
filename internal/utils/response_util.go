package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(statusCode int, resWriter http.ResponseWriter, data map[string]any) {
	resWriter.Header().Set("Content-Type", "application/json")
	resWriter.WriteHeader(statusCode)
	json.NewEncoder(resWriter).Encode(data)
}
