package server

import (
	"context"

	"joi-delivery-golang/cmd/api"
	"joi-delivery-golang/cmd/api/handler"
	"joi-delivery-golang/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server   *echo.Echo
	handlers handler.Handler
}

func NewServer(ctx context.Context) *Server {
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.CORS())

	s := &Server{
		server:   server,
		handlers: registerHandlers(),
	}

	api.BindRoutes(server, s.handlers)

	service.InitializeTestData()

	return s
}

func (s *Server) Start(port string) error {
	return s.server.Start(port)
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func registerHandlers() handler.Handler {
	userService := service.NewUserService()
	productService := service.NewProductService()
	cartService := service.NewCartService(userService, productService)
	orderService := service.NewOrderService()

	cartHandler := handler.NewCartHandler(cartService)
	inventHandler := handler.NewInventoryHandler()
	orderPlaceHandler := handler.NewOrderHandler(orderService, cartService)

	return handler.Handler{
		CartHandler:      cartHandler,
		InventoryHandler: inventHandler,
		OrderHandler:     orderPlaceHandler,
	}
}
