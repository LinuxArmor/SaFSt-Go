package fs

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
)

type FileSystem struct {
	pathfs.FileSystem
	dir string
}

func (fs *FileSystem) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	dir, err := ioutil.ReadDir(fs.dir + name)
	if err != nil {
		log.Println(err)
		i, _ := strconv.Atoi(err.Error())
		return nil, fuse.Status(i)
	}
	c := make([]fuse.DirEntry, len(dir))
	for i, entry := range dir {

		stat, _ := entry.Sys().(*syscall.Stat_t)
		c[i] = fuse.DirEntry{Mode: uint32(stat.Mode), Name: entry.Name(), Ino: stat.Ino}
	}
	return c, fuse.OK
}

func (fs *FileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	stat, err := os.Stat(fs.dir + name)
	if err != nil {
		log.Println(err)
		i, _ := strconv.Atoi(err.Error())
		return nil, fuse.Status(i)
	}
	attr := fuse.ToAttr(stat)
	attr.Uid = context.Uid
	attr.Gid = context.Gid
	attr.Owner = context.Owner
	return attr, fuse.OK
}

func NewFileSystem(dir string) *FileSystem {
	return &FileSystem{pathfs.NewDefaultFileSystem(), dir}
}
