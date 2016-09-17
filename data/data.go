package data

import (
	us "proj/user"
	"time"
)

const (
	PRECISION = 10000
	UNIT_PRECISION = 0.1 / (PRECISION / 10)
	latM1 = 0x00000001 //Latitude + 1
	latL1 = 0x00000010 //Latitude - 1
	lonM1 = 0x00000100 //Longitude + 1
	lonL1 = 0x00001000 //Latitude - 1
)

type Data struct {
	Users map[uint64]us.User
	BigData map[float64]map[float64]map[uint64]*us.User
}

func GetData() (dt Data) {
	dt.Users = make(map[uint64]us.User)
	return
}

func (Dt Data) GetUser(id uint64) (*us.User) {
	if usr, ok := Dt.Users[id]; ok == true {
		return &usr
	}
	return nil
}

func (Dt Data) AddUser(id uint64, pos us.Position, tm time.Duration) (*us.User) {
	Dt.Users[id] = us.NewUser(id, pos, tm)
	return &(Dt.Users[id])
}

func (Dt Data) Leaving(usr *us.User) {
	for idx, elem := range usr.Colliders {
		if distance(elem, usr) >= 10.0 { // check la distance.
			go elem.LeaveUser(usr.Uuid)
			go usr.LeaveUser(idx)
		}
	}
}

func (Dt Data)checkAround(usr *us.User, pos us.Position) {
	if mapUser, ok := Dt.BigData[pos.Lat][pos.Lon]; ok {
		for _, tmpUser := range mapUser {
			if _, ok := usr.Colliders[tmpUser.Uuid]; ok == false && distance(tmpUser, usr) < 10.0 {
				tmpUser.JoinUser(usr)
				usr.JoinUser(tmpUser)
			}
		}
	}
}

func (Dt Data) DeleteUserFromBD(usr *us.User) {
	if mapUser, ok := Dt.BigData[usr.Pos.Lat][usr.Pos.Lon]; ok {
		delete(mapUser, usr.Uuid)
	}
	if usr.Flag & latM1 == latM1 {
		if mapUser, ok := Dt.BigData[usr.Pos.Lat + UNIT_PRECISION][usr.Pos.Lon]; ok {
			delete(mapUser, usr.Uuid)
		}
	}
	if usr.Flag & latL1 == latL1 {
		if mapUser, ok := Dt.BigData[usr.Pos.Lat - UNIT_PRECISION][usr.Pos.Lon]; ok {
			delete(mapUser, usr.Uuid)
		}
	}
	if usr.Flag & lonM1 == lonM1 {
		if mapUser, ok := Dt.BigData[usr.Pos.Lat][usr.Pos.Lon + UNIT_PRECISION]; ok {
			delete(mapUser, usr.Uuid)
		}
	}
	if usr.Flag & lonL1 == lonL1 {
		if mapUser, ok := Dt.BigData[usr.Pos.Lat][usr.Pos.Lon - UNIT_PRECISION]; ok {
			delete(mapUser, usr.Uuid)
		}
	}
}

func (Dt Data) UpdateKey(usr *us.User) {
	pos := usr.GetKeys(PRECISION)
	if pos.Lon != usr.Key.Lon || pos.Lat != usr.Key.Lat {
		Dt.DeleteUserFromBD(usr)
	}
	// Check si la key est identique
		// Si c'est pas le cas Update et delete le pointer user de l'ancienne

	Dt.checkAround(usr, pos)
	// Rajouter les options
}