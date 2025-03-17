package service

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
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

func (s *TransactionService) GetByID(id uint) (*TransactionResponse, error) {
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

func (s *TransactionService) ListByPage(page, size int, transactionType string, startTime, endTime *time.Time) (*TransactionListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	offset := (page - 1) * size
	var transactions []model.Transaction
	var total int64
	db := config.DB.Model(&model.Transaction{}).Where("is_deleted = ?", false)
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
		return nil, err
	}
	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = TransactionResponse{
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
		}
	}
	return &TransactionListResponse{
		Total:        int(total),
		Page:         page,
		Size:         size,
		Transactions: responses,
	}, nil
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

func (s *TransactionService) ImportFromCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		return 0, err
	}
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}
	requiredColumns := []string{"type", "amount", "nameOrig", "oldBalanceOrig", "newBalanceOrig", "nameDest", "oldBalanceDest", "newBalanceDest", "isFraud"}
	for _, col := range requiredColumns {
		if _, exists := headerMap[col]; !exists {
			return 0, errors.New("CSV format error, there is no column " + col)
		}
	}
	var transactions []model.Transaction
	batchSize := 90
	totalImported := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return totalImported, err
		}
		amount, err := strconv.ParseFloat(record[headerMap["amount"]], 64)
		if err != nil {
			continue
		}
		oldBalanceOrig, _ := strconv.ParseFloat(record[headerMap["oldBalanceOrig"]], 64)
		newBalanceOrig, _ := strconv.ParseFloat(record[headerMap["newBalanceOrig"]], 64)
		oldBalanceDest, _ := strconv.ParseFloat(record[headerMap["oldBalanceDest"]], 64)
		newBalanceDest, _ := strconv.ParseFloat(record[headerMap["newBalanceDest"]], 64)
		transaction := model.Transaction{
			Type:           record[headerMap["type"]],
			Amount:         amount,
			NameOrig:       record[headerMap["nameOrig"]],
			OldBalanceOrig: oldBalanceOrig,
			NewBalanceOrig: newBalanceOrig,
			NameDest:       record[headerMap["nameDest"]],
			OldBalanceDest: oldBalanceDest,
			NewBalanceDest: newBalanceDest,
			IsFraud:        record[headerMap["isFraud"]] == "1",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		transactions = append(transactions, transaction)
		if len(transactions) >= batchSize {
			if err := config.DB.CreateInBatches(transactions, batchSize).Error; err != nil {
				return totalImported, err
			}
			totalImported += len(transactions)
			transactions = transactions[:0]
		}
	}
	if len(transactions) > 0 {
		if err := config.DB.CreateInBatches(transactions, len(transactions)).Error; err != nil {
			return totalImported, err
		}
		totalImported += len(transactions)
	}
	return totalImported, nil
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
