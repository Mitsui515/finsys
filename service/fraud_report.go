package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Mitsui515/finsys/model"
	"github.com/Mitsui515/finsys/repository"
	"github.com/Mitsui515/finsys/utils"
	"gorm.io/gorm"
)

type FraudReportService struct {
	fraudReportRepository repository.FraudReportRepository
	transactionRepository repository.TransactionRepository
}

func NewFraudReportService(db *gorm.DB) *FraudReportService {
	return &FraudReportService{
		fraudReportRepository: repository.NewFraudReportRepository(db),
		transactionRepository: repository.NewTransactionRepository(db),
	}
}

type FraudReportRequest struct {
	TransactionID uint `json:"transaction_id"`
}

type FraudReportResponse struct {
	ID            uint      `json:"id"`
	TransactionID uint      `json:"transaction_id"`
	Report        string    `json:"report"`
	GeneratedAt   time.Time `json:"generated_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FraudReportListResponse struct {
	Total   int64                 `json:"total"`
	Page    int                   `json:"page"`
	Size    int                   `json:"size"`
	Reports []FraudReportResponse `json:"reports"`
}

func (s *FraudReportService) Create(req *FraudReportRequest) (uint, error) {
	if req.TransactionID == 0 {
		return 0, utils.ErrInvalidTransactionID
	}
	transaction, err := s.transactionRepository.FindByID(req.TransactionID)
	if err != nil {
		return 0, err
	}
	_, err = s.fraudReportRepository.FindByTransactionID(req.TransactionID)
	if err == nil {
		return 0, errors.New("fraud report for this transaction already exists")
	}
	reportContent := generateFraudAnalysisReport(transaction)
	report := &model.FraudReport{
		TransactionID: req.TransactionID,
		Report:        reportContent,
	}
	if err := s.fraudReportRepository.Create(report); err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (s *FraudReportService) Update(id uint) (*FraudReportResponse, error) {
	report, err := s.fraudReportRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	transaction, err := s.transactionRepository.FindByID(report.TransactionID)
	if err != nil {
		return nil, err
	}
	report.Report = generateFraudAnalysisReport(transaction)
	if err := s.fraudReportRepository.Update(report); err != nil {
		return nil, err
	}
	return &FraudReportResponse{
		ID:            report.ID,
		TransactionID: report.TransactionID,
		Report:        report.Report,
		GeneratedAt:   report.GeneratedAt,
		UpdatedAt:     report.UpdatedAt,
	}, nil
}

func (s *FraudReportService) Delete(id uint) error {
	return s.fraudReportRepository.Delete(id)
}

func (s *FraudReportService) GetByID(id uint) (*FraudReportResponse, error) {
	report, err := s.fraudReportRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &FraudReportResponse{
		ID:            report.ID,
		TransactionID: report.TransactionID,
		Report:        report.Report,
		GeneratedAt:   report.GeneratedAt,
		UpdatedAt:     report.UpdatedAt,
	}, nil
}

func (s *FraudReportService) GetByTransactionID(transactionID uint) (*FraudReportResponse, error) {
	report, err := s.fraudReportRepository.FindByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}
	return &FraudReportResponse{
		ID:            report.ID,
		TransactionID: report.TransactionID,
		Report:        report.Report,
		GeneratedAt:   report.GeneratedAt,
		UpdatedAt:     report.UpdatedAt,
	}, nil
}

func (s *FraudReportService) List(page, size int) (*FraudReportListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	reports, total, err := s.fraudReportRepository.List(page, size)
	if err != nil {
		return nil, err
	}
	response := &FraudReportListResponse{
		Total:   total,
		Page:    page,
		Size:    size,
		Reports: make([]FraudReportResponse, len(reports)),
	}
	for i, report := range reports {
		response.Reports[i] = FraudReportResponse{
			ID:            report.ID,
			TransactionID: report.TransactionID,
			Report:        report.Report,
			GeneratedAt:   report.GeneratedAt,
			UpdatedAt:     report.UpdatedAt,
		}
	}
	return response, nil
}

func (s *FraudReportService) GenerateReport(transactionID uint) (*FraudReportResponse, error) {
	transaction, err := s.transactionRepository.FindByID(transactionID)
	if err != nil {
		return nil, err
	}
	existingReport, err := s.fraudReportRepository.FindByTransactionID(transactionID)
	if err == nil {
		return &FraudReportResponse{
			ID:            existingReport.ID,
			TransactionID: existingReport.TransactionID,
			Report:        existingReport.Report,
			GeneratedAt:   existingReport.GeneratedAt,
			UpdatedAt:     existingReport.UpdatedAt,
		}, nil
	} else if !errors.Is(err, utils.ErrFraudReportNotExists) {
		return nil, err
	}
	reportContent := generateFraudAnalysisReport(transaction)
	report := &model.FraudReport{
		TransactionID: transactionID,
		Report:        reportContent,
		GeneratedAt:   time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.fraudReportRepository.Create(report); err != nil {
		return nil, err
	}
	return &FraudReportResponse{
		ID:            report.ID,
		TransactionID: report.TransactionID,
		Report:        report.Report,
		GeneratedAt:   report.GeneratedAt,
		UpdatedAt:     report.UpdatedAt,
	}, nil
}

func generateFraudAnalysisReport(transaction *model.Transaction) string {
	var result string
	result = "# Transaction Fraud Analysis Report\n\n"
	result += "## Transaction Information\n\n"
	result += "- Transaction Type: " + transaction.Type + "\n"
	result += "- Amount: " + formatFloat(transaction.Amount) + "\n"
	result += "- Originator: " + transaction.NameOrig + "\n"
	result += "- Destination: " + transaction.NameDest + "\n\n"
	result += "## Fraud Risk Analysis\n\n"
	if transaction.Amount > 100000 {
		result += "- **High Risk**: Unusually large transaction amount (" + formatFloat(transaction.Amount) + ")\n"
	}
	origBalanceDiff := transaction.OldBalanceOrig - transaction.NewBalanceOrig
	if origBalanceDiff != transaction.Amount {
		result += "- **Anomaly**: Originator balance change (" + formatFloat(origBalanceDiff) + ") does not match transaction amount.\n"
	}
	destBalanceDiff := transaction.NewBalanceDest - transaction.OldBalanceDest
	if destBalanceDiff != transaction.Amount {
		result += "- **Anomaly**: Destination balance change (" + formatFloat(destBalanceDiff) + ") does not match transaction amount.\n"
	}
	if transaction.NewBalanceOrig < 0 {
		result += "- **Anomaly**: Originator's new balance is negative.\n"
	}
	if transaction.IsFraud {
		result += "\n## Conclusion\n\nThis transaction is flagged as **FRAUDULENT** by the system, with a fraud probability of: " + formatFloat(transaction.FraudProbability*100) + "%\n"
	} else {
		result += "\n## Conclusion\n\nThis transaction is determined to be **NORMAL** by the system, with a fraud probability of: " + formatFloat(transaction.FraudProbability*100) + "%\n"
	}
	result += "\nGenerated At: " + time.Now().Format(time.RFC3339) + "\n"
	return result
}

func formatFloat(num float64) string {
	return fmt.Sprintf("%.2f", num)
}
