package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

//DatabasesDir is a directory containing all the databases in the db handle
type DatabasesDir struct {
	fs *FS
}

var _ fs.Node = (*DatabasesDir)(nil)

func (d *DatabasesDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("[d *DatabasesDir Attr]")
	a.Inode = 4
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (d *DatabasesDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("[d *DatabasesDir Lookup] name=", name)
	if name == "afile" {
		d.fs.aFile = &File{}
		d.fs.aFile.content.Store("asdf")
		return d.fs.aFile, nil
	}
	return nil, syscall.ENOENT
}

func (d *DatabasesDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *DatabasesDir ReadDirAll]")
	return dirDirs, nil
}
