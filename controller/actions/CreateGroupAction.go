package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func CreateGroupAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	//Check Authorization from header
	key := req.Header.Get("key")
	token := req.Header.Get("token")
	if len(key) == 0 || len(token) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("missing headers"))
		return
	}
	subRoute := mux.CurrentRoute(req).Subrouter()
	url, err := subRoute.Get("check_authorization").URL("id", key, "password", token)
	subRoute.Methods(http.MethodPost)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(err.Error()))
		return
	}
	protocol := "http"
	if req.Proto == "HTTP/2" {
		protocol = "https"
	}
	authRes, authErr := http.Get(protocol + "://" + req.Host + url.Path)
	if authErr != nil {
		log.Fatal(authErr)
	}

	var apiResponse ApiResponse
	json.NewDecoder(authRes.Body).Decode(&apiResponse)
	if apiResponse.Status == Error {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(NewApiResponseError("unauthorized user"))
		return
	}
	//Authorized User FOUND
	//check if user already part of a group
	groupID := apiResponse.Data.(map[string]any)["group_id"]
	if groupID != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("already part of another group"))
		return
	}
	// not part of any group
	// create Group
	autoGenCode := getAuthGenCode(12)
	currentUserID, _ := strconv.Atoi(key)
	newGroup := model.Group{
		Code:      autoGenCode,
		CreatedAt: time.Now().UTC(),
		Users: []model.User{{
			ID: uint(currentUserID),
		}},
	}
	createRow := db.Create(&newGroup)
	if createRow.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(createRow.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(newGroup))
}

func getAuthGenCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
