package service

import (
	"joi-delivery-golang/internal/dto/request"
	"joi-delivery-golang/internal/dto/response"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (cs *OrderService) PlaceOrder(req request.OrderPlaceRequest) (*response.OrderPlaceResponse, error) {
	userCart := usersCart[req.UserID]

	res := response.OrderPlaceResponse{
		Cart: *userCart,
	}
	return &res, nil
}
