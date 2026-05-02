package config

type CustomProxyEntry struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Https    bool   `json:"https"`
	Authelia bool   `json:"authelia"`
}

type CustomProxy struct {
	db *Db
}

func NewCustomProxy(db *Db) *CustomProxy {
	return &CustomProxy{db: db}
}

func (c *CustomProxy) Add(name string, host string, port int, https bool, authelia bool) error {
	httpsInt := 0
	if https {
		httpsInt = 1
	}
	autheliaInt := 0
	if authelia {
		autheliaInt = 1
	}
	_, err := c.db.Exec("INSERT OR REPLACE INTO custom_proxy VALUES (?, ?, ?, ?, ?)", name, host, port, httpsInt, autheliaInt)
	return err
}

func (c *CustomProxy) Remove(name string) error {
	_, err := c.db.Exec("DELETE FROM custom_proxy WHERE name = ?", name)
	return err
}

func (c *CustomProxy) List() ([]CustomProxyEntry, error) {
	db := c.db.Open()
	defer db.Close()
	rows, err := db.Query("select name, host, port, https, authelia from custom_proxy")
	if err != nil {
		return nil, err
	}
	entries := make([]CustomProxyEntry, 0)
	defer rows.Close()
	for rows.Next() {
		var entry CustomProxyEntry
		var httpsInt int
		var autheliaInt int
		if err := rows.Scan(&entry.Name, &entry.Host, &entry.Port, &httpsInt, &autheliaInt); err != nil {
			return entries, err
		}
		entry.Https = httpsInt != 0
		entry.Authelia = autheliaInt != 0
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
