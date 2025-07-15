package repository

import (
	"errors"
	"time"

	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type FraudReportRepositoryImpl struct {
	db *gorm.DB
}

type MongoReport struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FraudReportID uint               `bson:"fraud_report_id"`
	Content       string             `bson:"content"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func NewFraudReportRepository(db *gorm.DB) FraudReportRepository {
	return &FraudReportRepositoryImpl{
		db: db,
	}
}

func (r *FraudReportRepositoryImpl) Create(report *model.FraudReport) error {
	tx := r.db.Begin()
	report.GeneratedAt = time.Now()
	report.UpdatedAt = time.Now()
	if err := tx.Create(report).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *FraudReportRepositoryImpl) Update(report *model.FraudReport) error {
	_, err := r.FindByID(report.ID)
	if err != nil {
		return err
	}
	tx := r.db.Begin()
	report.UpdatedAt = time.Now()
	if err := tx.Save(report).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *FraudReportRepositoryImpl) Delete(id uint) error {
	tx := r.db.Begin()
	now := time.Now()
	if err := tx.Model(&model.FraudReport{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *FraudReportRepositoryImpl) FindByID(id uint) (*model.FraudReport, error) {
	var report model.FraudReport
	err := r.db.Where("id = ? AND deleted_at = ?", id, time.Time{}).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrFraudReportNotExists
		}
		return nil, err
	}
	return &report, nil
}

func (r *FraudReportRepositoryImpl) FindByTransactionID(transactionID uint) (*model.FraudReport, error) {
	var report model.FraudReport
	err := r.db.Where("transaction_id = ? AND deleted_at = ?", transactionID, time.Time{}).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrFraudReportNotExists
		}
		return nil, err
	}
	return &report, nil
}

func (r *FraudReportRepositoryImpl) List(page, size int) ([]*model.FraudReport, int64, error) {
	var reports []*model.FraudReport
	var count int64
	offset := (page - 1) * size
	err := r.db.Model(&model.FraudReport{}).Where("deleted_at = ?", time.Time{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Where("deleted_at = ?", time.Time{}).Order("generated_at DESC").Offset(offset).Limit(size).Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}
	return reports, count, nil
}
