package middleware

import (
	"errors"
	"net/http"

	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"

	"github.com/labstack/echo/v4"
)

func ErrorHandlingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			if err != nil {
				return handleError(ctx, err)
			}
			return nil
		}
	}
}

func handleError(ctx echo.Context, err error) error {
	// Validation errors -> 400 Bad Request
	if errors.Is(err, errs.ErrValueIsInvalid) ||
		errors.Is(err, errs.ErrCommandIsInvalid) ||
		errors.Is(err, errs.ErrQueryIsInvalid) {
		return ctx.JSON(http.StatusBadRequest, servers.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Business logic conflicts -> 409 Conflict
	if errors.Is(err, errs.ErrVersionIsInvalid) {
		return ctx.JSON(http.StatusConflict, servers.Error{
			Code:    http.StatusConflict,
			Message: err.Error(),
		})
	}

	// Internal server errors -> 500
	return ctx.JSON(http.StatusInternalServerError, servers.Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	})
}
