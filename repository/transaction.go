package repository

import (
	"time"

	"github.com/Mitsui515/finsys/model"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	Update(transaction *model.Transaction) error
	Delete(id uint) error
	FindByID(id uint) (*model.Transaction, error)
	List(page, size int, transactionType string, startTime, endTime *time.Time) ([]*model.Transaction, int64, error)
}
