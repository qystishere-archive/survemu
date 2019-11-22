package utils

type Input struct {
	MoveLeft  bool
	MoveRight bool
	MoveUp    bool
	MoveDown  bool

	ShootStart bool
	ShootHold  bool
	Portrait   bool

	TouchMoveActive bool
	TouchMoveDir    *Vector2D
	TouchMoveLen    byte

	ToMouseDir *Vector2D
	ToMouseLen float32

	Ips []byte

	UseItem byte
}
