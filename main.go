package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var num int

func main() {
	app := pocketbase.New()
	fmt.Println("Run main")
	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		e.Router.GET("/events", func(c echo.Context) error {
			fmt.Println("Running events")
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Type")
			c.Response().Header().Set("Content-Type", "text/event-stream")
			c.Response().Header().Set("Cache-Control", "no-cache")
			c.Response().Header().Set("Connection", "keep-alive")

			for num <= 15 {
				_, err := fmt.Fprintf(c.Response().Writer, "data: %s\n\n", fmt.Sprintf("Event %d", num))
				if err != nil {
					log.Println("Error writing SSE event", err)
					break
				}
				time.Sleep(2 * time.Second)
				c.Response().Flush()
				num++
			}
			return nil
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
