package storage

type DiskStorage interface {
	Format(device string)
	BootExtend()
}
