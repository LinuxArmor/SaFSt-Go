package analyzing

import (
	"debug/elf"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

func CheckPermissions(elf_file *elf.File) int {
	//var libs []string
	var syms []elf.ImportedSymbol
	var ret int

	var err error
	//libs, err = elf_file.ImportedLibraries() in case we need list of imported libraries
	db, err := leveldb.OpenFile("/usr/local/var/safst/perm", nil)
	syms, err = elf_file.ImportedSymbols()
	fmt.Println(err)

	for i := 0; i < len(syms); {
		cur_sym := syms[i]
		if cur_sym.Name == "__libc_start_main" {
			i++
			continue
		}
		cur_func := cur_sym.Name + "@" + cur_sym.Library
		fmt.Println(cur_func)
		var byte_per []byte
		var per int

		byte_per, _ = db.Get([]byte(cur_func), nil)
		if len(byte_per) == 0 {
			per = 64
		} else {
			per, _ = strconv.Atoi(string(byte_per))
		}
		ret = ret | per
		i++
	}
	return ret
}
