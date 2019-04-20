package analyzing

import (
	"debug/elf"
	"fmt"
)

func CheckPermissions(elf_file *elf.File) int {
	var libs []string
	var syms []elf.ImportedSymbol

	var err error
	libs, err = elf_file.ImportedLibraries()
	syms, err = elf_file.ImportedSymbols()
	fmt.Println(err)
	fmt.Printf("%v\n", libs)
	fmt.Printf("syms: %v\n", syms)
	return 0
}
