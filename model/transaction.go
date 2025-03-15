package model

import "time"

type Transaction struct {
	ID               uint      `json:"id" gorm:"primary_key"`
	Type             string    `json:"type" gorm:"size:50;not null;index"`
	Amount           float64   `json:"amount" gorm:"not null"`
	NameOrig         string    `json:"nameOrig" gorm:"size:50;not null;index"`
	OldBalanceOrig   float64   `json:"oldBalanceOrig" gorm:"not null"`
	NewBalanceOrig   float64   `json:"newBalanceOrig" gorm:"not null"`
	NameDest         string    `json:"nameDest" gorm:"size:50;not null;index"`
	OldBalanceDest   float64   `json:"oldBalanceDest" gorm:"not null"`
	NewBalanceDest   float64   `json:"newBalanceDest" gorm:"not null"`
	IsFraud          bool      `json:"isFraud" gorm:"default:false"`
	FraudProbability float64   `json:"fraudProbability" gorm:"default:0"`
	IsDeleted        bool      `json:"isDeleted" gorm:"default:false"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	DeletedAt        time.Time `json:"deletedAt" gorm:"index"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}
