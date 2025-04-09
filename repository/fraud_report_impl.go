package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type FraudReportRepositoryImpl struct {
	db          *gorm.DB
	mongoClient *mongo.Database
	collection  string
}

type MongoReport struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FraudReportID uint               `bson:"fraud_report_id"`
	Content       string             `bson:"content"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

func NewFraudReportRepository(db *gorm.DB, mongoClient *mongo.Database) FraudReportRepository {
	return &FraudReportRepositoryImpl{
		db:          db,
		mongoClient: mongoClient,
		collection:  "fraud_reports",
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoReport := MongoReport{
		FraudReportID: report.ID,
		Content:       report.Report,
		CreatedAt:     report.GeneratedAt,
		UpdatedAt:     report.UpdatedAt,
	}
	_, err := r.mongoClient.Collection(r.collection).InsertOne(ctx, mongoReport)
	if err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"fraud_report_id": report.ID}
	update := bson.M{
		"$set": bson.M{
			"content":    report.Report,
			"updated_at": report.UpdatedAt,
		},
	}
	_, err = r.mongoClient.Collection(r.collection).UpdateOne(ctx, filter, update)
	if err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.mongoClient.Collection(r.collection).DeleteOne(ctx, bson.M{"fraud_report_id": id})
	if err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var mongoReport MongoReport
	err = r.mongoClient.Collection(r.collection).FindOne(ctx, bson.M{"fraud_report_id": id}).Decode(&mongoReport)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFraudReportNotExists
		}
		return nil, err
	}
	report.Report = mongoReport.Content
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var mongoReport MongoReport
	err = r.mongoClient.Collection(r.collection).FindOne(ctx, bson.M{"fraud_report_id": report.ID}).Decode(&mongoReport)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrFraudReportNotExists
		}
		return nil, err
	}
	report.Report = mongoReport.Content
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
	if len(reports) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		reportIDs := make([]uint, len(reports))
		for i, report := range reports {
			reportIDs[i] = report.ID
		}
		cursor, err := r.mongoClient.Collection(r.collection).Find(ctx, bson.M{"fraud_report_id": bson.M{"$in": reportIDs}})
		if err != nil {
			return nil, 0, err
		}
		defer cursor.Close(ctx)
		mongoReports := make(map[uint]string)
		for cursor.Next(ctx) {
			var mongoReport MongoReport
			if err := cursor.Decode(&mongoReport); err != nil {
				return nil, 0, err
			}
			mongoReports[mongoReport.FraudReportID] = mongoReport.Content
		}
		for _, report := range reports {
			if content, exists := mongoReports[report.ID]; exists {
				report.Report = content
			}
		}
	}
	return reports, count, nil
}
