package repository

import (
	"errors"
	"time"

	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/utils"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepositoryImpl) Update(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepositoryImpl) Delete(id uint) error {
	return r.db.Model(&model.Transaction{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *TransactionRepositoryImpl) FindByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrTransactionNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepositoryImpl) List(page, size int, transactionType string, startTime, endTime *time.Time) ([]*model.Transaction, int64, error) {
	var transactions []*model.Transaction
	var total int64
	db := r.db.Model(&model.Transaction{}).Where("is_deleted = ?", false)
	if transactionType != "" {
		db = db.Where("type = ?", transactionType)
	}
	if startTime != nil {
		db = db.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		db = db.Where("created_at <= ?", endTime)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}
