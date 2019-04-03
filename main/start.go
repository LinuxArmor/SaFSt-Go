package main

import (
	"flag"
	"log"
	"os"
	"syscall"
)

func checkDirectory(path string) {
	s, err := os.Stat(path)
	if os.IsNotExist(err) { // Check if the folder exists
		mkdirerr := os.Mkdir(path, 0700) // if not create it
		if mkdirerr != nil {
			log.Fatalf("Cannot open and create required folder: %s", path)
		}
		return
	}
	if err != nil { // Any other error should not happen but just in case
		log.Fatalf("Cannot open and create required folder: %s", path)
		return
	}
	if s.Mode().Perm() != 0700 { // Permissions needs to be as that for optimal security
		if chmoderr := os.Chmod(path, os.ModeDir|0700); chmoderr != nil {
			log.Fatalf("Cannot change permissions to  %s", path)
		}
	}

	if fi := s.Sys().(*syscall.Stat_t); fi.Uid != 0 || fi.Gid != 0 { // Permissions needs to be as that for optimal security
		if chownerr := os.Chown(path, 0, 0); chownerr != nil {
			log.Fatalf("Cannot change ownership of  %s", path)
		}
	}
}

func main() {
	flag.Parse() // parse the flags from the command line
	if len(flag.Args()) < 1 {
		log.Fatal("Mountpoint unspecified")
	}

	checkDirectory("/usr/local/var/safst")
	checkDirectory("/usr/local/var/safst/files")
}
