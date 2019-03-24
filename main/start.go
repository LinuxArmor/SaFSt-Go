package main

import (
	"flag"
	"github.com/LinuxArmor/SaFSt-Go/fs"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
)

func main() {
	flag.Parse() // parse the flags from the command line
	if len(flag.Args()) < 1 {
		log.Fatal("You need to specify a mountpoint!")
	}
	filesys := fs.NewFileSystem("/usr/local/var/safst")
	filesys.SetDebug(true)
	nfs := pathfs.NewPathNodeFs(filesys, nil)                        // create a nodefs
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil) // create a server from the nodefs

	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	println("Hello Docker!")
	server.Serve() // mount our filesystem
}
