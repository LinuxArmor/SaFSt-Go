package fs

import (
	"bazil.org/fuse"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

var mutex sync.Mutex = sync.Mutex{} // use a mutex to prevent reading and writing at the same time
var dbfolder, filestats string = "/usr/local/var/safst", dbfolder + "/files"

func getFile(path []byte) (fuse.Attr, error) {
	mutex.Lock()
	defer mutex.Unlock()

	db, err := leveldb.OpenFile(filestats, nil)

	if err != nil || db == nil {
		return fuse.Attr{}, err
	}

	defer db.Close()

	file, geterr := db.Get(path, &(opt.ReadOptions{}))

	if geterr != nil {
		return fuse.Attr{}, geterr
	}

	var f fuse.Attr

	marshalerr := json.Unmarshal(file, f)

	if marshalerr != nil {
		return fuse.Attr{}, marshalerr
	}

	return f, nil
}
