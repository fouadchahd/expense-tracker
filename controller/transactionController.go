package controller

import (
	"github.com/fouadelhamri/expense-tracker/controller/actions"
	"gorm.io/gorm"
	"net/http"
)

type TransactionController struct {
	db *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{db: db}
}

func (tc *TransactionController) CreateTransaction(res http.ResponseWriter, req *http.Request) {
	actions.SubmitTransactionAction(tc.db, res, req)
}

func (tc *TransactionController) GetTransactionsByUser(res http.ResponseWriter, req *http.Request) {
	actions.GetTransactionsListByUserID(tc.db, res, req)
}

func (tc *TransactionController) GetTransactionGroup(res http.ResponseWriter, req *http.Request) {
	actions.GetTransactionsListByGroupCodeAction(tc.db, res, req)
}
