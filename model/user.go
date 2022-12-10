package model

type User struct {
	ID           uint          `gorm:"primaryKey"`
	Username     string        `json:"username" gorm:"unique;size:255"`
	Password     string        `json:"password"`
	Role         string        `json:"role"`
	Transactions []Transaction `json:"-" gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Meta         []MetaItem    `json:"meta" gorm:"foreignKey:UserUsername;references:Username;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GroupID      *uint         `json:"group_id"`
}

type MetaItem struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	UserUsername string `gorm:"size:255"`
}

func (m *MetaItem) TableName() string {
	return "meta"
}
