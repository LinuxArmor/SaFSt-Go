// This package includes basic fuse-ops
package fs

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
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

func (fs *SaFStFileSystem) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	panic("implement me")
}

func (fs *SaFStFileSystem) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	panic("implement me")
}

func (fs *SaFStFileSystem) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) GetXAttr(name string, attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) ListXAttr(name string, context *fuse.Context) (attributes []string, code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	panic("implement me")
}

func (fs *SaFStFileSystem) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	panic("implement me")
}

func (fs *SaFStFileSystem) OnMount(nodeFs *pathfs.PathNodeFs) {
	panic("implement me")
}

func (fs *SaFStFileSystem) OnUnmount() {
	panic("implement me")
}

func (fs *SaFStFileSystem) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	panic("implement me")
}

func (fs *SaFStFileSystem) StatFs(name string) *fuse.StatfsOut {
	panic("implement me")
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
