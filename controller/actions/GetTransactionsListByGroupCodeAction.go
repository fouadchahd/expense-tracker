package actions

import (
	"encoding/json"
	"fmt"
	"github.com/fouadelhamri/expense-tracker/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func GetTransactionsListByGroupCodeAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	//look for group
	groupRepo := repository.NewGroupRepository(db)
	group, fetchErr := groupRepo.GetGroupTransactionsByCode(params["id"])
	if fetchErr != nil {
		json.NewEncoder(res).Encode(NewApiResponseError(fetchErr.Error()))
		return
	}
	//check user auth
	key := req.Header.Get("key")
	token := req.Header.Get("token")
	apiResponse, checkErr := CheckAuth(key, token, req)
	if checkErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(checkErr.Error()))
		return
	}
	userMap := apiResponse.Data.(map[string]interface{})
	fmt.Println("User Found : ", userMap)
	if (userMap["group_id"] == nil) || float64(group.ID) != userMap["group_id"].(float64) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("you are not part of this group"))
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(group.Transactions))
}
