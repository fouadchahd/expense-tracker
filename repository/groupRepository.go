package repository

import (
	"errors"
	"github.com/fouadelhamri/expense-tracker/model"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		db,
	}
}

func (gr *GroupRepository) GetGroupByCode(groupCode string) (*model.Group, error) {
	group := &model.Group{}
	fetchRow := gr.db.Preload("Users").Where("code = ?", groupCode).First(&group)
	if fetchRow.Error != nil {
		return nil, errors.New("group not found")
	}
	if fetchRow.RowsAffected == 1 {
		return group, nil
	}
	return nil, errors.New("group not found")
}
func (gr *GroupRepository) GetGroupTransactionsByCode(groupCode string) (*model.Group, error) {
	group := &model.Group{}
	fetchRow := gr.db.Preload("Users").Preload("Transactions").Where("code = ?", groupCode).First(&group)
	if fetchRow.Error != nil {
		return nil, errors.New("group not found")
	}
	if fetchRow.RowsAffected == 1 {
		return group, nil
	}
	return nil, errors.New("group not found")
}
