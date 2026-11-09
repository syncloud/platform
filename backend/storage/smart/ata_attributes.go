package smart

type AtaAttributes struct {
	Table []Attribute `json:"table"`
}

func (a *AtaAttributes) Raw(id int) *int64 {
	if a == nil {
		return nil
	}
	for _, attribute := range a.Table {
		if attribute.Id == id {
			value := attribute.Raw.Value
			return &value
		}
	}
	return nil
}
