package common

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func getPath(filePath, userName string) (path string, e error) {
	var usr *user.User
	if usr, e = user.Lookup(userName); e != nil {
		return
	}
	path = filePath
	// if the path is not absolute join to users home dir
	if !filepath.IsAbs(path) {
		path = filepath.Join(usr.HomeDir, path)
	}

	return
}

// Create returns an open file for writing.  If a relative path is
// passed in, an absolute path will be created by appending the
// home dir belonging to user identified by userName
func Create(path string, userName string) (f io.WriteCloser, e error) {

	if path, e = getPath(path, userName); e != nil {
		return
	}

	return os.Create(path)
}

// Open returns an open file for reading
func Open(path string, userName string) (f io.ReadCloser, e error) {

	if path, e = getPath(path, userName); e != nil {
		return
	}

	return os.Open(path)
}
