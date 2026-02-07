package services

import (
	"basic-rest-api-go/internal/dto"
	"basic-rest-api-go/internal/repositories"
)

type ReportService struct {
	repo        *repositories.TransactionRepository
	productRepo *repositories.ProductRepository
}

func NewReportService(repo *repositories.TransactionRepository, productRepo *repositories.ProductRepository) *ReportService {
	return &ReportService{repo: repo, productRepo: productRepo}
}

func (s *ReportService) GenerateSalesReport(request dto.ReportRequest) (*dto.ReportResponse, error) {
	// get all transactions within date range
	trxs, err := s.repo.GetTransactionsByDateRange(request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	if len(trxs) == 0 {
		return &dto.ReportResponse{
			TotalTransactions: 0,
			TotalRevenue:      0,
			ProductBestSeller: nil,
		}, nil
	}

	// get transaction details for each transaction
	var transactionIDs []int
	totalRevenue := 0
	for _, trx := range trxs {
		transactionIDs = append(transactionIDs, trx.ID)
		totalRevenue += trx.TotalAmount
	}

	// get best selling product
	bestSeller, err := s.GetBestSellingProduct(transactionIDs)
	if err != nil {
		return nil, err
	}

	return &dto.ReportResponse{
		TotalTransactions: len(trxs),
		TotalRevenue:      totalRevenue,
		ProductBestSeller: bestSeller,
	}, nil
}

func (s *ReportService) GetBestSellingProduct(transactionIDs []int) (*dto.ProductSales, error) {
	details, err := s.repo.GetTransactionDetails(transactionIDs)
	if err != nil {
		return nil, err
	}
	productSalesMap := make(map[int]int)

	for _, detail := range details {
		productSalesMap[detail.ProductID] += detail.Quantity
	}

	var bestSellerID int
	var bestSellerCount int
	for productID, count := range productSalesMap {
		if count > bestSellerCount {
			bestSellerCount = count
			bestSellerID = productID
		}
	}

	// get product name
	product, err := s.productRepo.GetByID(bestSellerID)
	if err != nil {
		return nil, err
	}

	return &dto.ProductSales{
		ProductName: product.Name,
		TotalSold:   bestSellerCount,
	}, nil
}
