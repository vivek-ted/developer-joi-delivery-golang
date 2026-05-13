package handler

import (
	"net/http"

	"joi-delivery-golang/internal/dto/request"
	"joi-delivery-golang/internal/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService *service.OrderService
	cartService  *service.CartService
}

func NewOrderHandler(orderService *service.OrderService, cartService *service.CartService) OrderHandler {
	return OrderHandler{
		orderService: orderService,
		cartService:  cartService,
	}
}

func (ch *OrderHandler) OrderPlace(c echo.Context) error {
	var req request.OrderPlaceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Success": false,
			"Message": "Invalid request format",
		})
	}

	if req.UserID == "" || req.CartID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"Success": false,
			"Message": "Missing required fields",
		})
	}

	resp, err := ch.orderService.PlaceOrder(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"Success": false,
			"Message": "Internal server error",
		})
	}

	return c.JSON(http.StatusOK, resp)
}
