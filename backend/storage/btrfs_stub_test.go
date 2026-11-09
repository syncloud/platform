package storage

type BtrfsStub struct {
	mounted bool
	reset   bool
}

func (b *BtrfsStub) Mounted() bool {
	return b.mounted
}

func (b *BtrfsStub) ResetErrors() error {
	b.reset = true
	return nil
}
