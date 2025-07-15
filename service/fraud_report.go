package service

import (
	"errors"
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
	TransactionID uint   `json:"transaction_id"`
	Report        string `json:"report"`
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
	if req.Report == "" {
		return 0, utils.ErrInvalidReport
	}
	_, err := s.transactionRepository.FindByID(req.TransactionID)
	if err != nil {
		return 0, err
	}
	report := &model.FraudReport{
		TransactionID: req.TransactionID,
		Report:        req.Report,
		GeneratedAt:   time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.fraudReportRepository.Create(report); err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (s *FraudReportService) Update(id uint, req *FraudReportRequest) (*FraudReportResponse, error) {
	if req.Report == "" {
		return nil, utils.ErrInvalidReport
	}
	report, err := s.fraudReportRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	report.Report = req.Report
	report.UpdatedAt = time.Now()
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

// TODO
func generateFraudAnalysisReport(transaction *model.Transaction) string {
	var result string
	result = "# 交易欺诈分析报告\n\n"
	result += "## 交易信息\n\n"
	result += "- 交易类型: " + transaction.Type + "\n"
	result += "- 交易金额: " + formatFloat(transaction.Amount) + "\n"
	result += "- 发起方: " + transaction.NameOrig + "\n"
	result += "- 接收方: " + transaction.NameDest + "\n\n"
	result += "## 欺诈风险分析\n\n"
	if transaction.Amount > 100000 {
		result += "- **高风险**: 交易金额异常大 (" + formatFloat(transaction.Amount) + ")\n"
	}
	origBalanceDiff := transaction.OldBalanceOrig - transaction.NewBalanceOrig
	if origBalanceDiff != transaction.Amount {
		result += "- **异常**: 发起方余额变化 (" + formatFloat(origBalanceDiff) + ") 与交易金额不符\n"
	}
	destBalanceDiff := transaction.NewBalanceDest - transaction.OldBalanceDest
	if destBalanceDiff != transaction.Amount {
		result += "- **异常**: 接收方余额变化 (" + formatFloat(destBalanceDiff) + ") 与交易金额不符\n"
	}
	if transaction.NewBalanceOrig < 0 {
		result += "- **异常**: 发起方新余额为负数\n"
	}
	if transaction.IsFraud {
		result += "\n## 结论\n\n此交易被系统标记为**欺诈交易**，欺诈概率为: " + formatFloat(transaction.FraudProbability*100) + "%\n"
	} else {
		result += "\n## 结论\n\n此交易被系统判定为**正常交易**，欺诈概率为: " + formatFloat(transaction.FraudProbability*100) + "%\n"
	}
	result += "\n生成时间: " + time.Now().Format(time.RFC3339) + "\n"
	return result
}

func formatFloat(num float64) string {
	return time.Now().Format(time.RFC3339)
}
