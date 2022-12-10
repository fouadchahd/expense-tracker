package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func LoginAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	body := make(map[string]string, 0)
	decoderErr := json.NewDecoder(req.Body).Decode(&body)
	if decoderErr != nil {
		if decoderErr == io.EOF {
			json.NewEncoder(res).Encode(NewApiResponseError("no request body provided"))
			return
		}
		json.NewEncoder(res).Encode(NewApiResponseError(decoderErr.Error()))
		return
	}
	//extract Request
	username, ok1 := body["username"]
	password, ok2 := body["password"]
	if !(ok1 && ok2) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("bad request"))
		return
	}
	userRepo := repository.NewUserRepository(db)
	userFound, fetchErr := userRepo.GetUserByUsername(username)
	if fetchErr != nil {
		json.NewEncoder(res).Encode(NewApiResponseError("user no found"))
		return
	}

	matchedPassword := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(password))
	if matchedPassword != nil {
		json.NewEncoder(res).Encode(NewApiResponseError("incorrect password"))
		return
	}
	userResponse := userRepo.ConvertUserModelToMapResponse(userFound)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(userResponse))
	return
}
