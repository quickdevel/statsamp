package response

import (
  "net/http"
  "encoding/json"
)

type StatusResponse struct {
  Status  string  `json:"status"`
  Desc    string  `json:"desc"`
}

func JSONResponse(w *http.ResponseWriter, code int, data interface{}) {
  encodedBytes, _ := json.Marshal(data)
  (*w).Header().Set("Content-Type", "application/json")
  (*w).WriteHeader(code)
  (*w).Write(encodedBytes)
}

func JSONResponseError(w *http.ResponseWriter, code int, desc string) {
  JSONResponse(w, code, StatusResponse { Status: "error", Desc: desc })
}

func JSONResponseSuccess(w *http.ResponseWriter, code int, desc string) {
  JSONResponse(w, code, StatusResponse { Status: "success", Desc: desc })
}
