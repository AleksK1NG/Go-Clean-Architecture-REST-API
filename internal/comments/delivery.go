package comments

import "github.com/labstack/echo"

// Comments handlers
type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	GetAllByNewsID() echo.HandlerFunc
}
