package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

//DatabaseDir is a directory containing all the databases in the db handle
type DatabaseDir struct {
	fs *FS
}

var _ fs.Node = (*DatabaseDir)(nil)

func (d *DatabaseDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("[d *DatabaseDir Attr]")
	a.Inode = d.fs.common.getInode()
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (d *DatabaseDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("[d *DatabaseDir Lookup] name=", name)
	return nil, syscall.ENOENT
}

func (d *DatabaseDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *DatabaseDir ReadDirAll]")
	return nil, nil
}
