package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func AuthorizationAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idString := vars["id"]
	hashedPasswordString := vars["password"]
	//Fetch User
	userRepo := repository.NewUserRepository(db)
	id, err := strconv.Atoi(idString)
	if err != nil || len(hashedPasswordString) == 0 {
		json.NewEncoder(res).Encode(NewApiResponseError("invalid credentials"))
		return
	}
	userFound, queryErr := userRepo.GetUserByCredentials(uint(id), hashedPasswordString)
	if queryErr != nil {
		json.NewEncoder(res).Encode(NewApiResponseError(queryErr.Error()))
		return
	}

	json.NewEncoder(res).Encode(NewApiResponseSuccess(userFound))
}
