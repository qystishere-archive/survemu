package game

import (
	"github.com/qystishere/survemu/game/data"
	"github.com/qystishere/survemu/game/object"
	"github.com/qystishere/survemu/game/utils"
)

func NewJoinedMsg(teamMode uint8, activePlayerID uint16, started, opts uint8, clients map[uint16]*Client) *SentMessage {
	m := NewSentMessage(data.MJoined, 4096)
	m.WriteUint8(teamMode)
	m.WriteUint16(activePlayerID)
	m.WriteUint8(started)
	m.WriteUint8(opts)

	m.WriteUint16(uint16(len(clients)))
	for _, client := range clients {
		m.WriteUint16(client.player.ID)
		m.WriteUint8(client.player.Team.ID)
		m.WriteString(client.player.Name)
	}
	return m
}

func NewMapMsg(width, height, shoreInset, grassInset uint16, seed, biome uint32,
	rivers []*object.River, places []*object.Place, objects map[uint16]*object.Object,
	groundPatches []*object.GroundPatch) *SentMessage {
	m := NewSentMessage(data.MMap, 15000)
	m.WriteUint16(width)
	m.WriteUint16(height)

	m.WriteUint16(shoreInset)
	m.WriteUint16(grassInset)

	m.WriteUint32(seed)
	m.WriteUint32(biome)

	m.WriteUint8(uint8(len(rivers)))
	for _, river := range rivers {
		m.WriteUint32(uint32(river.Width))
		m.WriteUint8(uint8(len(river.Points)))
		for _, point := range river.Points {
			m.WriteFloat(point.X, 0, 1024, 16)
			m.WriteFloat(point.Y, 0, 1024, 16)
		}
	}

	m.WriteUint8(uint8(len(places)))
	for _, place := range places {
		m.WriteString(place.Name)
		m.WriteFloat(place.Pos.X, 0, 1024, 16)
		m.WriteFloat(place.Pos.Y, 0, 1024, 16)
	}

	m.WriteUint16(uint16(len(objects)))
	for _, obj := range objects {
		m.WriteUint16(obj.ItemType)
		m.WriteFloat(obj.Pos.X, 0, 1024, 16)
		m.WriteFloat(obj.Pos.Y, 0, 1024, 16)
		m.WriteUnsignedBitInt32(uint32(obj.Ori), 2)
		m.WriteUnsignedBitInt32(uint32(obj.Scale), 6)
	}

	m.WriteUint8(uint8(len(groundPatches)))
	for _, groundPatch := range groundPatches {
		m.WriteFloat(groundPatch.Min.X, 0, 1024, 16)
		m.WriteFloat(groundPatch.Min.Y, 0, 1024, 16)

		m.WriteFloat(groundPatch.Max.X, 0, 1024, 16)
		m.WriteFloat(groundPatch.Max.Y, 0, 1024, 16)

		m.WriteUint32(groundPatch.Color)
		m.WriteUint32(uint32(groundPatch.Roughness))
		m.WriteUint32(uint32(groundPatch.OffsetDist))
	}
	return m
}

func NewPlayerInfoMsg(id uint16, teamID uint8, name string) *SentMessage {
	m := NewSentMessage(data.MPlayerInfo, 4096)
	m.WriteUint16(id)
	m.WriteUint8(teamID)
	m.WriteString(name)
	return m
}

func NewPingMsg(t uint8, pos *utils.Vector2D, useLoadout, isPing bool) *SentMessage {
	m := NewSentMessage(data.MEmote, 4096)
	m.WriteUint8(t)
	m.WriteFloat(pos.X, 0, 1024, 16)
	m.WriteFloat(pos.Y, 0, 1024, 16)
	m.WriteBool(useLoadout)
	m.WriteBool(isPing)
	m.WriteUnsignedBitInt32(0, 6)
	return m
}

func NewDisconnectMsg(reason string) *SentMessage {
	m := NewSentMessage(data.MDisconnect, 4096)
	m.WriteString(reason)
	return m
}

