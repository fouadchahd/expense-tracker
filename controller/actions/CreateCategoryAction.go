package actions

import (
	"encoding/json"
	"github.com/fouadelhamri/expense-tracker/model"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func CreateCategoryAction(db *gorm.DB, res http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	var body = make(map[string]any, 0)
	decoderErr := json.NewDecoder(req.Body).Decode(&body)
	if decoderErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		if decoderErr == io.EOF {
			json.NewEncoder(res).Encode(NewApiResponseError("no request body provided"))
			return
		}
		json.NewEncoder(res).Encode(NewApiResponseError(decoderErr.Error()))
		return
	}
	_, ok1 := body["label_ar"]
	_, ok2 := body["label_en"]
	_, ok3 := body["icon_key"]
	_, ok4 := body["parent_id"]
	_, ok5 := body["is_income"]
	if !(ok1 && ok2 && ok3 && ok4 && ok5) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError("Missing Properties : "))
		return
	}

	var newCat model.Category
	bytes, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(marshalErr.Error()))
		return
	}

	unmarshalErr := json.Unmarshal(bytes, &newCat)
	if unmarshalErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(unmarshalErr.Error()))
		return
	}
	insertRow := db.Create(&newCat)

	if insertRow.Error != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(NewApiResponseError(insertRow.Error.Error()))
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(NewApiResponseSuccess(newCat))

	//var cat model.Category
	//bytes, _ := json.Marshal(body)
	//json.Unmarshal(bytes, &cat)

}
