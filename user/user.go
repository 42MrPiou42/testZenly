package user

import (
	"time"
)

const (
	latM1 = 0x00000001 //Latitude + 1
	latL1 = 0x00000010 //Latitude - 1
	lonM1 = 0x00000100 //Longitude + 1
	lonL1 = 0x00001000 //Latitude - 1
)

type Position struct {
	Lon float64
	Lat float64
}

type User struct {
	Uuid uint64
	Pos Position
	Time time.Duration
	Colliders map[uint64]*User // Sort by Id User
	Key Position // truncaded position
	Flag uint8
}

func NewUser(id uint64, pos Position, tm time.Duration) (User) {
	return User{Uuid:id, Pos:pos, Time:tm, Colliders: make(map[uint64]*User)}
}

