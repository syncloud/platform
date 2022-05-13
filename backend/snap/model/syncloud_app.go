package model

type SyncloudApp struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Required    string `json:"required"`
	Ui          string `json:"ui"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}
