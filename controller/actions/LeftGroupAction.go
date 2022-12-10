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

func LeftGroupAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
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
				if len(group.Users) == 1 {
					//Remove the group
					deleteRow := db.Delete(&group)
					if deleteRow.Error != nil {
						res.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(res).Encode(NewApiResponseError(deleteRow.Error.Error()))
						return
					}
				}
				updateRow := db.Model(&model.User{}).Where("id = ?", currentUserID).Update("group_id", nil)
				if updateRow.Error != nil {
					res.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(res).Encode(NewApiResponseError(updateRow.Error.Error()))
					return
				}
				res.WriteHeader(http.StatusOK)
				json.NewEncoder(res).Encode(NewApiResponseSuccess("You left the group"))
				return
			}
		}
	}
	res.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(res).Encode(NewApiResponseError("You are not part of this group"))
	return
}
