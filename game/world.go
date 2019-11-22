package game

import (
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/apex/log"

	"github.com/qystishere/survemu/game/data"
	"github.com/qystishere/survemu/game/object"
	"github.com/qystishere/survemu/game/utils"
	"github.com/qystishere/survemu/pkg/rnd"
)

type World struct {
	TeamMode byte

	Width  uint16
	Height uint16

	ShoreInset uint16
	GrassInset uint16

	Seed  uint32
	Biome uint32

	Gas *utils.Gas

	Rivers        []*object.River
	Places        []*object.Place
	Objects       map[uint16]*object.Object
	GroundPatches []*object.GroundPatch

	Teams map[uint8]*object.Team

	Clients map[uint16]*Client

	addObjCh chan interface{}
	rmObjCh  chan interface{}

	teamSequence  uint8
	objIDSequence uint16
	fps           int

	_running int32
}

func NewWorld() *World {
	return &World{
		TeamMode: 1,

		Width:  200,
		Height: 200,

		ShoreInset: 48,
		GrassInset: 18,

		Seed:  43207811,
		Biome: 0,

		Gas: &utils.Gas{
			Mode:     0,
			Duration: 0,
			PosOld:   &utils.Vector2D{X: 100, Y: 100},
			PosNew:   &utils.Vector2D{X: 110, Y: 110},
			RadOld:   0,
			RadNew:   0,
		},

		Rivers:        []*object.River{},
		Places:        []*object.Place{},
		Objects:       map[uint16]*object.Object{},
		GroundPatches: []*object.GroundPatch{},

		Teams: map[uint8]*object.Team{},

		Clients: map[uint16]*Client{},

		addObjCh: make(chan interface{}, 5),
		rmObjCh:  make(chan interface{}, 5),

		teamSequence:  0,
		objIDSequence: 3000,
		fps:           30,
	}
}

func (w *World) Start() {
	w._running = 1

	dur := time.Millisecond * time.Duration(1000/w.fps)
	for {
		time.Sleep(dur)

		if !(w._running == 1) {
			return
		}

		w.Tick()
	}
}

func (w *World) Stop() {
	atomic.StoreInt32(&w._running, 0)
}

func (w *World) Running() bool {
	return atomic.LoadInt32(&w._running) == 1
}

func (w *World) NextObjID() uint16 {
	w.objIDSequence++
	if w.objIDSequence >= 32000 {
		w.objIDSequence = 2500
	}
	return w.objIDSequence
}

