package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/repository"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func RegisterAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	var body = make(map[string]any, 0)
	decodeErr := json.NewDecoder(req.Body).Decode(&body)

	if decodeErr != nil {
		if decodeErr == io.EOF {
			json.NewEncoder(res).Encode(NewApiResponseError("no request body provided"))
			return
		}
		json.NewEncoder(res).Encode(NewApiResponseError(decodeErr.Error()))
		return
	}
	//extract Request
	username, ok1 := body["username"]
	password, ok2 := body["password"]
	role, ok3 := body["role"]
	meta, ok4 := body["meta"]
	if !(ok1 && ok2 && ok3 && ok4) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("bad request"))
		return
	}

	//CreateUser
	userRepo := repository.NewUserRepository(db)
	userID, createUserErr := userRepo.CreateUser(username.(string), password.(string), role.(string), meta.(map[string]any))
	//Send Response
	if userID == 0 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(createUserErr.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(userID))
}
