package actions

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Response Constructor
const (
	Error   = "error"
	Success = "success"
)

type ApiStatus string

type ApiResponse struct {
	Data   any       `json:"data"`
	Status ApiStatus `json:"status"`
}

func NewApiResponseError(msg string) *ApiResponse {
	return &ApiResponse{
		Data:   msg,
		Status: Error,
	}
}

func NewApiResponseSuccess(data any) *ApiResponse {
	return &ApiResponse{
		Data:   data,
		Status: Success,
	}
}

func CheckAuth(key string, token string, req *http.Request) (*ApiResponse, error) {
	//Start checking auth
	if len(key) == 0 || len(token) == 0 {
		return nil, errors.New("missing headers")
	}
	subRoute := mux.CurrentRoute(req).Subrouter()
	url, err := subRoute.Get("check_authorization").URL("id", key, "password", token)
	subRoute.Methods(http.MethodGet)
	if err != nil {
		return nil, errors.New(err.Error())

	}
	protocol := "http"
	if req.Proto == "HTTP/2" {
		protocol = "https"
	}
	authRes, authErr := http.Get(protocol + "://" + req.Host + url.Path)
	if authErr != nil {
		log.Fatal(authErr)
	}
	apiResponse := new(ApiResponse)
	json.NewDecoder(authRes.Body).Decode(&apiResponse)
	if apiResponse.Status == Error {
		return nil, errors.New(apiResponse.Data.(string))
	}
	return apiResponse, nil
}
