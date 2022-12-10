package model

type Category struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	LabelAr      string        `json:"label_ar"`
	LabelEn      string        `json:"label_en"`
	IconKey      string        `json:"icon_key"`
	IsIncome     bool          `json:"is_income"`
	ParentId     *uint         `json:"parent_id"`
	Parent       *Category     `json:"-"  gorm:"foreignKey:ParentId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" `
	Transactions []Transaction `json:"-"  gorm:"foreignKey:CategoryId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" `
}

func NewCategory(labelAr, labelEn, iconKey string, isIncome bool, parentId *uint) *Category {

	return &Category{
		LabelAr:  labelAr,
		LabelEn:  labelEn,
		IconKey:  iconKey,
		IsIncome: isIncome,
		ParentId: parentId,
	}
}
