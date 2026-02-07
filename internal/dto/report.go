package dto

type ReportRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReportResponse struct {
	TotalTransactions int           `json:"total_transactions"`
	TotalRevenue      int           `json:"total_revenue"`
	ProductBestSeller *ProductSales `json:"product_best_seller"`
}

type ProductSales struct {
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}
