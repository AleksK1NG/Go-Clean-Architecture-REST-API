package news

import "github.com/labstack/echo"

// News Http Delivery
type Delivery interface {
	Create() echo.HandlerFunc
}
