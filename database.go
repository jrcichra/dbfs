package main

import (
	"context"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

//DatabaseDir is a directory containing all the databases in the db handle
type DatabaseDir struct {
	fs   *FS
	name string
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
	return &TableDir{fs: d.fs, name: name}, nil
}

func (d *DatabaseDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *DatabaseDir ReadDirAll]")
	// Generate a directory per line of 'show databases'
	log.Println("Going to use", d.name)
	_, err := d.fs.handle.ExecContext(context.TODO(), "USE "+d.name)
	if err != nil {
		panic(err)
	}
	res, err := d.fs.handle.QueryContext(context.TODO(), "SHOW TABLES")
	if err != nil {
		panic(err)
	}

	dirs := make([]fuse.Dirent, 0)

	var table string

	for res.Next() {
		res.Scan(&table)
		// Make a new directory for this table
		log.Println("found table:", table)
		dir := fuse.Dirent{Inode: d.fs.common.getInode(), Name: table, Type: fuse.DT_Dir}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}
