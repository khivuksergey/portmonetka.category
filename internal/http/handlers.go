package http

import (
	"github.com/khivuksergey/portmonetka.category/internal/core/port/service"
	"github.com/khivuksergey/portmonetka.category/internal/handler"
	"github.com/khivuksergey/portmonetka.common/middleware/authentication"
	"github.com/khivuksergey/portmonetka.common/middleware/error"
	"github.com/khivuksergey/webserver/logger"
	"github.com/spf13/viper"
)

type Handlers struct {
	error          *error.ErrorHandlingMiddleware
	authentication *authentication.AuthenticationMiddleware
	category       *handler.CategoryHandler
}

func newHandlers(services *service.Manager, logger logger.Logger) Handlers {
	return Handlers{
		error:          error.NewErrorHandlingMiddleware(),
		authentication: authentication.NewAuthenticationMiddleware(viper.GetString("JWT_SECRET"), logger),
		category:       handler.NewCategoryHandler(services, logger),
	}
}
