package user

import (
	"fmt"
	"errors"
	geo "github.com/kellydunn/golang-geo"
	"strconv"
)

const (
	PRECISION_TRUNC = 10000
)

type Position struct {
	Lon float64
	Lat float64
}

type User struct {
	Uuid int64
	Pos *geo.Point
	Key Position // truncaded position
	Time int64
	Colliders map[int64]*User // Sort by Id User
	Flag uint8
}

/*
** Functions for Position
 */

func checkLonAndLat(lon float64, lat float64) (error) {
	if lon < -180.0 || lon > 180.0 {
		return errors.New("Longitude must be on a range -180.0 to 180.0 value: " +
			strconv.FormatFloat(lon, 'E', -1, 64))
	} else if lat < -85.05115 || lat > 85.0 {
		return errors.New("Latitude must be on a range -85.0.. to 85.0 value: " +
			strconv.FormatFloat(lat, 'E', -1, 64))
	}
	return nil
}

func CreatePosition(lon float64, lat float64) (*Position, error) {
	if er := checkLonAndLat(lon, lat); er != nil {
		return nil, er
	}
	pos := &Position{Lon:lon, Lat:lat}
	return pos, nil
}

func (Pos *Position) TruncateMe() {
	Pos.Lon = float64(int(Pos.Lon * PRECISION_TRUNC)) / PRECISION_TRUNC
	Pos.Lat = float64(int(Pos.Lat * PRECISION_TRUNC)) / PRECISION_TRUNC
}

/*
** Functions for User
 */

func CreateUser(id int64) (*User) {
	usr := &User{}

	usr.Uuid = id
	usr.Colliders = make(map[int64]*User)
	return usr
}

func (Usr *User) SetTime(tm int64) {
	Usr.Time = tm
}

func (Usr *User) SetPosition(pos *Position) {
	Usr.Pos = geo.NewPoint(pos.Lat, pos.Lon)
}

func (Usr *User) RemoveCollider(idOther int64) {
	if Usr.Uuid < idOther {
		fmt.Println(Usr.Uuid, " a quitté ", idOther)
	}
	delete(Usr.Colliders, idOther)
}

func (Usr *User) AddCollider(otherUser *User) {
	Usr.Colliders[otherUser.Uuid] = otherUser
	if Usr.Uuid < otherUser.Uuid {
		fmt.Println(Usr.Uuid, " est avec ", otherUser.Uuid)
	}
}

func (Usr *User) CompareKey(pos *Position) (bool) {
	if Usr.Key.Lat == pos.Lat && Usr.Key.Lon == pos.Lon {
		return true
	}
	return false
}

func (User *User) UpdateKey(pos *Position) {
	User.Key.Lon = pos.Lon
	User.Key.Lat = pos.Lat
}

