package data

// MessageType
const (
	MNone        byte = 0
	MJoin             = 1
	MDisconnect       = 2
	MInput            = 3
	MEdit             = 4
	MJoined           = 5
	MPlayerInfo       = 6
	MUpdate           = 7
	MKill             = 8
	MGameOver         = 9
	MPickup           = 10
	MMap              = 11
	MSpectate         = 12
	MDropItem         = 13
	MEmote            = 14
	MPlayerStats      = 15
	MAdStatus         = 16
	MLoadout          = 17
)

// UpdateType
const (
	UTDeletedObjects uint16 = 1
	UTFullObjects           = 2
	UTActivePlayerId        = 4
	UTAliveCount            = 8
	UTGas                   = 16
	UTTeamData              = 32
	UTTeamInfos             = 64
	UTBullets               = 128
	UTExplosions            = 256
	UTEmotes                = 512
	UTPlanes                = 1024
	UTMapGlobalState        = 2048
)

// EmoteSlot
const (
	None   byte = 0
	Top         = 1
	Right       = 2
	Bottom      = 3
	Left        = 4
	Win         = 5
	Death       = 6
)

const ESCount = 7

// ObjectType
const (
	OTЕInvalid    byte = 0
	OTPlayer           = 1
	OTObstacle         = 2
	OTLoot             = 3
	OTLootSpawner      = 4
	OTDeadBody         = 5
	OTBuilding         = 6
	OTStructure        = 7
	OTDecal            = 8
	OTProjectile       = 9
	OTSmoke            = 10
	OTAirdrop          = 11
)

var BagSizes = [][]uint16{
	{120, 240, 330, 420},
	{90, 180, 240, 300},
	{90, 180, 240, 300},
	{15, 30, 60, 90},
	{49, 98, 147, 196},
	{10, 20, 30, 40},
	{2, 4, 6, 8},
	{90, 180, 240, 300},
	{3, 6, 9, 12},
	{3, 6, 9, 12},
	{2, 4, 6, 8},
	{10, 20, 30, 40},
	{5, 10, 15, 30},
	{1, 2, 3, 4},
	{2, 5, 10, 15},
	{1, 2, 3, 4},
	{1, 1, 1, 1},
	{1, 1, 1, 1},
	{1, 1, 1, 1},
	{1, 1, 1, 1},
	{1, 1, 1, 1},
}

const BSCount = 21

// WeaponSlot
const (
	WSPrimary   byte = 0
	WSSecondary      = 1
	WSThrowable      = 3
	WSMelee          = 2
	WSCount          = 4
)

// Anim
const (
	ANNone          byte = 0
	ANMelee              = 1
	ANCook               = 2
	ANThrow              = 3
	ANCrawlForward       = 4
	ANCrawlBackward      = 5
	ANRevive             = 6
)

// Action
const (
	ACNone    byte = 0
	ACReload       = 1
	ACUseItem      = 2
	ACRevive       = 3
)

// DamageType
const (
	DTPlayer   byte = 0
	DTBleeding      = 1
	DTGas           = 2
	DTAirdrop       = 3
)
