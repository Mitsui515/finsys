package model

import "time"

type FraudReport struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	TransactionID uint        `json:"transaction_id" gorm:"index;not null"`
	Transaction   Transaction `json:"-" gorm:"foreignKey:TransactionID"`
	Report        string      `json:"report" gorm:"type:text;not null"`
	GeneratedAt   time.Time   `json:"generated_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     time.Time   `json:"deleted_at"`
}

func (FraudReport) TableName() string {
	return "fraud_reports"
}
