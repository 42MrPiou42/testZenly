package main

import (
	"fmt"
	_ "proj/data"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	us "proj/user"
	"time"
	"strconv"
	"proj/data"
)

const (
	UUID = 0
	LATITUDE = 1
	LONGITUDE = 2
	TIME = 3
	BASE = 10
	BYTESIZE = 64
)

func parseRecord(rec []string) (id uint64, pos us.Position, tm time.Duration, er error) {
	if id, er = strconv.ParseUint(rec[UUID], BASE, BYTESIZE ); er != nil {
		return
	}
	if pos.Lat, er = strconv.ParseFloat(rec[LATITUDE], BYTESIZE); er != nil {
		return
	}
	if pos.Lon, er = strconv.ParseFloat(rec[LONGITUDE], BYTESIZE); er != nil {
		return
	}
	if tm, er = time.ParseDuration(rec[TIME] + "ns"); er != nil {
		return
	}
	return
}

func start(f *os.File) {
	title := true
	r := csv.NewReader(bufio.NewReader(f))
	Dt := data.GetData()

	for {
		if title == true {
			title = false
		} else {
			record, er := r.Read()
			if er == io.EOF {
				break
			}
			if len(record) != 4 {
				fmt.Println("csv entry not well formated: ", record)
				continue
			}
			id, pos, tm, er := parseRecord(record)
			if er != nil {
				fmt.Println("error parsing record: ", er)
				continue
			}
			usr := Dt.GetUser(id);
			if usr == nil {
				usr = Dt.AddUser(id, pos, tm)
			} else {
				usr.Pos = pos
				Dt.Leaving(usr)
			}
			Dt.UpdateKey(usr)
		}
	}
	fmt.Println("THE END :D")
	return
}

func main() {
	var er error
	var file *os.File

	if len(os.Args) != 2 {
		fmt.Println("Please usage ./a.out path_file.csv")
		return
	}
	if file, er = os.Open(os.Args[1]); er != nil {
		fmt.Println("Error: ", er);
	}
	defer file.Close()
	start(file)
	return
}

