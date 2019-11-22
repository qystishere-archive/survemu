package main

import (
	"github.com/apex/log"

	"github.com/qystishere/survemu/game"
	"github.com/qystishere/survemu/web"
)

func main() {
	log.SetLevel(log.DebugLevel)

	go game.Start(":8081")

	web.Start(":8080")
}
