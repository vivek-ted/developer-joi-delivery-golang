package request

type OrderPlaceRequest struct {
	UserID string `json:"userId"`
	CartID string `json:"cartId"`
}
