package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync/atomic"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
)

type FS struct {
	handle    *sql.DB
	aFile     *File
	databases *DatabasesDir
	common    *Common
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	return &RootDir{fs: f}, nil
}

// RootDir implements both Node and Handle for the root directory.
type RootDir struct {
	fs *FS
}

var _ fs.Node = (*RootDir)(nil)

func (d *RootDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("[d *Dir Attr]")
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}

var _ fs.NodeStringLookuper = (*RootDir)(nil)

func (d *RootDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Println("[d *RootDir Lookup] name=", name)
	if name == "afile" {
		d.fs.aFile = &File{}
		d.fs.aFile.content.Store("asdf")
		return d.fs.aFile, nil
	} else if name == "databases" {
		d.fs.databases = &DatabasesDir{fs: d.fs}
		return d.fs.databases, nil
	}
	return nil, syscall.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "afile", Type: fuse.DT_File},
	{Inode: 3, Name: "databases", Type: fuse.DT_Dir},
}

var _ fs.HandleReadDirAller = (*RootDir)(nil)

func (d *RootDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("[d *RootDir ReadDirAll]")
	return dirDirs, nil
}

type File struct {
	fuse    *fs.Server
	content atomic.Value
	count   uint64
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Println("f *File [Attr]")
	a.Inode = 2
	a.Mode = 0o444
	t := f.content.Load().(string)
	a.Size = uint64(len(t))
	return nil
}

var _ fs.NodeOpener = (*File)(nil)

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("f *File [Open]")
	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

var _ fs.Handle = (*File)(nil)

var _ fs.HandleReader = (*File)(nil)

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("f *File [Read]")
	t := f.content.Load().(string)
	fuseutil.HandleRead(req, resp, []byte(t))
	return nil
}
