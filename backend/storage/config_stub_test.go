package storage

type ConfigStub struct {
	diskDir string
}

func (c *ConfigStub) ExternalDiskDir() string {
	return c.diskDir
}
