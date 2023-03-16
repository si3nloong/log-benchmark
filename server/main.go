package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/article", func(c echo.Context) error {
		return c.JSON(http.StatusOK, json.RawMessage(`{"id": 0, "name": "My Great Article"}`))
	})
	e.GET("/articles", func(c echo.Context) error {
		var msg string
		msg += `[`
		for i := 0; i < 100; i++ {
			if i > 0 {
				msg += `,`
			}
			msg += fmt.Sprintf(`{"id": %d, "name": "My Great Article"}`, i+1)
		}
		msg += `]`
		return c.JSON(http.StatusOK, json.RawMessage(msg))
	})
	e.Logger.Fatal(e.Start(":1000"))
}
