package auth

import "github.com/labstack/echo"

// Auth Delivery interface
type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
}
