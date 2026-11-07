package model

type SmartHealth struct {
	Available            bool            `json:"available"`
	Message              string          `json:"message,omitempty"`
	Passed               *bool           `json:"passed,omitempty"`
	Reallocated          *int64          `json:"reallocated,omitempty"`
	Pending              *int64          `json:"pending,omitempty"`
	OfflineUncorrectable *int64          `json:"offline_uncorrectable,omitempty"`
	CrcErrors            *int64          `json:"crc_errors,omitempty"`
	PowerOnHours         *int64          `json:"power_on_hours,omitempty"`
	Temperature          *int64          `json:"temperature,omitempty"`
	SelfTest             SmartTestResult `json:"self_test"`
}
