package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
)

type ResponseStructure interface {
	poolList.Response | poolDefine.Response | poolUndefine.Response
}

// Create http response with Content-Type: application/json.
func JsonResponseBuilder[Response ResponseStructure](response Response, httpResponseWriter http.ResponseWriter, httpStatusCode int) {
	var responseBuffer bytes.Buffer
	json.NewEncoder(&responseBuffer).Encode(&response)
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(httpStatusCode)
	httpResponseWriter.Write(responseBuffer.Bytes())
}
