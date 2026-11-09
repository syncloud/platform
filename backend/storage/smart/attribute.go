package smart

type Attribute struct {
	Id   int          `json:"id"`
	Name string       `json:"name"`
	Raw  AttributeRaw `json:"raw"`
}
