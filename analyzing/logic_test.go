package analyzing

import (
	"debug/elf"
	"fmt"
)

func main() {
	err := InitDB()
	file, err := elf.Open("a.out")
	fmt.Print(err)
	fmt.Println(CheckPermissions(file)) // should return 72 because socket is unknown (64) and printf is 8
}
