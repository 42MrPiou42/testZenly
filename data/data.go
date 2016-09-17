package data

import (
	"proj/user"
	"time"
)


type Data struct {
	Users map[uint64]user.User
	BigData map[float64]map[float64]*user.User
}

func GetData() (dt Data) {
	dt.Users = make(map[uint64]user.User)
	return
}

func (Dt Data) GetUser(id uint64) (*user.User) {
	if usr, ok := Dt.Users[id]; ok == true {
		return &usr
	}
	return nil
}

func (Dt Data) AddUser(id uint64, pos user.Position, tm time.Duration) (*user.User) {
	Dt.Users[id] = user.NewUser(id, pos, tm)
	return &(Dt.Users[id])
}

func (Dt Data) Update(usr *user.User, ) {
	return
}

func (Dt Data) Leaving() {

}