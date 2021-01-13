package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

//TableDir is a directory containing all the tables in a database
type TableDir struct {
	fs    *FS
	name  string
	query *File
}

var _ fs.Node = (*TableDir)(nil)

func (d *TableDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("[d *TableDir Attr]")
	a.Inode = d.fs.common.getInode()
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (d *TableDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("[d *TableDir Lookup] name=", name)
	if name == "query" {
		d.query = &File{}
		d.query.content.Store("asdf")
	}
	return nil, syscall.ENOENT
}

func (d *TableDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *TableDir ReadDirAll]")
	// make a file for querying data
	content := make([]fuse.Dirent, 0)
	content = append(content, fuse.Dirent{Inode: d.fs.common.getInode(), Name: "query", Type: fuse.DT_File})
	return content, nil
}
