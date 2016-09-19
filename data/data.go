package data

import (
	us "proj/user"
	geo "github.com/kellydunn/golang-geo"
	"math"
	_ "fmt"
	"fmt"
)

const (
	UNIT_PRECISION = 0.1 / (us.PRECISION_TRUNC / 10)
	LIMIT_PRECISION = UNIT_PRECISION - math.SmallestNonzeroFloat64
	UNION_DISTANCE_KM = 0.01
	LATM1 = uint8(1) // Latitude + UNIT_PRECISION
	LATL1 = uint8(1) << 1 // Latitude - UNIT_PRECISION
	LONM1 = uint8(1) << 2 // Longitude + UNIT_PRECISION
	LONL1 = uint8(1) << 3 // Longitude - UNIT_PRECISION
)

type Data struct {
	Users map[int64]*us.User
	Chunks map[float64]map[float64]map[int64]*us.User
}

func (dt *Data) SetData() {
	dt.Users = make(map[int64]*us.User)
	dt.Chunks = make(map[float64]map[float64]map[int64]*us.User)
}

func (dt *Data) AddUser(id int64, pos *us.Position, tm int64) (*us.User) {
	usr, ok := dt.Users[id]
	if ok == false {
		usr = us.CreateUser(id)
	} else {
		if usr.Pos.GreatCircleDistance(geo.NewPoint(pos.Lat, pos.Lon)) < UNION_DISTANCE_KM {
			return nil
		}
	}
	usr.SetPosition(pos)
	usr.SetTime(tm)
	if ok == false {
		dt.Users[id] = usr
	}
	return usr
}

func (dt *Data) deleteOldChunks(usr *us.User) {
	if mapUser, ok := dt.Chunks[usr.Key.Lat][usr.Key.Lon]; ok {
		delete(mapUser, usr.Uuid)
	}
	if usr.Flag > 0 {
		if usr.Flag & LATM1 == LATM1 {
			usr.Key.Lat += UNIT_PRECISION
		}
		if usr.Flag & LATL1 == LATL1 {
			usr.Key.Lat -= UNIT_PRECISION
		}
		if usr.Flag & LONM1 == LONM1 {
			usr.Key.Lon += UNIT_PRECISION
		}
		if usr.Flag & LONL1 == LONL1 {
			usr.Key.Lon -= UNIT_PRECISION
		}
		if mapUser, ok := dt.Chunks[usr.Key.Lat][usr.Key.Lon]; ok {
			delete(mapUser, usr.Uuid)
		}
	}
}

func (dt *Data) addChunks(usr *us.User) {
	if _, ok := dt.Chunks[usr.Key.Lat][usr.Key.Lon][usr.Uuid]; ok == false {
		if _, ok := dt.Chunks[usr.Key.Lat]; ok == false {
			dt.Chunks[usr.Key.Lat] = make(map[float64]map[int64]*us.User)
		}
		if _, ok := dt.Chunks[usr.Key.Lat][usr.Key.Lon]; ok == false {
			dt.Chunks[usr.Key.Lat][usr.Key.Lon] = make(map[int64]*us.User)
		}
		dt.Chunks[usr.Key.Lat][usr.Key.Lon][usr.Uuid] = usr
	}
	usr.Flag = 0
	switch {
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat + UNIT_PRECISION, usr.Key.Lon + UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LATM1
		usr.Flag |= LONM1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat + UNIT_PRECISION, usr.Key.Lon)) < UNION_DISTANCE_KM:
		usr.Flag |= LATM1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat + UNIT_PRECISION, usr.Key.Lon - UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LATM1
		usr.Flag |= LONL1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat - UNIT_PRECISION, usr.Key.Lon + UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LATL1
		usr.Flag |= LONM1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat - UNIT_PRECISION, usr.Key.Lon)) < UNION_DISTANCE_KM:
		usr.Flag |= LATL1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat - UNIT_PRECISION, usr.Key.Lon - UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LATL1
		usr.Flag |= LONL1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat, usr.Key.Lon + UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LONM1
	case usr.Pos.GreatCircleDistance(geo.NewPoint(usr.Key.Lat, usr.Key.Lon - UNIT_PRECISION)) < UNION_DISTANCE_KM:
		usr.Flag |= LONL1
	default:
		usr.Flag = 0
	}
	pos, _ := us.CreatePosition(usr.Pos.Lng(), usr.Pos.Lat())
	if usr.Flag & LATM1 == LATM1 {
		pos.Lat += UNIT_PRECISION
	}
	if usr.Flag & LATL1 == LATL1 {
		pos.Lat -= UNIT_PRECISION
	}
	if usr.Flag & LONM1 == LONM1 {
		pos.Lon += UNIT_PRECISION
	}
	if usr.Flag & LONL1 == LONL1 {
		pos.Lon -= UNIT_PRECISION
	}
	if _, ok := dt.Chunks[pos.Lat][pos.Lon][usr.Uuid]; ok == false {
		if _, ok := dt.Chunks[pos.Lat]; ok == false {
			dt.Chunks[pos.Lat] = make(map[float64]map[int64]*us.User)
		}
		if _, ok := dt.Chunks[pos.Lat][pos.Lon]; ok == false {
			dt.Chunks[pos.Lat][pos.Lon] = make(map[int64]*us.User)
		}
		dt.Chunks[pos.Lat][pos.Lon][usr.Uuid] = usr
	}
}

func (dt *Data) UpdateChunks(usr *us.User) (error) {
	pos, er := us.CreatePosition(usr.Pos.Lng(), usr.Pos.Lat())
	if er != nil {
		return er
	}
	pos.TruncateMe()
	if usr.CompareKey(pos) == false {
		dt.deleteOldChunks(usr)
	}
	usr.UpdateKey(pos)
	dt.addChunks(usr)
	return nil
}

func (dt *Data) PrintData() {
	fmt.Println("Size nbr Users: ", len(dt.Users))
	fmt.Println("Size nbr Map Lat Key: ", len(dt.Chunks))
	for latkey, mapLon := range dt.Chunks {
		fmt.Println("Lat-Key: ", latkey)
		fmt.Println("Size relation lonKey: ", len(mapLon))
		for lonkey, mapUser := range mapLon {
			fmt.Println("Lon-Key: ", lonkey)
			fmt.Println("Size relation User: ", len(mapUser))
			for id, Usr := range dt.Chunks[latkey][lonkey] {
				fmt.Println("User id: ", id)
				fmt.Println("Nbr Collider: ", len(Usr.Colliders))
			}
		}
	}
}