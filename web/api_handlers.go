package web

import (
	"runtime"

	"github.com/labstack/echo"
)

func InitApiHandlers(e *echo.Echo) {
	e.GET("/api/site_info", SiteInfo)

	e.POST("/api/find_game", FindGame)
}

func SiteInfo(c echo.Context) error {
	return c.JSON(200, &SiteInfoResult{
		Pops: &Pops{
			Eu: "? players",
			Na: "offline",
			Sa: "offline",
			As: "offline",
		},
		PromptConsent: false,
	})
}

func FindGame(c echo.Context) error {
	if runtime.GOOS != "windows" {
		return c.JSON(200, &FindGameResult{
			Res: []*Game{
				{
					Zone:   "rus",
					GameID: "gameID",
					Hosts:  []string{":8081"},
					Addrs:  []string{"92.38.139.88:8081"},
					Data:   "data",
				},
			},
		})
	}

	return c.JSON(200, &FindGameResult{
		Res: []*Game{
			{
				Zone:   "rus",
				GameID: "gameID",
				Hosts:  []string{"localhost:8081"},
				Addrs:  []string{"127.0.0.1:8081"},
				Data:   "data",
			},
		},
	})
}
