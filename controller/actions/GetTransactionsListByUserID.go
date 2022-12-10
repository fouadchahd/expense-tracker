package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetTransactionsListByUserID(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()

	dayStr := queries.Get("day")
	monthStr := queries.Get("month")
	yearStr := queries.Get("year")
	d, _ := strconv.Atoi(dayStr)
	m, _ := strconv.Atoi(monthStr)
	y, _ := strconv.Atoi(yearStr)

	//Start checking auth
	key := req.Header.Get("key")
	token := req.Header.Get("token")
	if len(key) == 0 || len(token) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("missing headers"))
		return
	}
	subRoute := mux.CurrentRoute(req).Subrouter()
	url, err := subRoute.Get("check_authorization").URL("id", key, "password", token)
	subRoute.Methods(http.MethodGet)
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
	userID, _ := strconv.Atoi(key)
	var transactions = make([]model.Transaction, 0)

	dbTrans := db.Where(model.Transaction{UserId: uint(userID)})
	// Set Filters
	if y > 0 && y < time.Now().UTC().Year() {
		dbTrans.Where("year(submitted_at) = ?", y)
	}
	if m > 0 && m <= 12 {
		dbTrans.Where("month(submitted_at) = ?", m)
	}
	if d > 0 {
		dbTrans.Where("day(submitted_at) = ?", d)
	}

	fetchRow := dbTrans.Find(&transactions)
	if fetchRow.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(fetchRow.Error.Error()))
		return
	}
	json.NewEncoder(res).Encode(NewApiResponseSuccess(transactions))
	return

}
