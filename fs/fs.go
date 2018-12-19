package fs

import "github.com/hanwen/go-fuse/fuse"

type FileSystem struct {
	fuse.RawFileSystem
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		RawFileSystem: fuse.NewDefaultRawFileSystem(),
	}
}
