package controller

import (
	"github.com/fouadelhamri/expense-tracker/controller/actions"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db,
	}
}

func (uc *UserController) RegisterController(res http.ResponseWriter, req *http.Request) {
	actions.RegisterAction(uc.db, res, req)
}

func (uc *UserController) GetSingleUser(res http.ResponseWriter, req *http.Request) {
	actions.FetchUserByID(uc.db, res, req)
}
func (uc *UserController) CheckAuthorization(res http.ResponseWriter, req *http.Request) {
	actions.AuthorizationAction(uc.db, res, req)
}

func (uc *UserController) LoginController(res http.ResponseWriter, req *http.Request) {
	actions.LoginAction(uc.db, res, req)
}
