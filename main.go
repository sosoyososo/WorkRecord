package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
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
			nt := v.Date.Format("01-02 15:04")
			fmt.Printf("%d:\t%v\t%v\n", i, nt, v.Content)
		}
		return
	}

	if os.Args[1] == "-d" {
		if len(os.Args) > 3 {
			startStr := os.Args[2]
			start, err := strconv.Atoi(startStr)
			panicErr(err)
			endStr := os.Args[3]
			end, err := strconv.Atoi(endStr)
			panicErr(err)

			var tmpList []Record
			for i, v := range list {
				if i < start || i > end {
					tmpList = append(tmpList, v)
				}
			}
			list = tmpList
		} else if len(os.Args) > 2 {
			indexStr := os.Args[2]
			index, err := strconv.Atoi(indexStr)
			panicErr(err)
			var tmpList []Record
			for i, v := range list {
				if i != index {
					tmpList = append(tmpList, v)
				}
			}
			list = tmpList
		} else {
			err := os.Remove(path)
			panicErr(err)
		}
	} else {
		r := Record{os.Args[1], time.Now()}
		list = append(list, r)
	}

	buf, err = json.Marshal(list)
	panicErr(err)

	l, err := file.WriteAt(buf, 0)
	panicErr(err)

	err = file.Truncate(int64(l))
	panicErr(err)
}
