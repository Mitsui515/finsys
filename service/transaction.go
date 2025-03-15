package service

import (
	"errors"
	"time"

	"github.com/Mitsui515/finsys/config"
	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/utils"
	"gorm.io/gorm"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

type TransactionRequest struct {
	Type           string  `json:"type"`
	Amount         float64 `json:"amount"`
	NameOrig       string  `json:"nameOrig"`
	OldBalanceOrig float64 `json:"oldBalanceOrig"`
	NewBalanceOrig float64 `json:"newBalanceOrig"`
	NameDest       string  `json:"nameDest"`
	OldBalanceDest float64 `json:"oldBalanceDest"`
	NewBalanceDest float64 `json:"newBalanceDest"`
}

type TransactionResponse struct {
	ID             uint      `json:"id"`
	Type           string    `json:"type"`
	Amount         float64   `json:"amount"`
	NameOrig       string    `json:"nameOrig"`
	OldBalanceOrig float64   `json:"oldBalanceOrig"`
	NewBalanceOrig float64   `json:"newBalanceOrig"`
	NameDest       string    `json:"nameDest"`
	OldBalanceDest float64   `json:"oldBalanceDest"`
	NewBalanceDest float64   `json:"newBalanceDest"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	DeletedAt      time.Time `json:"deletedAt,omitempty"`
}

type TransactionListResponse struct {
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	Size         int                   `json:"size"`
	Transactions []TransactionResponse `json:"transactions"`
}

func (s *TransactionService) Create(req *TransactionRequest) (uint, error) {
	if err := validateTransaction(req); err != nil {
		return 0, err
	}
	transaction := model.Transaction{
		Type:           req.Type,
		Amount:         req.Amount,
		NameOrig:       req.NameOrig,
		OldBalanceOrig: req.OldBalanceOrig,
		NewBalanceOrig: req.NewBalanceOrig,
		NameDest:       req.NameDest,
		OldBalanceDest: req.OldBalanceDest,
		NewBalanceDest: req.NewBalanceDest,
		IsFraud:        false,
	}
	if err := config.DB.Create(&transaction).Error; err != nil {
		return 0, err
	}
	go s.predictFraud(&transaction)
	return transaction.ID, nil
}

func (*TransactionService) GetByID(id uint) (*TransactionResponse, error) {
	var transaction model.Transaction
	result := config.DB.First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrTransactionNotExists
		}
		return nil, result.Error
	}
	return &TransactionResponse{
		ID:             transaction.ID,
		Type:           transaction.Type,
		Amount:         transaction.Amount,
		NameOrig:       transaction.NameOrig,
		OldBalanceOrig: transaction.OldBalanceOrig,
		NewBalanceOrig: transaction.NewBalanceOrig,
		NameDest:       transaction.NameDest,
		OldBalanceDest: transaction.OldBalanceDest,
		NewBalanceDest: transaction.NewBalanceDest,
		CreatedAt:      transaction.CreatedAt,
	}, nil
}

func (s *TransactionService) Update(id uint, req *TransactionRequest) (*TransactionResponse, error) {
	if err := validateTransaction(req); err != nil {
		return nil, err
	}
	var transaction model.Transaction
	result := config.DB.First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrTransactionNotExists
		}
		return nil, result.Error
	}
	transaction.Type = req.Type
	transaction.Amount = req.Amount
	transaction.NameOrig = req.NameOrig
	transaction.OldBalanceOrig = req.OldBalanceOrig
	transaction.NewBalanceOrig = req.NewBalanceOrig
	transaction.NameDest = req.NameDest
	transaction.OldBalanceDest = req.OldBalanceDest
	transaction.NewBalanceDest = req.NewBalanceDest
	if err := config.DB.Save(&transaction).Error; err != nil {
		return nil, err
	}
	go s.predictFraud(&transaction)
	return &TransactionResponse{
		ID:             transaction.ID,
		Type:           transaction.Type,
		Amount:         transaction.Amount,
		NameOrig:       transaction.NameOrig,
		OldBalanceOrig: transaction.OldBalanceOrig,
		NewBalanceOrig: transaction.NewBalanceOrig,
		NameDest:       transaction.NameDest,
		OldBalanceDest: transaction.OldBalanceDest,
		NewBalanceDest: transaction.NewBalanceDest,
		CreatedAt:      transaction.CreatedAt,
		UpdatedAt:      transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) Delete(id uint) error {
	var transaction model.Transaction
	result := config.DB.First(&transaction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.ErrTransactionNotExists
		}
		return result.Error
	}
	transaction.IsDeleted = true
	return config.DB.Save(&transaction).Error
}

func validateTransaction(req *TransactionRequest) error {
	if req.Type == "" {
		return utils.ErrMissingType
	}
	if req.Amount <= 0 {
		return utils.ErrInvalidAmount
	}
	if req.NameOrig == "" {
		return utils.ErrInvalidOrig
	}
	if req.NameDest == "" {
		return utils.ErrInvalidDest
	}
	return nil
}

func (s *TransactionService) predictFraud(transaction *model.Transaction) {
	if transaction.Amount > 100000 {
		transaction.IsFraud = true
	} else {
		transaction.IsFraud = false
	}
	config.DB.Save(transaction)
}
