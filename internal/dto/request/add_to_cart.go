package request

type AddToCartRequest struct {
	UserID    string `json:"userId"`
	OutletID  string `json:"outletId"`
	ProductID string `json:"productId"`
	Qty       int    `json:"qty"`
}
