package backup

type File struct {
	Path string `json:"path"`
	File string `json:"file"`
}

type Auto struct {
	Auto string `json:"auto"`
	Day  int    `json:"day"`
	Hour int    `json:"hour"`
}

type AppBackup interface {
	List() ([]File, error)
	Create(app string)
	Restore(file string)
}
