package config

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"go.uber.org/zap"
	"log"
	"strconv"
)

const DefaultConfigDb = "/var/snap/platform/current/platform.db"

const True = "true"
const False = "false"

type Db struct {
	file   string
	logger *zap.Logger
}

func NewDb(file string, logger *zap.Logger) *Db {
	return &Db{file: file, logger: logger}
}

func (c *Db) File() string {
	return c.file
}

func (c *Db) Open() *sql.DB {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", c.file)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (c *Db) Init() error {
	db := c.Open()
	defer db.Close()

	_, err := db.Exec("create table config (key varchar primary key, value varchar)")
	if err != nil {
		return fmt.Errorf("unable to init db (%s): %s", c.file, err)
	}
	return nil
}

func (c *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	db := c.Open()
	defer db.Close()
	return db.Exec(query, args...)
}

func (c *Db) Upsert(key string, value string) {
	db := c.Open()
	defer db.Close()
	_, err := db.Exec("INSERT OR REPLACE INTO config VALUES (?, ?)", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Db) Delete(key string) {
	db := c.Open()
	defer db.Close()
	_, err := db.Exec("DELETE FROM config WHERE key = ?", key)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Db) GetOrNilString(key string) *string {
	db := c.Open()
	defer db.Close()
	var value string
	err := db.QueryRow("select value from config where key = ?", key).Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	}
	return &value
}

func (c *Db) Get(key string, defaultValue string) string {
	value := c.GetOrNilString(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *Db) GetBool(key string, defaultValue bool) bool {
	return c.Get(key, fromBool(defaultValue)) == True
}

func (c *Db) UpsertBool(key string, value bool) {
	c.Upsert(key, fromBool(value))
}

func fromBool(value bool) string {
	if value {
		return True
	}
	return False
}

func (c *Db) GetOrDefaultString(key string, defaultValue string) string {
	return c.Get(key, defaultValue)
}

func (c *Db) GetStringOrError(key string) (string, error) {
	value := c.GetOrNilString(key)
	if value == nil {
		return "", fmt.Errorf("%s is not found", key)
	}
	return *value, nil
}

func (c *Db) GetOrNilInt(key string) *int {
	value := c.GetOrNilString(key)
	if value == nil {
		return nil
	}
	i, err := strconv.Atoi(*value)
	if err != nil {
		return nil
	}
	return &i
}

func (c *Db) GetOrDefaultInt(key string, defaultValue int) int {
	value := c.GetOrNilInt(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *Db) GetOrNilInt64(key string) *int64 {
	value := c.GetOrNilString(key)
	if value == nil {
		return nil
	}
	i, err := strconv.ParseInt(*value, 10, 32)
	if err != nil {
		return nil
	}
	return &i
}

func (c *Db) GetOrDefaultInt64(key string, defaultValue int64) int64 {
	value := c.GetOrNilInt64(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *Db) List() map[string]string {
	db := c.Open()
	defer db.Close()
	values := make(map[string]string)
	rows, err := db.Query("select key, value from config")
	if err != nil {
		c.logger.Error("config query", zap.Error(err))
		return values
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var value string
		if err = rows.Scan(&key, &value); err != nil {
			c.logger.Error("Unable to scan results", zap.Error(err))
			continue
		}
		values[key] = value
	}
	return values
}
