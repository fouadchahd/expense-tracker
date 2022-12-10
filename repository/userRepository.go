package repository

import (
	"errors"
	"github.com/fouadelhamri/expense-tracker/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (ur *UserRepository) GetUserByID(ID uint) (*model.User, error) {
	u := model.User{
		ID: ID,
	}
	row := ur.db.Preload("Meta").First(&u)
	if row.Error != nil {
		return nil, row.Error
	}
	if row.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
func (ur *UserRepository) GetUserByCredentials(ID uint, pass string) (*model.User, error) {
	var u model.User
	row := ur.db.Preload("Meta").Where(&model.User{Password: pass, ID: ID}).First(&u)
	if row.Error != nil {
		return nil, row.Error
	}
	if row.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	var u model.User
	row := ur.db.Preload("Meta").Where(&model.User{
		Username: username,
	}).First(&u)
	if row.Error != nil {
		return nil, row.Error
	}
	if row.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func (ur *UserRepository) CreateUser(username string, password string, role string, meta map[string]any) (uint, error) {
	//Create Model object
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return 0, errors.New("trouble hashing the password provided")
	}
	u := model.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
		Meta:     make([]model.MetaItem, 0),
	}
	for key, val := range meta {
		metaItem := model.MetaItem{
			Name:  key,
			Value: val.(string),
		}
		u.Meta = append(u.Meta, metaItem)
	}
	insertROW := ur.db.Create(&u)

	//Error
	if insertROW.Error != nil {
		return 0, insertROW.Error
	}

	//inserted
	return u.ID, nil

}

func (ur *UserRepository) ConvertUserModelToMapResponse(userModal *model.User) map[string]any {
	//Convert Response
	userMap := make(map[string]any, 0)
	userMap["id"] = userModal.ID
	userMap["username"] = userModal.Username
	userMap["role"] = userModal.Role
	userMap["password"] = userModal.Password
	userMap["group_id"] = nil
	if userModal.GroupID != nil {
		userMap["group_id"] = userModal.GroupID
	}
	meta := make(map[string]string, 0)
	for _, metaItem := range userModal.Meta {
		meta[metaItem.Name] = metaItem.Value
	}
	userMap["meta"] = meta
	return userMap
}
