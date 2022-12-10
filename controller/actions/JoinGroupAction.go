package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"github.com/fouadelhamri/expense-tracker/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func JoinGroupAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	groupDestination := params["id"]
	//Check Group if found
	groupRepo := repository.NewGroupRepository(db)
	group, fetchErr := groupRepo.GetGroupByCode(groupDestination)
	if fetchErr != nil {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(NewApiResponseError(fetchErr.Error()))
		return
	}

	//Check Authorization from headers
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

	//check if part of this group
	currentUserID, _ := strconv.Atoi(key)

	if len(group.Users) > 0 {
		for _, user := range group.Users {
			if user.ID == uint(currentUserID) {
				res.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(res).Encode(NewApiResponseError("already part of this group"))
				return
			}
		}
	}
	//check if user already part of a group
	groupID := apiResponse.Data.(map[string]any)["group_id"]
	if groupID != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("you can only be part of one group at time"))
		return
	}
	// not part of any group
	db.Model(&model.User{
		ID: uint(currentUserID),
	}).Update("group_id", groupDestination)
}
