package manager

import (
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"fmt"
	"strconv"
	"proj/data"
	us "proj/user"
)
const (
	UUID = 0
	LATITUDE = 1
	LONGITUDE = 2
	TIME = 3
	BASE = 10
	BYTESIZE = 64
)

func parseRecord(rec []string) (id int64, pos *us.Position, tm int64, er error) {
	var lon float64
	var lat float64

	if id, er = strconv.ParseInt(rec[UUID], BASE, BYTESIZE ); er != nil {
		return
	}
	if lat, er = strconv.ParseFloat(rec[LATITUDE], BYTESIZE); er != nil {
		return
	}
	if lon, er = strconv.ParseFloat(rec[LONGITUDE], BYTESIZE); er != nil {
		return
	}
	if pos, er = us.CreatePosition(lon, lat); er != nil {
		return
	}
	if tm, er = strconv.ParseInt(rec[TIME], BASE, BYTESIZE); er != nil {
		return
	}
	return
}

func leavingColliders(usr *us.User) {
	for _, cld := range usr.Colliders {
		if cld.Pos.GreatCircleDistance(usr.Pos) > data.UNION_DISTANCE_KM {
			cld.RemoveCollider(usr.Uuid)
			usr.RemoveCollider(cld.Uuid)
		}
	}
}

func getColliders(mapUser map[int64]*us.User, usr *us.User) {
	for _, oneUser := range mapUser {
		if _, ok := usr.Colliders[oneUser.Uuid]; ok == false && oneUser.Uuid != usr.Uuid &&
			usr.Pos.GreatCircleDistance(oneUser.Pos) <= data.UNION_DISTANCE_KM {
			usr.AddCollider(oneUser)
			oneUser.AddCollider(usr)
		}
	}
}

func browseChunks(usr *us.User, dt *data.Data) {
	if mapUser, ok := dt.Chunks[usr.Key.Lat][usr.Key.Lon]; ok {
		getColliders(mapUser, usr)
	}
	if usr.Flag > 0 {
		pos, _ := us.CreatePosition(usr.Key.Lon, usr.Key.Lat)
		if usr.Flag & data.LATM1 == data.LATM1 {
			pos.Lat += data.UNIT_PRECISION
		}
		if usr.Flag & data.LATL1 == data.LATL1 {
			pos.Lat -= data.UNIT_PRECISION
		}
		if usr.Flag & data.LONM1 == data.LONM1 {
			pos.Lon += data.UNIT_PRECISION
		}
		if usr.Flag & data.LONL1 == data.LONL1 {
			pos.Lon -= data.UNIT_PRECISION
		}
		if mapUser, ok := dt.Chunks[pos.Lat][pos.Lon]; ok {
			getColliders(mapUser, usr)
		}
	}
}

func process(rec []string, dt *data.Data) {
	id, pos, tm, er := parseRecord(rec)
	if er != nil {
		return
	}
	usr := dt.AddUser(id, pos, tm) // Add User if exist, else he was created
	leavingColliders(usr)
	if er = dt.UpdateChunks(usr); er != nil {
		return
	}
	browseChunks(usr, dt)
	return
}

func Manager(f *os.File) (error) {
	title := true
	r := csv.NewReader(bufio.NewReader(f))
	dt := &data.Data{}

	dt.SetData()
	for i := 0; i >= 0; i++ {
		record, er := r.Read()
		if er == io.EOF {
			dt.PrintData()
			return nil
		} else if er != nil {
			return er
		}
		if title == true {
			title = false
		} else {
			if len(record) != 4 {
				fmt.Println("csv entry not well formated: ", record)
				continue
			}
			process(record, dt)
		}
	}
	return nil
}