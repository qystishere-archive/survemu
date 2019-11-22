package web

import (
	"github.com/apex/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Start(addr string) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

	e.Static("/", "./resources/web/share")

	InitApiHandlers(e)

	log.Fatal(e.Start(addr).Error())
}