// This package includes basic fuse-ops
package fs

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"io/ioutil"
	"log"
	"os"
	"syscall"
)

type FileSystem struct {
	pathfs.FileSystem
	dir   string
	debug bool
}

// Sets a new debug value which is used for logging actions.
func (fs *FileSystem) SetDebug(debug bool) {
	fs.debug = debug
}

// Called when a dir is opened
func (fs *FileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	dir, err := ioutil.ReadDir(fs.dir + name)

	if err != nil {
		if fs.debug {
			log.Println("Received an error while opening a directory: " + name)
			log.Println(err)
		}
		return nil, fuse.ToStatus(err)
	}
	if fs.debug {
		log.Println("Opened directory: " + name)
	}
	c := make([]fuse.DirEntry, len(dir))
	for i, entry := range dir {
		stat, _ := entry.Sys().(*syscall.Stat_t)
		c[i] = fuse.DirEntry{Mode: uint32(stat.Mode), Name: entry.Name(), Ino: stat.Ino}
	}
	if fs.debug {
		log.Println(c)
	}
	return c, fuse.OK
}

func (fs *FileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	stat, err := os.Stat(fs.dir + name)
	if err != nil {
		if fs.debug {
			log.Println("Received an error while getting attributes: " + name)
			log.Println(err)
		}
		return nil, fuse.ToStatus(err)
	}
	if fs.debug {
		log.Println("Got attributes about: " + name)
		log.Println(stat)
	}
	attr := fuse.ToAttr(stat)
	attr.Uid = context.Uid
	attr.Gid = context.Gid
	attr.Owner = context.Owner
	return attr, fuse.OK
}

func (fs *FileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	file, err := os.Open(name)

	if err != nil {
		if fs.debug {
			log.Println("Received an error while opening a file: " + name)
			log.Println(err)
		}
		return nil, fuse.ToStatus(err)
	}
	if fs.debug {
		log.Println("Opened a file: " + name)
		log.Println(file)
	}
	return nodefs.NewLoopbackFile(file), fuse.OK
}

func NewFileSystem(dir string) *FileSystem {
	return &FileSystem{pathfs.NewDefaultFileSystem(), dir, false}
}
