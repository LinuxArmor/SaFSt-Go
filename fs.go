// This package includes basic fuse-ops
package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

// The file system
type SaFStFileSystem struct {
	dir string
}

// Returns the root directory
func (fs *SaFStFileSystem) Root() (fs.Node, error) {
	return &Dir{[]byte("/")}, nil
}

// Represents a directory in the file system
type Dir struct {
	path []byte
}

// Checks the RWX based of the mode and the person trying to access
// who: 0 for the owner of the file
//      1 for the group of the owner of the file
//      2 for others
// returns 3 bools representing read, write and execute respectively
func getRWX(who uint8, mode os.FileMode) (bool, bool, bool) {
	switch who {
	case 0: // owner
		return mode&(1<<8) == 1, mode&(1<<7) == 1, mode&(1<<6) == 1
	case 1: // group
		return mode&(1<<5) == 1, mode&(1<<4) == 1, mode&(1<<3) == 1
	case 2: // other
		return mode&(1<<2) == 1, mode&(1<<1) == 1, mode&1 == 1
	default:
		return false, false, false
	}
}

// Checks if the permissions match
// fuid, guid: the uid of the file and the accessor respectively
// fgid, ggid: the gid of the file and the accessor respectively
// fmode, gmode: the uid of the file and the accessor respectively
// rwx: 1 to read 2 to write 4 to execute, can be combined together
func checkPerms(fuid, guid uint32, fgid, ggid uint32, fmode, gmode os.FileMode, rwx uint8) (bul bool) {
	if guid == 0 {
		return true
	}

	bul = false

	if fuid == guid {
		r, w, x := getRWX(0, gmode)
		if rwx&1 == 1 {
			bul = bul && r
		}
		if rwx&2 == 1 {
			bul = bul && w
		}
		if rwx&4 == 1 {
			bul = bul && x
		}
	} else if fgid == ggid {
		r, w, x := getRWX(1, gmode)
		if rwx&1 == 1 {
			bul = bul && r
		}
		if rwx&2 == 1 {
			bul = bul && w
		}

		if rwx&4 == 1 {
			bul = bul && x
		}
	} else {
		r, w, x := getRWX(2, gmode)
		if rwx&1 == 1 {
			bul = bul && r
		}
		if rwx&2 == 1 {
			bul = bul && w
		}

		if rwx&4 == 1 {
			bul = bul && x
		}
	}

	return
}

// Implements mkdir inside the folder
func (d Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	dat, err := getFile(d.path)

	if err != nil {
		return nil, err
	}

	if checkPerms(dat.Uid, req.Uid, dat.Gid, req.Gid, dat.Mode, req.Mode, 2) {
		log.Println(path.Join(FileFolder, string(d.path), req.Name))
		err = os.Mkdir(path.Join(FileFolder, string(d.path), req.Name), os.ModeDir|0700)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		npath := []byte(path.Join(string(d.path), req.Name))

		var attr fuse.Attr
		s, err := os.Stat(path.Join(FileFolder, string(d.path), req.Name))

		if err != nil {
			log.Println(err)
			return nil, err
		}

		f := s.Sys().(*syscall.Stat_t)
		attr.Atime = time.Now()
		attr.Inode = 0 // Dynamic inode would be used
		attr.Mode = os.ModeDir | os.FileMode(f.Mode)
		attr.Gid = req.Gid
		attr.Uid = req.Uid
		attr.Blocks = uint64(f.Blocks)
		attr.BlockSize = uint32(f.Blksize)
		attr.Crtime = time.Unix(f.Ctim.Unix())
		attr.Mtime = time.Unix(f.Mtim.Unix())
		attr.Nlink = uint32(f.Nlink)
		attr.Rdev = uint32(f.Rdev)
		attr.Size = uint64(f.Size)

		err = PutFile(npath, attr)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		log.Println("Successfully Mkdir")
		return Dir{npath}, nil
	}
	return nil, syscall.EACCES
}

// Implements ls inside the folder
func (d Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	npath := path.Join(FileFolder, string(d.path))
	entries, err := ioutil.ReadDir(npath)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	var dirents []fuse.Dirent
	for _, entry := range entries {
		stat := entry.Sys().(*syscall.Stat_t)
		dirents = append(dirents, fuse.Dirent{stat.Ino, fuse.DirentType(entry.Mode()), entry.Name()})
	}
	return dirents, nil
}

// Implements lookup for files inside the folder
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

// Returns the attributes of the folder
func (d Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	f, err := getFile(d.path)
	if err != nil {
		return err
	}
	attr.Atime = f.Atime
	attr.Inode = 0 // Dynamic inode would be used
	attr.Mode = os.ModeDir | f.Mode
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

// Returns the attributes of the file
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