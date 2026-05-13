package api

import (
	"joi-delivery-golang/cmd/api/handler"

	"github.com/labstack/echo/v4"
)

func BindRoutes(server *echo.Echo, hand handler.Handler) {
	// Cart routes
	server.GET("/cart/view", hand.GetCart)
	server.POST("/cart/product", hand.AddToCart)

	// Inventory routes
	server.POST("/inventory/health", hand.InventoryHealth)

	//Order route
	server.POST("/order", hand.OrderPlace)
}
