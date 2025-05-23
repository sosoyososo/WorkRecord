package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

type ActionType int

const (
	ActionTypeHelp ActionType = iota
	ActionTypeAdd
	ActionTypeList
	ActionTypeUpdate
	ActionTypeDel
	ActionTypeCopy
	ActionTypeMax
)

var (
	acList = map[ActionType]func([]string, []Record) (bool, []Record){
		ActionTypeHelp:   handleHelp,
		ActionTypeAdd:    handleAdd,
		ActionTypeList:   handleList,
		ActionTypeUpdate: handleUpdate,
		ActionTypeDel:    handleDel,
		ActionTypeCopy:   handleCopyToClipboard,
	}
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
func handleCopyToClipboard(args []string, list []Record) (bool, []Record) {
	// -c index
	if len(args) < 3 {
		panicErr(errors.New("command format fail, should be : '-c index'"))
	}
	// check index
	indexStr := args[2]
	index, err := strconv.Atoi(indexStr)
	panicErr(err)
	if index < 0 || index >= len(list) {
		panicErr(errors.New("out of index"))
	}
	// copy content to clipboard
	content := list[index].Content
	err = clipboard.WriteAll(content)
	panicErr(err)
	fmt.Println("Content copied to clipboard")
	return false, list

}

func handleDel(args []string, list []Record) (bool, []Record) {
	if len(args) > 3 { // -d indexStart indexEnd //删除多个
		startStr := args[2]
		start, err := strconv.Atoi(startStr)
		panicErr(err)
		endStr := args[3]
		end, err := strconv.Atoi(endStr)
		panicErr(err)

		var tmpList []Record
		for i, v := range list {
			if i < start || i > end {
				tmpList = append(tmpList, v)
			}
		}
		list = tmpList
	} else if len(args) > 2 { //-d index //删除一个
		indexStr := args[2]
		index, err := strconv.Atoi(indexStr)
		panicErr(err)
		var tmpList []Record
		for i, v := range list {
			if i != index {
				tmpList = append(tmpList, v)
			}
		}
		list = tmpList
	} else { //-d //删除全部
		list = []Record{}
	}
	return true, list
}

func handleList(args []string, list []Record) (bool, []Record) {
	for i, v := range list {
		nt := v.Date.Format("01-02 15:04")
		fmt.Printf("%d:\t%v\t%v\n", i, nt, v.Content)
	}
	return false, list
}

func handleAdd(args []string, list []Record) (bool, []Record) {
	r := Record{args[1], time.Now()}
	list = append(list, r)
	return true, list
}

func handleUpdate(args []string, list []Record) (bool, []Record) {
	// -e index freshContent
	if len(args) < 4 {
		panicErr(errors.New("command format fail, should be : '-e index freshContent'"))
	}

	// check index
	indexStr := args[2]
	index, err := strconv.Atoi(indexStr)
	panicErr(err)
	if len(list) <= index {
		panicErr(errors.New("out of index"))
	}

	// update content
	list[index].Content = args[3]
	list[index].Date = time.Now()
	return true, list
}

func handleHelp(args []string, list []Record) (bool, []Record) {
	fmt.Println("Usage:")
	fmt.Println("  workrecord [content]          Add new record")
	fmt.Println("  workrecord -l                 List all records")
	fmt.Println("  workrecord -d [index]         Delete record at index")
	fmt.Println("  workrecord -d [start] [end]   Delete records from start to end")
	fmt.Println("  workrecord -e [index] [text]  Update record at index")
	fmt.Println("  workrecord -c [index]         Copy record content to clipboard")
	fmt.Println("  workrecord -h                 Show this help")
	return false, list
}

func main() {
	// create new file for user or read content from file for user
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

	args := os.Args
	acType := ActionTypeMax

	if len(args) < 2 || args[1] == "-l" { //无参数或者首参数是-l，是添加内容
		acType = ActionTypeList
	} else if args[1] == "-d" { //首参是-d，是删除内容
		acType = ActionTypeDel
	} else if args[1] == "-e" { // 首参是-e，是更改内容
		acType = ActionTypeUpdate
	} else if args[1] == "-c" { // 首参是-c，是复制内容
		acType = ActionTypeCopy
	} else if args[1] == "-h" { // 首参是-h，是帮助
		acType = ActionTypeHelp
	} else { //其余情况都是为了添加内容
		acType = ActionTypeAdd
	}

	ac := acList[acType]
	needChange, list := ac(args, list)
	if !needChange {
		return
	}

	buf, err = json.Marshal(list)
	panicErr(err)

	l, err := file.WriteAt(buf, 0)
	panicErr(err)

	err = file.Truncate(int64(l))
	panicErr(err)
}
