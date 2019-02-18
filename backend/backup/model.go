package backup

type AppBackup interface {
	List() ([]string, error)
	Create(app string, file string)
	Restore(app string, file string)
}
