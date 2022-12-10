package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"gorm.io/gorm"
	"net/http"
)

func SubmitTransactionAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	var transaction *model.Transaction
	marshalErr := json.NewDecoder(req.Body).Decode(&transaction)

	if marshalErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(marshalErr.Error()))
		return
	}

	if _, validErr := transaction.IsValid(); validErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(validErr.Error()))
		return
	}
	// INSERT To Database
	insertRow := db.Create(&transaction)
	if insertRow.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(insertRow.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(transaction))
}
