package utils

import "net/http"

// Create http response with no content. Use case is for "PUT", "PATCH" and "DELETE" method.
func NoContentResponseBuilder(httpResponseWriter http.ResponseWriter) {
	httpResponseWriter.WriteHeader(http.StatusNoContent)
	httpResponseWriter.Write(nil)
}
