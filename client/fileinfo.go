package client

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"strings"

	"github.com/murphybytes/ucp/common"
)

type fileInfo struct {
	host  string
	user  string
	port  int
	path  string
	local bool
	read  bool
}

func parseUserHost(userHost string) (user, host string, e error) {
	parts := strings.Split(userHost, "@")
	if len(parts) != 2 {
		e = errors.New("Invalid host specification - " + userHost)
		return
	}

	user = parts[0]
	host = parts[1]
	e = nil
	return
}

// filespec takes the form
// [user@host[:port]:]/path/to/file
// is optional user@host is not supplied /path/to/file is assumed
// to be local with the current user
func newFileInfo(filespec string, read bool) (fi *fileInfo, e error) {
	fi = &fileInfo{
		port: common.DefaultPort,
		read: read,
	}

	parts := strings.Split(filespec, ":")
	if len(parts) == 1 {
		fi.local = true
		fi.path = filespec
		fi.host = "localhost"

		var userInfo *user.User
		if userInfo, e = user.Current(); e != nil {
			return nil, e
		}

		fi.user = userInfo.Username
	}

	if len(parts) == 2 {
		fi.path = parts[1]
		if fi.user, fi.host, e = parseUserHost(parts[0]); e != nil {
			return nil, e
		}
	}

	if len(parts) == 3 {
		fi.path = parts[2]
		if fi.port, e = strconv.Atoi(parts[1]); e != nil {
			e = errors.New("Invalid port format in file spec '" + filespec + "'")
			return nil, e
		}

		if fi.user, fi.host, e = parseUserHost(parts[0]); e != nil {
			return nil, e
		}
	}

	return fi, nil
}

func (fi *fileInfo) getConnectString() (connectString string, e error) {
	if fi.local {
		e = errors.New("Connect string not available for local file info.")
		return
	}
	connectString = fmt.Sprint(fi.host, ":", fi.port)
	return

}
