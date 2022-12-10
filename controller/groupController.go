package controller

import (
	"github.com/fouadelhamri/expense-tracker/controller/actions"
	"gorm.io/gorm"
	"net/http"
)

type GroupController struct {
	db *gorm.DB
}

func NewGroupController(db *gorm.DB) *GroupController {
	return &GroupController{db: db}
}

func (gc *GroupController) CreateGroupByUser(res http.ResponseWriter, req *http.Request) {
	actions.CreateGroupAction(gc.db, res, req)
}

func (gc *GroupController) JoinGroup(res http.ResponseWriter, req *http.Request) {
	actions.JoinGroupAction(gc.db, res, req)
}

func (gc *GroupController) LeftGroup(res http.ResponseWriter, req *http.Request) {
	actions.LeftGroupAction(gc.db, res, req)
}
