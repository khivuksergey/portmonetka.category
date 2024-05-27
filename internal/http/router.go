package http

import (
	"github.com/khivuksergey/portmonetka.category/config"
	"github.com/khivuksergey/portmonetka.category/docs"
	"github.com/khivuksergey/portmonetka.category/internal/core/port/service"
	"github.com/khivuksergey/webserver/logger"
	"github.com/khivuksergey/webserver/router"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Router struct {
	*echo.Echo
}

func NewRouter(cfg *config.Configuration, services *service.Manager, logger logger.Logger) http.Handler {
	handlers := newHandlers(services, logger)

	e := router.NewEchoRouter().
		WithConfig(cfg.Router).
		UseMiddleware(handlers.error.HandleError).
		UseHealthCheck().
		UseSwagger(docs.SwaggerInfo, cfg.Swagger)

	categories := e.Group("users/:userId/categories", handlers.authentication.AuthenticateJWT)
	categories.GET("", handlers.category.GetCategories)
	categories.POST("", handlers.category.CreateCategory)
	categories.DELETE("/:categoryId", handlers.category.DeleteCategory)
	categories.PATCH("/:categoryId", handlers.category.UpdateCategory)

	return e
}
