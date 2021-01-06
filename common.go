package main

type Common struct {
	curInode uint64
}

func (c *Common) getInode() uint64 {
	c.curInode++
	return c.curInode
}
