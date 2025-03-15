package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"size:20;not null;uniqueIndex"`
	Password  string `json:"password" gorm:"size:100;not null"`
	Email     string `json:"email" gorm:"size:100;not null;uniqueIndex"`
	IsAdmin   bool   `json:"isAdmin" gorm:"default:false"`
	IsDeleted bool   `json:"isDeleted" gorm:"default:false"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
	DeletedAt int64  `json:"deletedAt" gorm:"index"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if len(u.Password) > 60 {
		return nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
