// +build darwin dragonfly freebsd linux netbsd openbsd

package main

import (
	"syscall"
)

// Returns total and free bytes available in a directory.
func DiskSpace(path string) (total, free int64, err error) {
	s := syscall.Statfs_t{}
	err = syscall.Statfs(path, &s)
	if err != nil {
		return
	}
	total = s.Bsize * int64(s.Blocks)
	free = s.Bsize * int64(s.Bfree)
	return
}