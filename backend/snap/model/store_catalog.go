package model

type StoreCatalog struct {
	Apps []StoreCatalogApp `json:"apps"`
}

type StoreCatalogApp struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}
