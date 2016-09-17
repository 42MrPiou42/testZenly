package user

import (
	"time"
	"fmt"
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

func (us *User)LeaveUser(idOther uint64) {
	if us.Uuid < idOther {
		fmt.Println(us.Uuid, " a quitté ", idOther)
	}
	delete(us.Colliders, idOther)
}

func (us *User)JoinUser(usOther *User) {
	if (us.Uuid < usOther.Uuid) {
		fmt.Println(us.Uuid, " est avec ", usOther.Uuid)
	}
	us.Colliders[usOther.Uuid] = usOther
}

func (us *User)GetKeys(precision uint) (pos Position) {
	pos.Lon = float64(int(us.Pos.Lon * precision)) / precision
	pos.Lat = float64(int(us.Pos.Lat * precision)) / precision
	return
}