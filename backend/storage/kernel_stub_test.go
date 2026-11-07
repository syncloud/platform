package storage

import "github.com/syncloud/platform/storage/kernel"

type KernelStub struct {
	dmesg string
}

func (k *KernelStub) Events() (*kernel.Events, error) {
	return kernel.NewEvents(kernel.ParseLines(k.dmesg)), nil
}
