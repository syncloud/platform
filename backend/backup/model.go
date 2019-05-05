package backup

type File struct {
	Path string `json:"path"`
	File string `json:"file"`
}

type AppBackup interface {
	List() ([]File, error)
	Create(app string)
	Restore(file string)
}
