package controller

import (
	"github.com/fouadelhamri/expense-tracker/controller/actions"
	"gorm.io/gorm"
	"net/http"
)

type CategoryController struct {
	db *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{
		db,
	}
}

func (cc *CategoryController) CreateCategory(res http.ResponseWriter, req *http.Request) {
	actions.CreateCategoryAction(cc.db, res, req)
}

func (cc *CategoryController) GetCategories(res http.ResponseWriter, req *http.Request) {
	actions.GetAllCategoriesAction(cc.db, res, req)
}