func (w *World) Tick() {
	gDeletedClients := map[uint16]*Client{}
	gDeletedObjects := map[uint16]interface{}{}

	for len(w.addObjCh) > 0 {
		obj := <-w.addObjCh

		switch obj := obj.(type) {
		case *Client:
			obj.player.ID = w.NextObjID()

			w.Clients[obj.player.ID] = obj

			if len(w.Teams) == 0 {
				w.teamSequence++
				team := object.NewTeam(w.teamSequence)
				w.Teams[w.teamSequence] = team
			}

			for _, team := range w.Teams {
				team.Players[obj.player.ID] = obj.player
				obj.player.Team = team
				break
			}

			obj.Send(NewJoinedMsg(
				WorldBase.TeamMode,
				obj.player.ID,
				1,
				2,
				w.Clients,
			))

			obj.Send(NewMapMsg(
				w.Width,
				w.Height,
				w.ShoreInset,
				w.GrassInset,
				w.Seed,
				w.Biome,
				w.Rivers,
				w.Places,
				w.Objects,
				w.GroundPatches,
			))

			fullObjects := map[uint16]interface{}{}
			for _, client := range w.Clients {
				client.Send(NewPlayerInfoMsg(obj.player.ID, obj.player.Team.ID, obj.player.Name))

				if client == obj || obj.player.Pos.Distance(client.player.Pos) < 100 {
					fullObjects[client.player.ID] = client
				}
			}

			obj.Send(NewUpdateMsg(
				data.UTActivePlayerId+data.UTTeamInfos+data.UTAliveCount+data.UTFullObjects,
				map[uint16]interface{}{},
				fullObjects,
				map[uint16]interface{}{},
				obj.player,
				byte(len(w.Clients)),
				w.Gas,
				w.Teams,
				[]*object.Bullet{},
				[]*object.Emote{},
				10,
			))
		case *object.Object:
			obj.ID = w.NextObjID()

			w.Objects[obj.ID] = obj
		default:
			log.Errorf("[GAME] [WORLD] Add unknown object: %v", obj)
		}
	}

	for len(w.rmObjCh) > 0 {
		obj := <-w.rmObjCh

		switch obj := obj.(type) {
		case *Client:
			delete(w.Clients, obj.player.ID)
			gDeletedClients[obj.player.ID] = obj
		case *object.Object:
			delete(w.Objects, obj.ID)
			gDeletedObjects[obj.ID] = obj
		default:
			log.Errorf("[GAME] [WORLD] Remove unknown object: %v", obj)
		}
	}

	// Input
	for _, client := range w.Clients {
		if client.player.Ori.X != client.input.ToMouseDir.X || client.player.Ori.Y != client.input.ToMouseDir.Y {
			client.player.Ori.X = client.input.ToMouseDir.X
			client.player.Ori.Y = client.input.ToMouseDir.Y

			client.player.RequirePartUpdate = true
		}

		multiplier := float32(1)
		if (client.input.MoveLeft && client.input.MoveUp) || (client.input.MoveLeft && client.input.MoveDown) ||
			(client.input.MoveRight && client.input.MoveUp) || (client.input.MoveRight && client.input.MoveDown) {
			multiplier = 0.75
		}

		if client.input.MoveLeft {
			client.player.Pos.X -= client.player.Speed * multiplier
		}
		if client.input.MoveRight {
			client.player.Pos.X += client.player.Speed * multiplier
		}
		if client.input.MoveUp {
			client.player.Pos.Y += client.player.Speed * multiplier
		}
		if client.input.MoveDown {
			client.player.Pos.Y -= client.player.Speed * multiplier
		}

		if client.input.MoveLeft || client.input.MoveRight || client.input.MoveUp || client.input.MoveDown {
			client.player.RequirePartUpdate = true
		}

		if client.player.Action != nil {
			since := time.Since(client.player.Action.CreatedAt)

			client.player.Action.Time = float32(since.Seconds())

			if since >= time.Second*time.Duration(client.player.Action.Duration) {
				client.player.Action = nil
				client.player.ActionType = 0
			}
		}

		if client.input.UseItem != 0 {
			switch client.input.UseItem {
			case 5:
				if client.player.Action != nil {
					break
				}
				client.player.ActionType = 1
				client.player.Action = &utils.Action{
					Time:      0,
					Duration:  3,
					Type:      1,
					CreatedAt: time.Now(),
					TargetID:  0,
				}
				client.player.RequireFullUpdate = true
			case 7: // pickup
				var nearestLoot *object.Object
				for _, obj := range w.Objects {
					if obj.Type == data.OTLoot {
						pos := client.player.Pos.Distance(obj.Pos)

						if pos <= 1.2 {
							if nearestLoot != nil && client.player.Pos.Distance(nearestLoot.Pos) > pos {
								nearestLoot = obj
							} else {
								nearestLoot = obj
							}
						}
					}
				}
				if nearestLoot != nil {
					client.pickup = nearestLoot

					client.player.RequireFullUpdate = true
				}
			case 11, 12, 13, 14:
				client.player.CurrentWSlotID = client.input.UseItem - 11
				client.player.RequireFullUpdate = true
			default:
				item, ok := data.GetItemNameByID(uint16(client.input.UseItem))
				if ok {
					if strings.Index(item, "scope") != -1 {
						client.player.Scope = uint8(client.input.UseItem)
					}
				}
			}
		}

		if client.input.ShootStart || client.input.ShootHold {
			weapon := client.player.WSlots[client.player.CurrentWSlotID]
			if weapon != nil && time.Since(weapon.LastUsedAt) >= weapon.Duration {
				switch client.player.CurrentWSlotID {
				case 0, 1:
					client.player.AnimType = data.ANNone

					source := &utils.Vector2D{
						X: client.player.Pos.X + (client.player.Ori.X * 2),
						Y: client.player.Pos.Y + (client.player.Ori.Y * 2),
					}

					switch weapon.ID {
					case 70:
						source1 := &utils.Vector2D{
							X: client.player.Pos.X + 2,
							Y: client.player.Pos.Y + 2,
						}
						source2 := &utils.Vector2D{
							X: client.player.Pos.X - 2,
							Y: client.player.Pos.Y - 2,
						}
						ori := &utils.Vector2D{
							X: client.player.Ori.X,
							Y: client.player.Ori.Y,
						}
						ori1 := &utils.Vector2D{
							X: client.player.Ori.X + 2,
							Y: client.player.Ori.Y + 2,
						}
						ori2 := &utils.Vector2D{
							X: client.player.Ori.X - 2,
							Y: client.player.Ori.Y - 2,
						}
						client.player.Bullets = []*object.Bullet{
							object.NewBullet(
								client.player.ID, 163, source, ori, weapon.ID,
							),
							object.NewBullet(
								client.player.ID, 163, source1, ori1, weapon.ID,
							),
							object.NewBullet(
								client.player.ID, 163, source2, ori2, weapon.ID,
							),
						}
					default:
						client.player.Bullets = []*object.Bullet{
							object.NewBullet(client.player.ID, uint8(rnd.NextLimit(155, 190)), source, client.player.Ori, weapon.ID,
							),
						}
					}
				case 2:
					if client.input.ShootStart {
						client.player.AnimType = data.ANMelee
					} else {
						client.player.AnimType = data.ANNone
					}
				case 3:
					if client.input.ShootStart {
						client.player.AnimType = data.ANCook
					} else {
						client.player.AnimType = data.ANNone
					}
				}

				client.player.AnimSeq += 1
				if client.player.AnimSeq > 7 {
					client.player.AnimSeq = 0
				}

				weapon.LastUsedAt = time.Now()
			}
			client.player.RequireFullUpdate = true
		} else {
			client.player.AnimType = data.ANNone
		}
	}

	// Update
	for id, client := range w.Clients {
		updateType := uint16(0)

		deletedObjects := map[uint16]interface{}{}
		fullObjects := map[uint16]interface{}{}
		partObjects := map[uint16]interface{}{}
		bullets := []*object.Bullet{}
		emotes := []*object.Emote{}

		deletePickup := func() {
			delete(w.Objects, client.pickup.ID)
			deletedObjects[client.pickup.ID] = client.pickup
		}

		drop := func(id byte, count byte) {
			if id == 0 || id == 49 || id == 131 {
				return
			}

			newObject := object.NewObject(
				w.NextObjID(),
				data.OTLoot,
				uint16(id),
				&utils.Vector2D{
					X: client.player.Pos.X,
					Y: client.player.Pos.Y,
				},
				0,
				count,
				0,
			)

			w.Objects[newObject.ID] = newObject
			fullObjects[newObject.ID] = newObject
		}

		if client.drop != nil {
			switch {
			case client.drop.ID >= 50 && client.drop.ID <= 66:
				weapon := client.player.WSlots[2]

				if weapon != nil && weapon.ID == client.drop.ID {
					client.player.WSlots[2] = &utils.Weapon{
						ID:   49,
						Ammo: 0,

						Duration: time.Millisecond * 250,
					}

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 67 && client.drop.ID <= 112:
				weapon := client.player.WSlots[client.drop.WSlotID]

				if weapon != nil && weapon.ID == client.drop.ID {
					if client.player.CurrentWSlotID == client.drop.WSlotID {
						client.player.CurrentWSlotID = 2
					}

					client.player.WSlots[client.drop.WSlotID] = nil

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 113 && client.drop.ID <= 118:
				weapon := client.player.WSlots[3]

				if weapon != nil && weapon.ID == client.drop.ID {
					client.player.CurrentWSlotID = 2

					client.player.WSlots[3] = nil

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 119 && client.drop.ID <= 126:
				slotItem := int(client.drop.ID - 119)
				amount := math.Ceil(float64(client.player.Items[client.drop.ID-119]) / 2)
				if amount > 0 {
					client.player.Items[slotItem] -= uint16(amount)
					drop(client.drop.ID, byte(amount))
				}
			case client.drop.ID >= 127 && client.drop.ID <= 130:
				slotItem := 12 + int(client.drop.ID-127)
				amount := math.Ceil(float64(12+client.player.Items[client.drop.ID-127]) / 2)
				if amount > 0 {
					client.player.Items[slotItem] -= uint16(amount)
					drop(client.drop.ID, byte(amount))
				}
			case client.drop.ID >= 132 && client.drop.ID <= 134:
				if client.player.Backpack != 131 {
					client.player.Backpack = 131

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 135 && client.drop.ID <= 137:
				if client.player.Helmet != 0 {
					client.player.Helmet = 0

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 138 && client.drop.ID <= 140:
				if client.player.Chest != 0 {
					client.player.Chest = 0

					drop(client.drop.ID, 1)
				}
			case client.drop.ID >= 141 && client.drop.ID <= 145:
				slotItem := 16 + int(client.drop.ID-141)
				amount := client.player.Items[client.drop.ID-141]
				if amount > 0 {
					if client.player.Scope == client.drop.ID {
						client.player.Scope = 141
					}
					client.player.Items[slotItem] = amount
					drop(client.drop.ID, byte(amount))
				}
			}

			client.player.RequireFullUpdate = true
		}

		if client.pickup != nil {
			slot := byte(0)
			slotItem := -1

			switch {
			case client.pickup.ItemType >= 50 && client.pickup.ItemType <= 66:
				slot = 2
			case client.pickup.ItemType >= 67 && client.pickup.ItemType <= 112:
				for slotID, weapon := range client.player.WSlots {
					if weapon == nil {
						slot = byte(slotID)
						break
					}
				}
			case client.pickup.ItemType == 115:
				slotItem = -2
			case client.pickup.ItemType >= 113 && client.pickup.ItemType != 115 && client.pickup.ItemType <= 118:
				slot = 3
			case client.pickup.ItemType >= 119 && client.pickup.ItemType <= 126:
				slotItem = int(client.pickup.ItemType - 119)
			case client.pickup.ItemType >= 127 && client.pickup.ItemType <= 130:
				slotItem = 12 + int(client.pickup.ItemType-127)
			case client.pickup.ItemType == 131: // first backpack
				slotItem = -2
			case client.pickup.ItemType >= 132 && client.pickup.ItemType <= 134:
				slotItem = -2
				if client.player.Backpack < byte(client.pickup.ItemType) {
					drop(client.player.Backpack, 1)
				} else {
					break
				}
				deletePickup()
				client.player.Backpack = byte(client.pickup.ItemType)
			case client.pickup.ItemType >= 135 && client.pickup.ItemType <= 137:
				slotItem = -2
				if client.player.Helmet < byte(client.pickup.ItemType) {
					drop(client.player.Helmet, 1)
				} else {
					break
				}
				deletePickup()
				client.player.Helmet = byte(client.pickup.ItemType)
			case client.pickup.ItemType >= 138 && client.pickup.ItemType <= 140:
				slotItem = -2
				if client.player.Chest < byte(client.pickup.ItemType) {
					drop(client.player.Chest, 1)
				} else {
					break
				}
				deletePickup()
				client.player.Chest = byte(client.pickup.ItemType)
			case client.pickup.ItemType >= 141 && client.pickup.ItemType <= 145:
				slotItem = 16 + int(client.pickup.ItemType-141)
				client.player.Scope = byte(client.pickup.ItemType)
			}

			if slotItem == -1 {
				weapon := client.player.WSlots[slot]
				if weapon != nil {
					drop(weapon.ID, 1)
				}

				deletePickup()

				newWeapon := &utils.Weapon{
					ID:   uint8(client.pickup.ItemType),
					Ammo: 100,
				}

				if newWeapon.ID >= 50 && newWeapon.ID <= 66 {
					newWeapon.Duration = time.Millisecond * 250
				}

				client.player.WSlots[slot] = newWeapon

				client.player.CurrentWSlotID = slot
			} else if slotItem != -2 {
				amount := data.BagSizes[slotItem][client.player.Backpack-131] - client.player.Items[slotItem]
				if uint16(client.pickup.Count) <= amount {
					if uint16(client.pickup.Count) < amount {
						amount = uint16(client.pickup.Count)
					}
					deletePickup()
				} else {
					client.pickup.Count -= uint8(amount)
				}
				client.player.Items[slotItem] += amount
			}
		}

		// Delete outdated

		for _, obj := range client.player.LoadedObjects {
			pos := client.player.Pos.Distance(obj.Pos)

			if _, ok := gDeletedObjects[obj.ID]; ok || pos >= 100 {
				deletedObjects[obj.ID] = obj

				delete(client.player.LoadedObjects, obj.ID)
			}
		}

		for _, player := range client.player.LoadedPlayers {
			client2, ok := w.Clients[player.ID]
			if !ok {
				if client3, ok2 := gDeletedClients[player.ID]; ok2 {
					deletedObjects[player.ID] = client3

					delete(client.player.LoadedPlayers, player.ID)
				}
			} else {
				if player.Pos.Distance(client.player.Pos) >= 100 {
					deletedObjects[player.ID] = client2

					delete(client.player.LoadedPlayers, player.ID)
				}
			}
		}

		// Interactions

		for _, obj := range w.Objects {
			pos := client.player.Pos.Distance(obj.Pos)

			if pos < 100 {
				_, loaded := client.player.LoadedObjects[obj.ID]

				if !loaded {
					fullObjects[obj.ID] = obj

					client.player.LoadedObjects[obj.ID] = obj
				}
			}
		}

		for someID, someClient := range w.Clients {
			if id == someID { // Local
				if someClient.emote != nil {
					emotes = append(emotes, someClient.emote)
				}

				if someClient.player.Bullets != nil {
					for _, bullet := range someClient.player.Bullets {
						bullets = append(bullets, bullet)
					}
				}

				if someClient.player.RequireFullUpdate {
					fullObjects[someClient.player.ID] = someClient
				} else if someClient.player.RequirePartUpdate {
					partObjects[someClient.player.ID] = someClient
				}
			} else {
				pos := client.player.Pos.Distance(someClient.player.Pos)

				if pos < 100 {
					if someClient.emote != nil {
						emotes = append(emotes, someClient.emote)
					}

					if someClient.player.Bullets != nil {
						for _, bullet := range someClient.player.Bullets {
							bullets = append(bullets, bullet)
						}
					}

					_, loaded := client.player.LoadedPlayers[someClient.player.ID]

					if !loaded {
						fullObjects[someClient.player.ID] = someClient

						client.player.LoadedPlayers[someClient.player.ID] = someClient.player
					} else {
						if someClient.player.RequireFullUpdate {
							fullObjects[someClient.player.ID] = someClient
						} else if someClient.player.RequirePartUpdate {
							partObjects[someClient.player.ID] = someClient
						}
					}
				}
			}
		}

		if len(deletedObjects) > 0 {
			updateType |= data.UTDeletedObjects
		}

		if len(fullObjects) > 0 {
			updateType |= data.UTFullObjects
		}

		if len(bullets) > 0 {
			updateType |= data.UTBullets
		}

		if len(emotes) > 0 {
			updateType |= data.UTEmotes
		}

		client.Send(NewUpdateMsg(
			updateType,
			deletedObjects,
			fullObjects,
			partObjects,
			client.player,
			55,
			w.Gas,
			w.Teams,
			bullets,
			emotes,
			10,
		))
	}

	for _, client := range w.Clients {
		client.input.ShootStart = false
		client.input.UseItem = 0

		client.drop = nil
		client.pickup = nil
		client.emote = nil

		client.player.Bullets = nil

		client.player.RequirePartUpdate = false
		client.player.RequireFullUpdate = false
	}
}
