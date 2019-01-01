package main

import (
	"debug/elf"
	"fmt"
)

func check_permissions(elf_file *elf.File) int {
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

func main() {
	var e *elf.File
	var err error
	e, err = elf.Open("a.out")
	fmt.Println(err)
	check_permissions(e)
}
