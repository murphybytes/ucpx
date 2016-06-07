package common

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
)

// Create returns an open file for writing.  If a relative path is
// passed in, an absolute path will be created by appending the
// home dir belonging to user identified by userName
func Create(path string, userName string) (f io.WriteCloser, e error) {
	var usr *user.User
	if usr, e = user.Lookup(userName); e != nil {
		return
	}

	// if the path is not absolute join to users home dir
	if !filepath.IsAbs(path) {
		path = filepath.Join(usr.HomeDir, path)
	}

	return os.Create(path)
}
