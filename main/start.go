package main

import (
	"../fs"
	"flag"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"os"
)

func main() {
	flag.Parse() // parse the flags from the command line
	if len(flag.Args()) < 1 {
		log.Fatal("You need to specify a mountpoint!")
	}
	file, err := os.Open("/usr/local/var/safst")
	if err != nil {
		mkdirerr := os.Mkdir("/usr/local/var/safst", 0700)
		if mkdirerr != nil {
			log.Fatal("Cannot open and create /usr/local/var/safst")
		}
	}
	if file != nil {
		err := file.Close()
		if err != nil {
			log.Fatal("Couldn't close /usr/local/var/safst")
		}
	}
	filesys := fs.NewFileSystem("/usr/local/var/safst")
	filesys.SetDebug(true)
	nfs := pathfs.NewPathNodeFs(filesys, pathfs.PathNodeFsOptions{}) // create a nodefs
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil) // create a server from the nodefs

	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	log.Println("Mounting The FileSystem...")
	server.Serve() // mount our filesystem
}
