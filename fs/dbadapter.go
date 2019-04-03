package fs

import (
	"bazil.org/fuse"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

var mutex sync.Mutex = sync.Mutex{} // use a mutex to prevent reading and writing at the same time
const DbFolder, FileFolder string = "/usr/local/var/safst", DbFolder + "/files"

var Db, Err = leveldb.OpenFile(DbFolder, &opt.Options{})

// Sets up the root folder in the database
func SetUpRoot() error {
	exists, err := Db.Has([]byte("/"), nil)

	if err != nil {
		return err
	}

	if !exists {
		var attr fuse.Attr
		s, err := os.Stat(FileFolder)

		if err != nil {
			return err
		}

		f := s.Sys().(*syscall.Stat_t)
		attr.Atime = time.Unix(f.Atim.Unix())
		attr.Inode = 0 // Dynamic inode would be used
		attr.Mode = os.FileMode(f.Mode)
		attr.Gid = f.Gid
		attr.Uid = f.Uid
		attr.Blocks = uint64(f.Blocks)
		attr.BlockSize = uint32(f.Blksize)
		attr.Crtime = time.Unix(f.Ctim.Unix())
		attr.Mtime = time.Unix(f.Mtim.Unix())
		attr.Nlink = uint32(f.Nlink)
		attr.Rdev = uint32(f.Rdev)
		attr.Size = uint64(f.Size)
		js, err := json.Marshal(attr)

		if err != nil {
			return err
		}

		err = Db.Put([]byte("/"), js, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

// Gets the file at the specfied path and returns its saved attributes.
// The root folder should always be specified.
// Strong Consistency is applied
func getFile(path []byte) (fuse.Attr, error) {
	file, err := Db.Get(path, &(opt.ReadOptions{}))

	if err != nil {
		log.Println(err)
		if err == leveldb.ErrNotFound {
			return fuse.Attr{}, fuse.Errno(syscall.ENOENT)
		}
		return fuse.Attr{}, err
	}

	var f fuse.Attr

	marshalerr := json.Unmarshal(file, f)

	if marshalerr != nil {
		return fuse.Attr{}, marshalerr
	}

	return f, nil
}
