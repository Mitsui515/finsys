package repository

import (
	"github.com/Mitsui515/finsys/model"
)

type FraudReportRepository interface {
	Create(report *model.FraudReport) error
	Update(report *model.FraudReport) error
	Delete(id uint) error
	FindByID(id uint) (*model.FraudReport, error)
	FindByTransactionID(transactionID uint) (*model.FraudReport, error)
	List(page, size int) ([]*model.FraudReport, int64, error)
}
