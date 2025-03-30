package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Mitsui515/finsys/model"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) Delete(id string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *UserRepositoryImpl) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? AND is_deleted = ?", username, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByGitHubID(githubID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("github_id = ? AND is_deleted = ?", githubID, false).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) List(page, size int) ([]*model.User, int64, error) {
	var users []*model.User
	var count int64
	offset := (page - 1) * size
	err := r.db.Model(&model.User{}).Where("is_deleted = ?", false).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("is_deleted = ?", false).Offset(offset).Limit(size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
