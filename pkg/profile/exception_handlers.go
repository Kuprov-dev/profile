package profile

import (
	"encoding/json"
	"net/http"
)

func makeInternalServerErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusInternalServerError)
	body, _ := json.Marshal(ResponseBody{Status: 500, Message: "Internal error."})
	(*w).Write(body)
}

func makeBadGatewayErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadGateway)
	body, _ := json.Marshal(ResponseBody{Status: 502, Message: "Bad gateway."})
	(*w).Write(body)
}

func makeServiceUnavailableErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusServiceUnavailable)
	body, _ := json.Marshal(ResponseBody{Status: 503, Message: "External service is unavailable."})
	(*w).Write(body)
}

func makeBadRequestErrorResponse(w *http.ResponseWriter, errMsg string) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadRequest)
	body, _ := json.Marshal(ResponseBody{Status: 400, Message: errMsg})
	(*w).Write(body)
}

func makeUnathorisedErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(ResponseBody{Status: 401, Message: "Unathorised error."})
	(*w).Write(body)
}

func makeForbiddenErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(ResponseBody{Status: 403, Message: "Forbidden error."})
	(*w).Write(body)
}

func makeNotFoundErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(ResponseBody{Status: 404, Message: "Not found error."})
	(*w).Write(body)
}
