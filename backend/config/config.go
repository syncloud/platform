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
	file           string
	oldConfigFile  string
	redirectDomain string
	redirectUrl    string
}

var OldConfig string
var DefaultConfigDb string

func init() {
	OldConfig = fmt.Sprintf("%s/user_platform.cfg", os.Getenv("SNAP_COMMON"))
	DefaultConfigDb = fmt.Sprintf("%s/platform.db", os.Getenv("SNAP_COMMON"))
}

func New(file string, oldConfigFile string, redirectDomain string, redirectUrl string) (*PlatformUserConfig, error) {
	config := &PlatformUserConfig{
		file:           file,
		oldConfigFile:  oldConfigFile,
		redirectDomain: redirectDomain,
		redirectUrl:    redirectUrl,
	}
	err := config.ensureDb()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *PlatformUserConfig) ensureDb() error {
	_, err := os.Stat(c.file)
	if os.IsNotExist(err) {
		err = c.initDb()
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(c.oldConfigFile)
	if err == nil {
		c.migrate()
	}
	return nil
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
			c.Upsert(fmt.Sprintf("%s.%s", section, key), dbValue)
		}
	}
	c.SetWebSecretKey(uuid.New().String())
	err = os.Rename(c.oldConfigFile, fmt.Sprintf("%s.bak", c.oldConfigFile))
	if err != nil {
		log.Println("Cannot backup old config: ", c.oldConfigFile, err)
	}
}

func (c *PlatformUserConfig) SetWebSecretKey(value string) {
	c.Upsert("platform.web_secret_key", value)
}

func (c *PlatformUserConfig) initDb() error {
	db := c.open()
	defer db.Close()

	initDbSql := "create table config (key varchar primary key, value varchar)"
	_, err := db.Exec(initDbSql)
	if err != nil {
		return fmt.Errorf("unable to init db (%s): %s", c.file, err)
	}
	return nil
}

func (c *PlatformUserConfig) open() *sql.DB {
	db, err := sql.Open("sqlite3", c.file)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (c *PlatformUserConfig) UpdateRedirectDomain(domain string) {
	c.Upsert("redirect.domain", domain)
}

func (c *PlatformUserConfig) UpdateRedirectApiUrl(apiUrl string) {
	c.Upsert("redirect.api_url", apiUrl)
}

func (c *PlatformUserConfig) SetUserEmail(userEmail string) {
	c.Upsert("redirect.user_email", userEmail)
}

func (c *PlatformUserConfig) SetUserUpdateToken(userUpdateToken string) {
	c.Upsert("redirect.user_update_token", userUpdateToken)
}

func (c *PlatformUserConfig) GetRedirectDomain() string {
	return c.Get("redirect.domain", c.redirectDomain)
}

func (c *PlatformUserConfig) GetRedirectApiUrl() string {
	return c.Get("redirect.api_url", c.redirectUrl)
}

func (c PlatformUserConfig) GetUpnp() bool {
	result := c.Get("platform.upnp", DbTrue)
	return c.toBool(result)
}

func (c *PlatformUserConfig) IsRedirectEnabled() bool {
	result := c.Get("platform.redirect_enabled", DbFalse)
	return c.toBool(result)
}

func (c *PlatformUserConfig) GetExternalAccess() bool {
	result := c.Get("platform.external_access", DbFalse)
	return c.toBool(result)
}

func (c *PlatformUserConfig) SetRedirectEnabled(enabled bool) {
	c.Upsert("platform.redirect_enabled", c.fromBool(enabled))
}

func (c *PlatformUserConfig) UpdateUserDomain(domain string) {
	c.Upsert("platform.user_domain", domain)
}

func (c *PlatformUserConfig) UpdateDomainToken(token string) {
	c.Upsert("platform.domain_update_token", token)
}

func (c *PlatformUserConfig) Upsert(key string, value string) {
	db := c.open()
	defer db.Close()
	_, err := db.Exec("INSERT OR REPLACE INTO config VALUES (?, ?)", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *PlatformUserConfig) Get(key string, defaultValue string) string {
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
