// This package includes basic fuse-ops
package fs

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

const URDONLY_MODE = os.ModeDir | 0700

// Returns the root directory
func (fs *SaFStFileSystem) Root() (fs.Node, error) {
	return &Dir{[]byte("/")}, nil
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

func convertPathError(err *os.PathError) error {
	if os.IsNotExist(err) {
		return fuse.ENOENT
	} else if os.IsPermission(err) {
		return fuse.EPERM
	} else if os.IsExist(err) {
		return fuse.EEXIST
	} else {
		return err
	}
}

// Checks if the permissions match
// fuid, guid: the uid of the file and the accessor respectively
// fgid, ggid: the gid of the file and the accessor respectively
// fmode, gmode: the uid of the file and the accessor respectively
// rwx: 1 to read 2 to write 4 to execute, can be combined together
func checkPerms(fuid, guid uint32, fgid, ggid uint32, mode os.FileMode, rwx uint8) (bul bool) {
	if guid == 0 {
		return true
	}

	bul = true

	if fuid == guid {
		r, w, x := getRWX(0, mode)
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
		r, w, x := getRWX(1, mode)
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
		r, w, x := getRWX(2, mode)
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

// Represents a directory in the file system
type Dir struct {
	path []byte
}

// Implements mknod inside the folder
func (d Dir) Mknod(ctx context.Context, req *fuse.MknodRequest) (fs.Node, error) {
	npath := path.Join(string(d.path), req.Name)
	dat, err := GetFile(d.path)

	if err != nil {
		return nil, err
	}

	if checkPerms(dat.Uid, req.Uid, dat.Gid, req.Gid, dat.Mode, 2) && req.Mode.IsRegular() {
		err = syscall.Mknod(npath, 0700, int(req.Rdev))
		if err != nil {
			return nil, err
		}

		var attr fuse.Attr

		attr.Rdev = req.Rdev
		attr.Atime = time.Now()
		attr.Mode = req.Mode
		attr.Gid = req.Gid
		attr.Uid = req.Uid
		attr.Mtime = time.Now()
		attr.Ctime = time.Now()
		attr.Inode = 0

		f, err := os.Stat(npath)

		if err != nil {
			return nil, convertPathError(err.(*os.PathError))
		}

		file := f.Sys().(*syscall.Stat_t)
		attr.BlockSize = uint32(file.Blksize)
		attr.Blocks = uint64(file.Blocks)
		attr.Nlink = uint32(file.Nlink)
		attr.Size = uint64(file.Size)

		err = PutFile([]byte(npath), attr)

		if err != nil {
			return nil, err
		}
		return &File{[]byte(npath)}, nil
	} else {
		return nil, fuse.EPERM
	}
}

// Implements mkdir inside the folder
func (d Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	dat, err := GetFile(d.path)

	if err != nil {
		return nil, err
	}

	if checkPerms(dat.Uid, req.Uid, dat.Gid, req.Gid, dat.Mode, 2) {
		err = os.Mkdir(path.Join(FileFolder, string(d.path), req.Name), URDONLY_MODE)

		if err != nil {
			return nil, convertPathError(err.(*os.PathError))
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
		attr.Mode = os.ModeDir | req.Mode
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
		return &Dir{npath}, nil
	}
	return nil, fuse.EPERM
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
	file, err := GetFile(p)
	if err != nil {
		return nil, err
	}
	if file.Mode.IsDir() {
		return &Dir{p}, nil
	}
	return &File{p}, nil
}

// Returns the attributes of the folder
func (d Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	f, err := GetFile(d.path)
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

// Implements remove inside the folder
func (d Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	p := []byte(path.Join(string(d.path), req.Name))
	node, err := GetFile(p)

	if err != nil {
		return err
	}

	if checkPerms(node.Uid, req.Uid, node.Gid, req.Gid, node.Mode, 2) {
		err = os.Remove(path.Join(FileFolder, string(p)))

		if err != nil {
			return convertPathError(err.(*os.PathError))
		}

		err = DeleteFile(p)

		if err != nil {
			return err
		}

		return nil
	}
	return fuse.EPERM

}

type File struct {
	path []byte
}

// Returns the attributes of the file
func (fi File) Attr(ctx context.Context, attr *fuse.Attr) error {
	f, err := GetFile(fi.path)
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
