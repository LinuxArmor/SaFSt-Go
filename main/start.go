package main

import (
	fs2 "../fs"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"log"
	"os"
	"syscall"
)

// Checks if the necessary folders exist before launching
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
	if fs2.Err != nil {
		log.Fatal(fs2.Err)
	}

	defer fs2.Db.Close()

	flag.Parse() // parse the flags from the command line
	if len(flag.Args()) < 1 {
		log.Fatal("Mountpoint unspecified")
	}

	checkDirectory("/usr/local/var/safst")
	checkDirectory("/usr/local/var/safst/files")

	err := fs2.SetUpRoot()

	if err != nil {
		log.Fatal(err)
	}

	c, err := fuse.Mount(
		flag.Arg(0),
		fuse.FSName("SaFSt"),
		fuse.Subtype("SaFSt-FileSystem"),
		fuse.AllowOther(),
		fuse.AllowDev(),
		fuse.LocalVolume(),
		fuse.VolumeName("SaFSt"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, fs2.NewFileSystem(flag.Arg(0)))

	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready

	if c.MountError != nil {
		log.Fatal(err)
	}
}
