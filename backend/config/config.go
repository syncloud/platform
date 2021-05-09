package config

import (
	"database/sql"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DbTrue = "true"
var DbFalse = "false"
var OldBoolTrue = "True"
var OldBoolFalse = "False"

type PlatformUserConfig struct {
	file          string
	oldConfigFile string
}

func New(file string, oldConfigFile string) *PlatformUserConfig {
	return &PlatformUserConfig{
		file:          file,
		oldConfigFile: oldConfigFile,
	}
}

func (c *PlatformUserConfig) EnsureDb() {
	_, err := os.Stat(c.file)
	if os.IsNotExist(err) {
		c.initDb()
	}

	_, err = os.Stat(c.oldConfigFile)
	if err == nil {
		c.migrate()
	}

}

func (c *PlatformUserConfig) migrate() {

	oldConfig, err := configparser.NewConfigParserFromFile(c.oldConfigFile)
	if err != nil {
		log.Println("Cannot load config: ", c.oldConfigFile, err)
		return
	}

	for _, section := range oldConfig.Sections() {
		dict, err := oldConfig.Items(section)
		if err != nil {
			log.Println("Cannot read sections config: ", c.oldConfigFile, err)
			return
		}
		for key, value := range dict {
			dbValue := value
			if value == OldBoolTrue || value == OldBoolFalse {
				dbValue = c.fromBool(value == OldBoolTrue)
			}
			c.upsert(fmt.Sprintf("%s.%s", section, key), dbValue)
		}
	}
	c.SetWebSecretKey(uuid.New().String())
	err = os.Rename(c.oldConfigFile, fmt.Sprintf("%s.bak", c.oldConfigFile))
	if err != nil {
		log.Println("Cannot backup old config: ", c.oldConfigFile, err)
	}
}

func (c *PlatformUserConfig) SetWebSecretKey(value string) {
	c.upsert("platform.web_secret_key", value)
}

func (c *PlatformUserConfig) initDb() {
	db := c.open()
	defer db.Close()

	initDbSql := "create table config (key varchar primary key, value varchar)"
	_, err := db.Exec(initDbSql)
	if err != nil {
		log.Printf("%q: %s\n", err, initDbSql)
		return
	}

}

func (c *PlatformUserConfig) open() *sql.DB {
	db, err := sql.Open("sqlite3", c.file)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (c *PlatformUserConfig) UpdateRedirect(domain string, api_url string) {
	c.upsert("redirect.domain", domain)
	c.upsert("redirect.api_url", api_url)
}

func (c *PlatformUserConfig) GetRedirectDomain() string {
	return c.get("redirect.domain", "syncloud.it")
}

func (c *PlatformUserConfig) GetRedirectApiUrl() string {
	return c.get("redirect.api_url", "https://api.syncloud.it")
}

func (c PlatformUserConfig) GetUpnp() bool {
	result := c.get("platform.upnp", DbTrue)
	return c.toBool(result)
}

func (c *PlatformUserConfig) IsRedirectEnabled() bool {
	result := c.get("platform.redirect_enabled", DbFalse)
	return c.toBool(result)
}

func (c *PlatformUserConfig) GetExternalAccess() bool {
	result := c.get("platform.external_access", DbFalse)
	return c.toBool(result)
}

func (c *PlatformUserConfig) upsert(key string, value string) {
	db := c.open()
	defer db.Close()
	log.Printf("setting %s=%s", key, value)
	_, err := db.Exec("INSERT OR REPLACE INTO config VALUES (?, ?)", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *PlatformUserConfig) get(key string, defaultValue string) string {
	db := c.open()
	defer db.Close()
	var value string
	err := db.QueryRow("select value from config where key = ?", key).Scan(&value)
	switch {
	case err == sql.ErrNoRows:
		return defaultValue
	case err != nil:
		log.Fatal(err)
	}
	return value
}

func (c *PlatformUserConfig) toBool(dbValue string) bool {
	return dbValue == DbTrue
}

func (c *PlatformUserConfig) fromBool(value bool) string {
	if value {
		return DbTrue
	} else {
		return DbFalse
	}
}
