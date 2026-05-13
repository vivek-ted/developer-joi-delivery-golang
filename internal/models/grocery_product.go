package models

type GroceryProduct struct {
	Product
	ExpiryDate     int     `json:"expiryDate"`
	Threshold      int     `json:"threshold"`
	AvailableStock int     `json:"availableStock"`
	SellingPrice   float64 `json:"sellingPrice"`
	Weight         float64 `json:"weight"`
	Discount       float64 `json:"discount"`
	Store          Store   `json:"store"`
	Qty            int     `json:"qty"`
	Category       string  `json:"category"`
}
