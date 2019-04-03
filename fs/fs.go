// This package includes basic fuse-ops
package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"fmt"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"
)

type SaFStFileSystem struct {
	dir   string
	debug bool
}

func (fs *SaFStFileSystem) Root() (fs.Node, error) {
	return &Dir{[]byte("/")}, nil
}

type Dir struct {
	path []byte
}

func (d Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	f, err := getFile(d.path)
	if err != nil {
		return err
	}
	attr.Atime = time.Now()
	attr.Inode = 0 // Dynamic inode would be used
	attr.Mode = f.Mode
	attr.Gid = f.Gid
	attr.Uid = f.Uid
	attr.Blocks = f.Blocks
	attr.BlockSize = f.BlockSize
	attr.Crtime = f.Crtime
	attr.Flags = f.Flags
	attr.Mtime = f.Mtime
	attr.Nlink = f.Nlink
	attr.Rdev = f.Rdev
	attr.Size = f.Size
	attr.Valid = f.Valid
	return nil
}

// Sets a new debug value which is used for logging actions.
func (fs *SaFStFileSystem) SetDebug(debug bool) {
	fs.debug = debug
}

// Called when a dir is opened
func (fs *SaFStFileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
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

func (fs *SaFStFileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
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

func (fs *SaFStFileSystem) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
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

func (fs *SaFStFileSystem) String() string {
	var debug string

	if debug = "out"; fs.debug {
		debug = ""
	}

	return fmt.Sprintf("SaFStFileSystem - With%s Debug - Directory: %s", debug, fs.dir)
}

func NewFileSystem(dir string) *SaFStFileSystem {
	return &SaFStFileSystem{dir, false}
}
