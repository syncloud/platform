package storage

import "os"

type Linker struct {
}

func NewLinker() *Linker {
	return &Linker{}
}
func (d *Linker) RelinkDisk(link string, target string) error {

	err := os.Chmod(target, 0o755)
	if err != nil {
		return err
	}

	fi, err := os.Lstat(link)
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		err = os.Remove(link)
		if err != nil {
			return err
		}
	}
	err = os.Symlink(target, link)
	return err
}
