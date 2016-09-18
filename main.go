package main

import (
	"fmt"
	"os"
	mn "proj/manager"
)

func main() {
	var er error
	var file *os.File

	if len(os.Args) != 2 {
		fmt.Println("Please usage ./a.out path_file.csv")
		return
	}
	if file, er = os.Open(os.Args[1]); er != nil {
		fmt.Println("Error: ", er);
		return
	}
	defer file.Close()
	if er := mn.Manager(file); er != nil {
		fmt.Println("Error on Manager: ", er)
	}
	return
}

