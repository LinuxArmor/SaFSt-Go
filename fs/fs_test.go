package fs_test

import (
	fs "."
	"bazil.org/fuse"
	fs2 "bazil.org/fuse/fs"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	"testing"
)

const (
	testFolder = "/tmp/entry"
)

func checkDirectory(path string) {
	s, err := os.Stat(path)
	if os.IsNotExist(err) { // Check if the folder exists
		mkdirerr := os.Mkdir(path, 0700) // if not create it
		if mkdirerr != nil {
			log.Fatalf("Cannot open and create required folder: %s", path)
		}
		return
	}
	if err != nil { // Any other error should not happen but just in case
		log.Fatalf("Cannot open and create required folder: %s", path)
		return
	}
	if s.Mode().Perm() != 0700 { // Permissions needs to be as that for optimal security
		if chmoderr := os.Chmod(path, os.ModeDir|0700); chmoderr != nil {
			log.Fatalf("Cannot change permissions to  %s", path)
		}
	}

	if fi := s.Sys().(*syscall.Stat_t); fi.Uid != 0 || fi.Gid != 0 { // Permissions needs to be as that for optimal security
		if chownerr := os.Chown(path, 0, 0); chownerr != nil {
			log.Fatalf("Cannot change ownership of  %s", path)
		}
	}
}

func startFileSystem(c *fuse.Conn) {
	err := fs2.Serve(c, fs.NewFileSystem(testFolder))
	if err != nil {
		fmt.Printf("Couldn't serve the filesystem / the filesystem was unmounted %s\n", err)
	}
}

func TestMain(m *testing.M) {
	if os.Getuid() != 0 {
		panic("Tests need to be run on root!")
	}

	checkDirectory(fs.DbFolder)
	checkDirectory(fs.FileFolder)

	if err := os.Mkdir(testFolder, os.ModeDir|0700); err != nil && err.(*os.PathError).Err != syscall.EEXIST {
		panic(fmt.Errorf("couldn't create the test folder: %s", err))
	}

	defer func() {
		err := fuse.Unmount(testFolder)
		if err != nil {
			fmt.Printf("Failed to unmount the test folder: %s\n", err)
		}
		err = os.Remove(fs.DbFolder)
		if err != nil {
			fmt.Println("couldn't remove the database folder!")
		}
		err = os.Remove(testFolder)

		if err != nil {
			fmt.Println("Couldn't remove the test folder!")
		}
	}()

	if err := fs.SetUpRoot(); err != nil {
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

	defer c.Close()

	os.Exit(m.Run())
}

func TestDir_Attr(t *testing.T) {
	fi, err := os.Stat(testFolder)

	if err != nil {
		t.Errorf("Can't check the attributes of the test folder: %s", err)
	}

	fi2, err := os.Stat(fs.FileFolder)

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
	} else if fi.Mode() != os.ModeDir|0755 {
		t.Errorf("The new folder doesn't have the file mode given: %s", fi.Mode())
	}

	if fi, err := os.Stat(path.Join(fs.FileFolder, "new")); err != nil {
		t.Errorf("Errors getting attributes about the new folder in the file folder: %s", err)
	} else if fi.Mode() != os.ModeDir|0700 {
		t.Errorf("The new folder in the file folder doesn't have the file folder file mode")
	}

	t.Logf("Successfully created a new directory!")
}
