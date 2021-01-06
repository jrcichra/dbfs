package main

import (
	"context"
	"log"
	"os"

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
	a.Inode = d.fs.common.getInode()
	a.Mode = os.ModeDir | 0o555
	return nil
}

func (d *DatabasesDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("[d *DatabasesDir Lookup] name=", name)
	return &DatabaseDir{fs: d.fs}, nil
}

func (d *DatabasesDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *DatabasesDir ReadDirAll]")

	// Generate a directory per line of 'show databases'
	res, _ := d.fs.handle.Query("SHOW DATABASES")

	dirs := make([]fuse.Dirent, 0)

	var database string
	for res.Next() {
		res.Scan(&database)
		// Make a new directory for this table
		log.Println("found database:", database)
		dir := fuse.Dirent{Inode: d.fs.common.getInode(), Name: database, Type: fuse.DT_Dir}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}
