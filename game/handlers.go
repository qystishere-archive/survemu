package game

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	"github.com/qystishere/survemu/pkg/bitbuf"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Play(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		panic(err)
		return err
	}
	defer func() {
		_ = ws.Close()
	}()
	log.Debugf("[GAME] [HANDLERS] Connected: %v", ws.RemoteAddr())

	client := NewClient(ws)
	m := &ReceivedMessage{
		Message: &Message{},
	}
	for {
		t, rawMsg, err := ws.ReadMessage()
		if err != nil || t != websocket.BinaryMessage {
			break
		}

		if len(rawMsg) < 2 {
			break
		}

		m.ID = rawMsg[0]
		m.Reader = bitbuf.NewReader(rawMsg[1:])

		client.Handle(m)
	}
	log.Debugf("[GAME] [HANDLERS] Disconnected: %v", ws.RemoteAddr())
	return nil
}
