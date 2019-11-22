package object

type Team struct {
	ID uint8

	Players map[uint16]*Player
}

func NewTeam(id byte) *Team {
	return &Team{
		ID:      id,
		Players: map[uint16]*Player{},
	}
}