func NewUpdateMsg(t uint16, dObjs map[uint16]interface{}, fObjs map[uint16]interface{}, pObjs map[uint16]interface{}, player *object.Player,
	aliveCount byte, gas *utils.Gas, teams map[uint8]*object.Team, bullets []*object.Bullet, emotes []*object.Emote, ack byte) *SentMessage {
	m := NewSentMessage(data.MUpdate, 4096)
	m.WriteUint16(t)

	if (t & data.UTDeletedObjects) != 0 {
		m.WriteUint16(uint16(len(dObjs)))

		for _, obj := range dObjs {
			switch obj := obj.(type) {
			case *object.Object:
				m.WriteUint16(obj.ID)
			case *Client:
				m.WriteUint16(obj.player.ID)
			}
		}
	}

	if (t & data.UTFullObjects) != 0 {
		m.WriteUint16(uint16(len(fObjs)))

		for _, obj := range fObjs {
			switch obj := obj.(type) {
			case *object.Object:
				m.WriteUint8(obj.Type)
				m.WriteUint16(obj.ID)

				switch obj.Type {
				case data.OTLoot:
					// Part
					m.WriteFloat(obj.Pos.X, 0, 1024, 16)
					m.WriteFloat(obj.Pos.Y, 0, 1024, 16)

					// Full
					m.WriteUint8(uint8(obj.ItemType))
					m.WriteUint8(obj.Count)
					m.WriteUnsignedBitInt32(uint32(obj.Layer), 2)
					m.WriteUnsignedBitInt32(uint32(obj.Old), 6)
				}
			case *Client:
				m.WriteUint8(data.OTPlayer)
				m.WriteUint16(obj.player.ID)

				// Part
				m.WriteFloat(obj.player.Pos.X, 0, 1024, 16)
				m.WriteFloat(obj.player.Pos.Y, 0, 1024, 16)

				m.WriteFloat(obj.player.Ori.X, -1.0001, 1.0001, 8)
				m.WriteFloat(obj.player.Ori.Y, -1.0001, 1.0001, 8)

				// Full
				m.WriteSignedBitInt32(int32(obj.player.Layer), 2)

				m.WriteBool(obj.player.Dead)
				m.WriteBool(obj.player.Downed)

				m.WriteSignedBitInt32(int32(obj.player.AnimType), 3)
				m.WriteSignedBitInt32(int32(obj.player.AnimSeq), 3)

				m.WriteSignedBitInt32(int32(obj.player.ActionType), 2)

				m.WriteUint8(obj.player.Skin)
				m.WriteUint8(obj.player.Backpack)
				m.WriteUint8(obj.player.Helmet)
				m.WriteUint8(obj.player.Chest)
				if weapon := obj.player.WSlots[obj.player.CurrentWSlotID]; weapon != nil {
					m.WriteUint8(weapon.ID)
				} else {
					m.WriteUint8(0)
				}

				m.WriteBool(obj.player.WearingPan)
				m.WriteBool(obj.player.Frozen)
				m.WriteSignedBitInt32(obj.player.FrozenOri, 2)
			}
		}
	}

	m.WriteUint16(uint16(len(pObjs)))
	for _, obj := range pObjs {
		switch obj := obj.(type) {
		case *object.Object:
			m.WriteUint16(obj.ID)

			switch obj.ItemType {
			case data.OTLoot:
				// Part
				m.WriteFloat(obj.Pos.X, 0, 1024, 16)
				m.WriteFloat(obj.Pos.Y, 0, 1024, 16)
			}
		case *Client:
			m.WriteUint16(obj.player.ID)

			// Part
			m.WriteFloat(obj.player.Pos.X, 0, 1024, 16)
			m.WriteFloat(obj.player.Pos.Y, 0, 1024, 16)

			m.WriteFloat(obj.player.Ori.X, -1.0001, 1.0001, 8)
			m.WriteFloat(obj.player.Ori.Y, -1.0001, 1.0001, 8)
		}
	}

	if (t & data.UTActivePlayerId) != 0 {
		m.WriteUint16(player.ID)
	}

	m.WriteBool(player.ScopedIn)

	m.WriteFloat(player.Health, 0, 100, 8)

	boosted := player.Boost > 0
	m.WriteBool(boosted)
	if boosted {
		m.WriteFloat(player.Boost, 0, 100, 8)
	}

	inAction := player.Action != nil
	m.WriteBool(inAction)
	if inAction {
		m.WriteFloat(player.Action.Time, 0, 10, 12)
		m.WriteFloat(player.Action.Duration, 0, 10, 12)
		m.WriteUint8(player.Action.Type)
		m.WriteUint16(player.Action.TargetID)
	}

	m.WriteBool(true)
	if true {
		m.WriteUint8(player.Scope)
		// FIXME
		for _, amount := range player.Items {
			m.WriteUint16(amount)
		}
	}

	m.WriteBool(true)
	if true {
		m.WriteUnsignedBitInt32(uint32(player.CurrentWSlotID), 2)
		m.WriteUnsignedBitInt32(0, 6)
		for _, weapon := range player.WSlots {
			if weapon == nil {
				m.WriteUint8(0)
				m.WriteUint8(0)
			} else {
				m.WriteUint8(weapon.ID)
				m.WriteUint8(weapon.Ammo)
			}
		}
	}

	m.WriteBool(false)
	if false {
		m.WriteUint8(0)
	}

	m.WriteUnsignedBitInt32(0, 2)

	if (t & data.UTAliveCount) != 0 {
		m.WriteUint8(aliveCount)
	}

	m.WriteUint16(uint16(0))

	if (t & data.UTGas) != 0 {
		m.WriteUint8(gas.Mode)
		m.WriteUint32(uint32(gas.Duration))

		m.WriteUint16(uint16(gas.PosOld.X))
		m.WriteUint16(uint16(gas.PosOld.Y))
		m.WriteUint16(uint16(gas.PosNew.X))
		m.WriteUint16(uint16(gas.PosNew.Y))

		m.WriteUint16(uint16(gas.RadOld))
		m.WriteUint16(uint16(gas.RadNew))
	}

	if (t & data.UTTeamInfos) != 0 {
		m.WriteUint8(uint8(len(teams)))

		for _, team := range teams {
			m.WriteUint8(team.ID)
			m.WriteUint8(uint8(len(team.Players)))
			for _, player := range team.Players {
				m.WriteUint16(player.ID)
			}
		}
	}

	if (t & data.UTTeamData) != 0 {

	}

	if (t & data.UTBullets) != 0 {
		m.WriteUint8(uint8(len(bullets)))

		for _, bullet := range bullets {
			m.WriteUint16(bullet.PlayerID)
			m.WriteUint8(bullet.Type)

			m.WriteFloat(bullet.Pos.X, 0, 1024, 16)
			m.WriteFloat(bullet.Pos.Y, 0, 1024, 16)

			m.WriteFloat(bullet.Dir.X, -1.0001, 1.0001, 8)
			m.WriteFloat(bullet.Dir.Y, -1.0001, 1.0001, 8)

			m.WriteUnsignedBitInt32(uint32(bullet.Layer), 2)

			m.WriteFloat(bullet.VarianceT, 0, 1, 5)
			m.WriteUnsignedBitInt32(uint32(bullet.DistAdjID), 4)

			m.WriteBool(bullet.ClipDistance)
			if bullet.ClipDistance {
				m.WriteFloat(bullet.Distance, 0, 128, 8)
			}

			m.WriteBool(bullet.ShotFX > 0)
			if bullet.ShotFX > 0 {
				m.WriteUint8(bullet.ShotFX)
			}

			m.WriteUnsignedBitInt32(uint32(bullet.ShotOffHand), 7)
			m.WriteBool(bullet.LastShot)

			m.WriteUnsignedBitInt32(uint32(bullet.ReflectCount), 2)

			m.WriteBool(bullet.ReflectCount > 0)
			if bullet.ReflectCount > 0 {
				m.WriteUint16(bullet.ReflectObjID)
			}
		}
	}

	if (t & data.UTExplosions) != 0 {

	}

	if (t & data.UTEmotes) != 0 {
		m.WriteUint8(uint8(len(emotes)))

		for _, emote := range emotes {
			m.WriteUint8(emote.Type)

			if emote.IsPing {
				m.WriteUint8(1)
			} else {
				m.WriteUint8(0)
			}

			m.WriteUint16(emote.PlayerID)
			m.WriteFloat(emote.Pos.X, 0, 1024, 16)
			m.WriteFloat(emote.Pos.Y, 0, 1024, 16)
		}
	}

	if (t & data.UTPlanes) != 0 {

	}

	if (t & data.UTMapGlobalState) != 0 {

	}

	m.WriteUint8(ack)
	return m
}
