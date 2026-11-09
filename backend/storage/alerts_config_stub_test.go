package storage

type AlertsConfigStub struct {
	seen []string
}

func (c *AlertsConfigStub) StorageAlertsSeen() []string {
	return c.seen
}

func (c *AlertsConfigStub) SetStorageAlertsSeen(keys []string) {
	c.seen = keys
}
