package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"time"
)

type Record struct {
	Content string    `json:"content" bson:"content"`
	Date    time.Time `json:"date" bson:"date"`
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("参数不够")
	}

	usr, err := user.Current()
	panicErr(err)
	path := usr.HomeDir + "/.workRecord-rc"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	panicErr(err)

	buf, err := ioutil.ReadAll(file)
	panicErr(err)

	var list []Record
	if len(buf) > 0 {
		err = json.Unmarshal(buf, &list)
		panicErr(err)
	}

	if os.Args[1] == "-l" {
		for i, v := range list {
			fmt.Printf("%d: %v %v\n", i, v.Date, v.Content)
		}
		return
	}

	if os.Args[1] == "-d" {
		err := os.Remove(path)
		panicErr(err)
		return
	}

	r := Record{os.Args[1], time.Now()}
	list = append(list, r)

	buf, err = json.Marshal(list)
	panicErr(err)

	_, err = file.WriteAt(buf, 0)
	panicErr(err)
}
