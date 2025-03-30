package repository

import (
	"github.com/Mitsui515/finsys/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	List(page, size int) ([]*model.User, int64, error)
}
