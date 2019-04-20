package main

import (
	"bufio"
	"fmt"
	"github.com/golang/leveldb"
	"os"
	"strings"
)

func main() {
	var db *leveldb.DB

	var funct, per string
	var inp string
	var arr []string

	db, _ = leveldb.Open("/home/alwin/go/src/SaFSt/libdb", nil)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("func per: ")
	inp, _ = reader.ReadString('\n')
	for inp != "\n" {
		arr = strings.Split(strings.Trim(inp, "\n"), " ")
		funct = arr[0]
		per = arr[1]
		key := funct + "@libc.so.6"
		_ = db.Set([]byte(key), []byte(per), nil)
		fmt.Println("func lib per: ")
		inp, _ = reader.ReadString('\n')
	}
	_ = os.Stdin.Close()
	_ = db.Close()

}
