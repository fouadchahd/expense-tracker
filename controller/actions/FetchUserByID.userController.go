package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func FetchUserByID(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idString := vars["id"]
	//Fetch User
	userRepo := repository.NewUserRepository(db)
	id, err := strconv.Atoi(idString)
	if err != nil {
		json.NewEncoder(res).Encode(NewApiResponseError(err.Error()))
		return
	}
	userFound, queryErr := userRepo.GetUserByID(uint(id))
	if queryErr != nil {
		json.NewEncoder(res).Encode(NewApiResponseError(queryErr.Error()))
		return
	}
	//Convert Response
	userMap := userRepo.ConvertUserModelToMapResponse(userFound)

	json.NewEncoder(res).Encode(NewApiResponseSuccess(userMap))
}
