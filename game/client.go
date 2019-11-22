package game

import (
	"github.com/apex/log"
	"github.com/gorilla/websocket"

	"github.com/qystishere/survemu/game/data"
	"github.com/qystishere/survemu/game/object"
	"github.com/qystishere/survemu/game/utils"
)

type Client struct {
	input *utils.Input

	drop   *utils.Drop
	pickup *object.Object
	emote  *object.Emote

	player *object.Player
	world  *World

	ws *websocket.Conn
}

func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		ws: ws,
		input: &utils.Input{
			TouchMoveDir: &utils.Vector2D{},
			ToMouseDir:   &utils.Vector2D{},
		},
	}
}

func (c *Client) Handle(m *ReceivedMessage) {
	defer func() {
		if r := recover(); r != nil {
			c.Close()
		}
	}()

	switch m.ID {
	case data.MJoin:
		protocol, _ := m.ReadUint32()
		privData, _ := m.ReadString(16)
		name, _ := m.ReadString(16)
		_, _ = m.ReadBits(uint(15-len(name)) * 8)

		emotes := [data.ESCount]byte{}
		for key := range emotes {
			emotes[key], _ = m.ReadUint8()
		}

		c.player = object.NewPlayer(name, emotes)

		WorldBase.addObjCh <- c

		c.world = WorldBase
		log.Debugf("[GAME] [CLIENT] Handle MJoin (%v, %v, %v)", protocol, privData, name)
	case data.MInput:
		_, _ = m.ReadUint8() // seq
		c.input.MoveLeft, _ = m.ReadBool()
		c.input.MoveRight, _ = m.ReadBool()
		c.input.MoveUp, _ = m.ReadBool()
		c.input.MoveDown, _ = m.ReadBool()

		if shootStart, _ := m.ReadBool(); shootStart {
			c.input.ShootStart = true
		}
		c.input.ShootHold, _ = m.ReadBool()

		c.input.Portrait, _ = m.ReadBool()
		c.input.TouchMoveActive, _ = m.ReadBool()

		if c.input.TouchMoveActive {
			c.input.TouchMoveDir.X, _ = m.ReadFloat(-1.0001, 1.0001, 8)
			c.input.TouchMoveDir.Y, _ = m.ReadFloat(-1.0001, 1.0001, 8)
			c.input.TouchMoveLen, _ = m.ReadUint8()
		}

		c.input.ToMouseDir.X, _ = m.ReadFloat(-1.0001, 1.0001, 10)
		c.input.ToMouseDir.Y, _ = m.ReadFloat(-1.0001, 1.0001, 10)
		c.input.ToMouseLen, _ = m.ReadFloat(0, 64, 8)

		inputsLen, _ := m.ReadInt32Bits(4)

		if len(c.input.Ips) != int(inputsLen) {
			c.input.Ips = make([]byte, inputsLen)
		}

		for ; inputsLen > 1; inputsLen-- {
			c.input.Ips[inputsLen], _ = m.ReadUint8()
		}

		if useItem, _ := m.ReadUint8(); useItem != 0 {
			c.input.UseItem = useItem
			log.Debugf("useItem: %v", useItem)
		}
	case data.MEmote:
		t, _ := m.ReadUint8()
		pos := &utils.Vector2D{}
		pos.X, _ = m.ReadFloat(0, 1024, 16)
		pos.Y, _ = m.ReadFloat(0, 1024, 16)
		useLoadout, _ := m.ReadBool()
		isPing, _ := m.ReadBool()
		_, _ = m.ReadBits(6)

		if useLoadout {
			t = c.player.Emotes[t]
		}

		c.emote = &object.Emote{
			Type:       t,
			Pos:        pos,
			UseLoadout: useLoadout,
			IsPing:     isPing,

			PlayerID: c.player.ID,
		}
	case data.MDropItem:
		id, _ := m.ReadUint8()
		wSlotID, _ := m.ReadUint8()

		c.drop = &utils.Drop{
			ID:      id,
			WSlotID: wSlotID,
		}
		log.Debugf("DROP: %v", c.drop)
	default:
		log.Debugf("[GAME] [CLIENT] Handle: %v (%x)", m.ID, m.Data())
	}
}

func (c *Client) Send(m *SentMessage) {
	err := c.ws.WriteMessage(websocket.BinaryMessage, append([]byte{m.ID}, m.Data()...))
	if err != nil {
		c.Close()
	}
}

func (c *Client) Close() {
	_ = c.ws.Close()

	if c.world != nil {
		c.world.rmObjCh <- c
	}
}
