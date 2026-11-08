package smart

type AtaSmartData struct {
	SelfTest     *SelfTestState   `json:"self_test"`
	Capabilities *AtaCapabilities `json:"capabilities"`
}
