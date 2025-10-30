package usecases

import (
	"context"

	"github.com/sainaif/animalsys/internal/core/entities"
	"github.com/sainaif/animalsys/internal/core/interfaces"
)

type FinanceUseCase struct {
	financeRepo interfaces.FinanceRepository
	auditRepo   interfaces.AuditLogRepository
}

func NewFinanceUseCase(
	financeRepo interfaces.FinanceRepository,
	auditRepo interfaces.AuditLogRepository,
) *FinanceUseCase {
	return &FinanceUseCase{
		financeRepo: financeRepo,
		auditRepo:   auditRepo,
	}
}

func (uc *FinanceUseCase) Create(ctx context.Context, req *entities.FinanceCreateRequest, createdBy string) (*entities.FinanceTransaction, error) {
	transaction := entities.NewFinanceTransaction(
		req.Type,
		req.Category,
		req.Amount,
		req.Date,
		req.Description,
	)

	transaction.Source = req.Source
	transaction.Reference = req.Reference
	transaction.FiscalYear = req.FiscalYear
	transaction.Notes = req.Notes
	transaction.CreatedBy = createdBy

	if err := uc.financeRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(createdBy, "", "", entities.ActionCreate, "finance", transaction.ID.Hex(), "Financial transaction recorded")
	uc.auditRepo.Create(ctx, auditLog)

	return transaction, nil
}

func (uc *FinanceUseCase) GetByID(ctx context.Context, id string) (*entities.FinanceTransaction, error) {
	return uc.financeRepo.GetByID(ctx, id)
}

func (uc *FinanceUseCase) List(ctx context.Context, filter *entities.FinanceFilter) ([]*entities.FinanceTransaction, int64, error) {
	return uc.financeRepo.List(ctx, filter)
}

func (uc *FinanceUseCase) Update(ctx context.Context, id string, req *entities.FinanceUpdateRequest, updatedBy string) (*entities.FinanceTransaction, error) {
	transaction, err := uc.financeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Type != "" {
		transaction.Type = req.Type
	}
	if req.Category != "" {
		transaction.Category = req.Category
	}
	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	if req.Date != "" {
		transaction.Date = req.Date
	}
	if req.Description != "" {
		transaction.Description = req.Description
	}
	if req.Source != "" {
		transaction.Source = req.Source
	}
	if req.Reference != "" {
		transaction.Reference = req.Reference
	}
	if req.FiscalYear != "" {
		transaction.FiscalYear = req.FiscalYear
	}
	if req.Notes != "" {
		transaction.Notes = req.Notes
	}
	transaction.UpdatedBy = updatedBy

	if err := uc.financeRepo.Update(ctx, id, transaction); err != nil {
		return nil, err
	}

	// Audit
	auditLog := entities.NewAuditLog(updatedBy, "", "", entities.ActionUpdate, "finance", id, "Financial transaction updated")
	uc.auditRepo.Create(ctx, auditLog)

	return transaction, nil
}

func (uc *FinanceUseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	if err := uc.financeRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit
	auditLog := entities.NewAuditLog(deletedBy, "", "", entities.ActionDelete, "finance", id, "Financial transaction deleted")
	uc.auditRepo.Create(ctx, auditLog)

	return nil
}

func (uc *FinanceUseCase) GetByDateRange(ctx context.Context, startDate, endDate string, limit, offset int) ([]*entities.FinanceTransaction, int64, error) {
	return uc.financeRepo.GetByDateRange(ctx, startDate, endDate, limit, offset)
}

func (uc *FinanceUseCase) GetByCategory(ctx context.Context, category entities.FinanceCategory, limit, offset int) ([]*entities.FinanceTransaction, int64, error) {
	return uc.financeRepo.GetByCategory(ctx, category, limit, offset)
}

func (uc *FinanceUseCase) GetSummary(ctx context.Context, startDate, endDate string) (map[string]float64, error) {
	return uc.financeRepo.GetSummary(ctx, startDate, endDate)
}

func (uc *FinanceUseCase) GetByFiscalYear(ctx context.Context, fiscalYear string) ([]*entities.FinanceTransaction, error) {
	return uc.financeRepo.GetByFiscalYear(ctx, fiscalYear)
}

func (uc *FinanceUseCase) GetFinancialReport(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	// Get all transactions in date range
	transactions, _, err := uc.financeRepo.GetByDateRange(ctx, startDate, endDate, 0, 0)
	if err != nil {
		return nil, err
	}

	// Get summary
	summary, err := uc.financeRepo.GetSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate detailed statistics
	totalIncome := 0.0
	totalExpense := 0.0
	incomeByCategory := make(map[string]float64)
	expenseByCategory := make(map[string]float64)

	for _, tx := range transactions {
		if tx.Type == entities.TransactionTypeIncome {
			totalIncome += tx.Amount
			incomeByCategory[string(tx.Category)] += tx.Amount
		} else {
			totalExpense += tx.Amount
			expenseByCategory[string(tx.Category)] += tx.Amount
		}
	}

	report := map[string]interface{}{
		"period": map[string]string{
			"start": startDate,
			"end":   endDate,
		},
		"summary": map[string]float64{
			"total_income":  totalIncome,
			"total_expense": totalExpense,
			"net_balance":   totalIncome - totalExpense,
		},
		"income_by_category":  incomeByCategory,
		"expense_by_category": expenseByCategory,
		"transaction_count":   len(transactions),
		"raw_summary":         summary,
	}

	return report, nil
}
