// This package includes basic fuse-ops
package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"path"
	"syscall"
	"time"
)

type SaFStFileSystem struct {
	dir string
}

func (fs *SaFStFileSystem) Root() (fs.Node, error) {
	return &Dir{[]byte("/")}, nil
}

type Dir struct {
	path []byte
}

func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	npath := path.Join(FileFolder, string(d.path))
	entries, err := ioutil.ReadDir(npath)

	if err != nil {
		return nil, err
	}
	var dirents []fuse.Dirent
	for _, entry := range entries {
		stat := entry.Sys().(*syscall.Stat_t)
		dirents = append(dirents, fuse.Dirent{stat.Ino, fuse.DirentType(entry.Mode()), entry.Name()})
	}
	return dirents, nil
}

func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	p := []byte(path.Join(string(d.path), name))
	file, err := getFile(p)
	if err != nil {
		return nil, err
	}
	if file.Mode.IsDir() {
		return Dir{p}, nil
	}
	return File{p}, nil
}

func (d Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Printf("Dir Attributes of %s\n", d.path)
	f, err := getFile(d.path)
	if err != nil {
		return err
	}
	attr.Atime = f.Atime
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

type File struct {
	path []byte
}

func (fi File) Attr(ctx context.Context, attr *fuse.Attr) error {
	f, err := getFile(fi.path)
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

func NewFileSystem(dir string) *SaFStFileSystem {
	return &SaFStFileSystem{dir}
}
