package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"

	//Need mysql
	_ "github.com/go-sql-driver/mysql"
)

func run(dsn, mountpoint string) error {
	// connect to the db
	h, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	handle, err := h.Conn(context.TODO())
	if err != nil {
		panic(err)
	}
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("db"),
		fuse.Subtype("dbfs"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	srv := fs.New(c, nil)
	filesys := &FS{
		handle: handle,
		common: &Common{curInode: 2},
	}
	if err := srv.Serve(filesys); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	dsn := flag.Arg(0)
	mountpoint := flag.Arg(1)

	if err := run(dsn, mountpoint); err != nil {
		log.Fatal(err)
	}
}
