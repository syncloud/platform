package backup

type AppBackup interface {
	List() ([]string, error)
	Create(app string)
	Restore(app string, file string)
}
