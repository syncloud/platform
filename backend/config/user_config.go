package config

import (
	"database/sql"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

const DbTrue = "true"
const DbFalse = "false"
const OldBoolTrue = "True"
const OldBoolFalse = "False"
const DefaultRedirectDomain = "syncloud.it"

type UserConfig struct {
	file          string
	oldConfigFile string
}

var OldConfig string
var DefaultConfigDb string

func init() {
	OldConfig = fmt.Sprintf("%s/user_platform.cfg", os.Getenv("SNAP_COMMON"))
	DefaultConfigDb = fmt.Sprintf("%s/platform.db", os.Getenv("SNAP_DATA"))
}

func NewUserConfig(file string, oldConfigFile string) (*UserConfig, error) {
	config := &UserConfig{
		file:          file,
		oldConfigFile: oldConfigFile,
	}
	err := config.ensureDb()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *UserConfig) ensureDb() error {
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

func (c *UserConfig) migrate() {

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

func (c *UserConfig) SetWebSecretKey(key string) {
	c.Upsert("platform.web_secret_key", key)
}

func (c *UserConfig) initDb() error {
	db := c.open()
	defer db.Close()

	initDbSql := "create table config (key varchar primary key, value varchar)"
	_, err := db.Exec(initDbSql)
	if err != nil {
		return fmt.Errorf("unable to init db (%s): %s", c.file, err)
	}
	return nil
}

func (c *UserConfig) open() *sql.DB {
	db, err := sql.Open("sqlite3", c.file)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (c *UserConfig) SetRedirectDomain(domain string) {
	c.Upsert("redirect.domain", domain)
}

func (c *UserConfig) UpdateRedirectApiUrl(apiUrl string) {
	c.Upsert("redirect.api_url", apiUrl)
}

func (c *UserConfig) SetUserEmail(userEmail string) {
	c.Upsert("redirect.user_email", userEmail)
}

func (c *UserConfig) SetUserUpdateToken(userUpdateToken string) {
	c.Upsert("redirect.user_update_token", userUpdateToken)
}

func (c *UserConfig) GetRedirectDomain() string {
	return c.Get("redirect.domain", DefaultRedirectDomain)
}

func (c *UserConfig) GetRedirectApiUrl() string {
	return fmt.Sprintf("https://api.%s", c.GetRedirectDomain())
}

func (c UserConfig) GetUpnp() bool {
	result := c.Get("platform.upnp", DbTrue)
	return c.toBool(result)
}

func (c *UserConfig) IsRedirectEnabled() bool {
	result := c.Get("platform.redirect_enabled", DbFalse)
	return c.toBool(result)
}

func (c *UserConfig) GetExternalAccess() bool {
	result := c.Get("platform.external_access", DbFalse)
	return c.toBool(result)
}

func (c *UserConfig) SetRedirectEnabled(enabled bool) {
	c.Upsert("platform.redirect_enabled", c.fromBool(enabled))
}

func (c *UserConfig) SetActivated() {
	c.Upsert("platform.activated", DbTrue)
}

func (c *UserConfig) IsActivated() bool {
	return c.toBool(c.Get("platform.activated", DbFalse))
}

func (c *UserConfig) IsCertbotStaging() bool {
	return c.toBool(c.Get("certbot.staging", DbFalse))
}

func (c *UserConfig) SetDomain(domain string) {
	c.Upsert("platform.domain", domain)
}

func (c *UserConfig) GetDomain() *string {
	return c.GetOrNil("platform.domain")
}

func (c *UserConfig) setDeprecatedUserDomain(domain string) {
	c.Upsert("platform.user_domain", domain)
}

func (c *UserConfig) getDeprecatedUserDomain() *string {
	return c.GetOrNil("platform.user_domain")
}

func (c *UserConfig) UpdateDomainToken(token string) {
	c.Upsert("platform.domain_update_token", token)
}

func (c *UserConfig) Upsert(key string, value string) {
	db := c.open()
	defer db.Close()
	_, err := db.Exec("INSERT OR REPLACE INTO config VALUES (?, ?)", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *UserConfig) Delete(key string) {
	db := c.open()
	defer db.Close()
	_, err := db.Exec("DELETE FROM config WHERE key = ?", key)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *UserConfig) GetOrNil(key string) *string {
	db := c.open()
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

func (c *UserConfig) List() map[string]string {
	db := c.open()
	defer db.Close()
	rows, err := db.Query("select key, value from config")
	values := make(map[string]string, 0)
	if err != nil {
		log.Println(err)
		return values
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var value string

		if err := rows.Scan(&key, &value); err != nil {
			log.Println("Unable to scan results:", err)
			continue
		}
		values[key] = value
	}
	return values
}

func (c *UserConfig) Get(key string, defaultValue string) string {
	value := c.GetOrNil(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *UserConfig) toBool(dbValue string) bool {
	return dbValue == DbTrue
}

func (c *UserConfig) fromBool(value bool) string {
	if value {
		return DbTrue
	} else {
		return DbFalse
	}
}

func (c *UserConfig) GetDkimKey() *string {
	return c.GetOrNil("dkim_key")
}

func (c *UserConfig) GetDomainUpdateToken() *string {
	return c.GetOrNil("platform.domain_update_token")
}

func (c *UserConfig) SetExternalAccess(enabled bool) {
	c.Upsert("platform.external_access", c.fromBool(enabled))
}

func (c *UserConfig) SetUpnp(enabled bool) {
	c.Upsert("platform.upnp", c.fromBool(enabled))
}

func (c *UserConfig) SetPublicIp(publicIp string) {
	c.Upsert("platform.public_ip", publicIp)
}

func (c *UserConfig) DeletePublicIp() {
	c.Delete("platform.public_ip")
}

func (c *UserConfig) SetManualCertificatePort(manualCertificatePort int) {
	c.Upsert("platform.manual_certificate_port", strconv.Itoa(manualCertificatePort))
}

func (c *UserConfig) SetManualAccessPort(manualAccessPort int) {
	c.Upsert("platform.manual_access_port", strconv.Itoa(manualAccessPort))
}

func (c *UserConfig) GetCustomDomain() *string {
	return c.GetOrNil("platform.custom_domain")
}

func (c *UserConfig) SetCustomDomain(domain string) {
	c.Upsert("platform.custom_domain", domain)
}

func (c *UserConfig) GetDeviceDomain() string {
	result := "localhost"
	if c.IsRedirectEnabled() {
		domain := c.GetDomain()
		if domain != nil {
			result = *domain
		} else {
			userDomain := c.getDeprecatedUserDomain()
			if userDomain != nil {
				result = fmt.Sprintf("%s.%s", *userDomain, c.GetRedirectDomain())
			}
		}
	} else {
		customDomain := c.GetCustomDomain()
		if customDomain != nil {
			result = *customDomain
		}
	}
	return result
}
