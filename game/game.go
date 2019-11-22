package game

import (
	"github.com/apex/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/qystishere/survemu/game/data"
	"github.com/qystishere/survemu/game/object"
	"github.com/qystishere/survemu/game/utils"
)

var WorldBase = NewWorld()

func init() {
	row := 0
	j := 0
	for i := uint16(0); i < 96; i++ {
		WorldBase.Objects[i] = object.NewObject(
			i, data.OTLoot, 50+i, &utils.Vector2D{X: 65 + float32(j*5), Y: 130 - float32(row*5)}, 0, 100, 1)
		if j >= 15 {
			row++
			j = 0
		} else {
			j++
		}
	}
}

func Start(addr string) {
	go WorldBase.Start()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())

	e.GET("/play", Play)

	log.Fatal(e.Start(addr).Error())
}
