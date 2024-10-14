package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Create http response with Content-Type: application/json.
func JsonResponseBuilder(response interface{}, httpResponseWriter http.ResponseWriter, httpStatusCode int) {
	var responseBuffer bytes.Buffer
	json.NewEncoder(&responseBuffer).Encode(&response)
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(httpStatusCode)
	httpResponseWriter.Write(responseBuffer.Bytes())
}
