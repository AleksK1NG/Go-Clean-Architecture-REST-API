package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server")
	//http.HandleFunc("/speed", func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(200)
	//	w.Write([]byte(r.RemoteAddr))
	//})
	//
	//log.Fatal(http.ListenAndServe(":8080", nil))

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// Routes
	e.GET("/speed", hello)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
