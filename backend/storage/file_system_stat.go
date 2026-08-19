package storage

import "syscall"

type FileSystemStat struct {
}

func NewFileSystemStat() *FileSystemStat {
	return &FileSystemStat{}
}

func (f *FileSystemStat) Stat(path string) (totalKB uint64, freeKB uint64, device uint64, err error) {
	var stat syscall.Stat_t
	err = syscall.Stat(path, &stat)
	if err != nil {
		return 0, 0, 0, err
	}
	var statfs syscall.Statfs_t
	err = syscall.Statfs(path, &statfs)
	if err != nil {
		return 0, 0, 0, err
	}
	blockSize := uint64(statfs.Bsize)
	return statfs.Blocks * blockSize / 1024, statfs.Bavail * blockSize / 1024, uint64(stat.Dev), nil
}
