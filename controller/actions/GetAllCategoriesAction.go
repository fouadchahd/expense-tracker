package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"gorm.io/gorm"
	"net/http"
)

func GetAllCategoriesAction(db *gorm.DB, res http.ResponseWriter, _ *http.Request) {
	categories := make([]model.Category, 0)
	row := db.Model(&model.Category{}).Find(&categories)

	if row.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(row.Error.Error()))
		return
	}
	json.NewEncoder(res).Encode(NewApiResponseSuccess(categories))

}
