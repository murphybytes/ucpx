package client

import (
	"os/user"
	"testing"
)

func TestNewFileInfo(t *testing.T) {

	fileInfo, _ := newFileInfo("john@foo.com:/home/john/source")
	if fileInfo == nil {
		t.Error("Expect file info not to be nil")
	}

	if fileInfo.host != "foo.com" {
		t.Error("For host we expected foo.com, but got " + fileInfo.host)
	}

	if fileInfo.user != "john" {
		t.Error("Expected user john, got " + fileInfo.user)
	}

	if fileInfo.path != "/home/john/source" {
		t.Error("Expected /home/john/source but got " + fileInfo.path)
	}

	if fileInfo.port != 9191 {
		t.Errorf("Expected port to be 9191 but got %d", fileInfo.port)
	}
}

func TestNewFileInfoExplicitPort(t *testing.T) {
	fi, e := newFileInfo("john@foo.com:1234:/home/xxx")
	if e != nil {
		t.Error("Expected error to be nil")
	}

	if fi == nil {
		t.Error("Expected non nil file info")
	}

	if fi.host != "foo.com" {
		t.Error("For host we expected foo.com, but got " + fi.host)
	}

	if fi.user != "john" {
		t.Error("Expected user john, got " + fi.user)
	}

	if fi.path != "/home/xxx" {
		t.Error("Expected /home/xxx but got " + fi.path)
	}

	if fi.port != 1234 {
		t.Errorf("Expected port to be 1234 but got %d", fi.port)
	}

}

func TestErrors(t *testing.T) {
	_, e := newFileInfo("john:/home/xxx")
	if e == nil {
		t.Error("Missing host should have caused error")
	}

	_, e = newFileInfo("jam@foo.com::/home/xxx")
	if e == nil {
		t.Error("Missing port should have caused error")
	}

}

func TestNewFileNoUserHost(t *testing.T) {
	fi, e := newFileInfo("/home/xxx")
	u, _ := user.Current()
	if e != nil {
		t.Error("Didn't expect error")
	}
	if fi.user != u.Username {
		t.Error("Incorrect username ")
	}

	if fi.host != "localhost" {
		t.Error("Unexpected host")
	}

	if fi.port != 9191 {
		t.Error("Unexpected port")
	}

	if fi.path != "/home/xxx" {
		t.Error("Path is in error")
	}

}
