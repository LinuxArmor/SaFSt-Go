package main

import (
	"bazil.org/fuse"
	fs2 "bazil.org/fuse/fs"
	"fmt"
	"os"
	"path"
	"testing"
)

const (
	testFolder = "/tmp/entry"
)

func startFileSystem(c *fuse.Conn) {
	err := fs2.Serve(c, NewFileSystem(testFolder))
	if err != nil {
		return
	}
}

func TestMain(m *testing.M) {
	if os.Getuid() != 0 {
		panic("Tests need to be run on root!")
	}
	if err := os.Mkdir(testFolder, os.ModeDir|0700); err != nil && err != os.ErrExist {
		panic("Couldn't create the test folder!")
	}

	checkDirectory(DbFolder)
	checkDirectory(FileFolder)

	err := SetUpRoot()

	if err != nil {
		panic(err)
	}

	c, err := fuse.Mount(
		testFolder,
		fuse.FSName("SaFSt"),
		fuse.Subtype("SaFSt-FileSystem"),
		fuse.AllowOther(),
		fuse.LocalVolume(),
		fuse.VolumeName("SaFSt"))

	if err != nil {
		panic(err)
	}

	go startFileSystem(c)

	<-c.Ready

	if c.MountError != nil {
		panic(err)
	}

	r := m.Run()

	c.Close()

	err = fuse.Unmount(testFolder)

	if err != nil {
		fmt.Println("Couldn't unmount the test folder!")
	}

	err = os.Remove(testFolder)

	if err != nil {
		fmt.Println("Couldn't remove test folder!")
	}

	os.Exit(r)
}

func TestDir_Attr(t *testing.T) {
	fi, err := os.Stat(testFolder)

	if err != nil {
		t.Errorf("Can't check the attributes of the test folder: %s", err)
	}

	fi2, err := os.Stat(FileFolder)

	if err != nil {
		t.Errorf("Can't check the attributes of the file folder: %s", err)
	}

	if fi.Mode() != fi2.Mode() {
		t.Logf("Test Folder Mode: %d\nFile Folder Mode: %d", fi.Mode(), fi2.Mode())
		t.Errorf("The test folder and the file folder do not share the file mode!")
	}

	t.Logf("Successfully checked the attributes!")
}

func TestDir_Mkdir(t *testing.T) {
	newFolder := path.Join(testFolder, "new")

	if err := os.Mkdir(newFolder, os.ModeDir|0777); err != nil {
		t.Errorf("Couldn't create new folder in the FileSystem: %s", err)
	}

	if fi, err := os.Stat(newFolder); err != nil {
		t.Errorf("Error getting attributes about the new folder: %s", err)
	} else if fi.Mode() != os.ModeDir|0777 {
		t.Error("The new folder doesn't have the file mode given")
	}

	if fi, err := os.Stat(path.Join(FileFolder, "new")); err != nil {
		t.Errorf("Errors getting attributes about the new folder in the file folder: %s", err)
	} else if fi.Mode() != os.ModeDir|0700 {
		t.Errorf("The new folder in the file folder doesn't have the file folder file mode")
	}
}
