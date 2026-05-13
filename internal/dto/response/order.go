package response

import "joi-delivery-golang/internal/models"

type OrderPlaceResponse struct {
	Cart models.Cart `json:"cart"`
}
