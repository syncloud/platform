package smart

import "github.com/syncloud/platform/storage/model"

const (
	AttributeReallocated          = 5
	AttributePending              = 197
	AttributeOfflineUncorrectable = 198
	AttributeCrcErrors            = 199
	NotAvailable                  = "SMART is not available for this disk"
)

type Output struct {
	Smartctl     Info           `json:"smartctl"`
	SmartStatus  *Status        `json:"smart_status"`
	AtaSmartData *AtaSmartData  `json:"ata_smart_data"`
	Attributes   *AtaAttributes `json:"ata_smart_attributes"`
	PowerOnTime  *PowerOnTime   `json:"power_on_time"`
	Temperature  *Temperature   `json:"temperature"`
	SelfTestLog  *SelfTestLog   `json:"ata_smart_self_test_log"`
}

func (o *Output) Supported() bool {
	return o.SmartStatus != nil
}

func (o *Output) Succeeded() bool {
	return o.Smartctl.ExitStatus&0x7 == 0
}

func (o *Output) Message() string {
	for _, message := range o.Smartctl.Messages {
		if message.Severity == "error" {
			return message.String
		}
	}
	return NotAvailable
}

func (o *Output) Health() model.SmartHealth {
	passed := o.SmartStatus.Passed
	health := model.SmartHealth{
		Available:            true,
		Passed:               &passed,
		Reallocated:          o.Attributes.Raw(AttributeReallocated),
		Pending:              o.Attributes.Raw(AttributePending),
		OfflineUncorrectable: o.Attributes.Raw(AttributeOfflineUncorrectable),
		CrcErrors:            o.Attributes.Raw(AttributeCrcErrors),
	}
	if o.PowerOnTime != nil {
		hours := o.PowerOnTime.Hours
		health.PowerOnHours = &hours
	}
	if o.Temperature != nil {
		current := o.Temperature.Current
		health.Temperature = &current
	}
	health.SelfTest = o.selfTest()
	return health
}

func (o *Output) selfTest() model.SmartTestResult {
	result := model.SmartTestResult{}
	if o.AtaSmartData != nil && o.AtaSmartData.SelfTest != nil {
		status := o.AtaSmartData.SelfTest.Status
		if status.RemainingPercent != nil {
			result.Running = true
			result.RemainingPercent = *status.RemainingPercent
		} else {
			result.LastStatus = status.String
			result.LastPassed = status.Passed
		}
	}
	if o.SelfTestLog != nil && len(o.SelfTestLog.Standard.Table) > 0 {
		last := o.SelfTestLog.Standard.Table[0]
		hours := last.LifetimeHours
		result.LastType = last.Type.String
		result.LastStatus = last.Status.String
		result.LastPassed = last.Status.Passed
		result.LastHours = &hours
	}
	return result
}
