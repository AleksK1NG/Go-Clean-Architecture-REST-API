package news

import "github.com/labstack/echo"

// News Http Delivery
type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
}
