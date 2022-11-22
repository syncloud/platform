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
	"time"
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
	OldConfig = fmt.Sprintf("/var/snap/platform/common/user_platform.cfg")
	DefaultConfigDb = fmt.Sprintf("/var/snap/platform/current/platform.db")
}

func NewUserConfig(file string, oldConfigFile string) *UserConfig {
	return &UserConfig{
		file:          file,
		oldConfigFile: oldConfigFile,
	}
}

func (c *UserConfig) Load() {
	err := c.ensureDb()
	if err != nil {
		log.Fatalln(err)
	}
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
		c.migrateV1()
	}
	c.migrateV2()
	return nil
}

func (c *UserConfig) migrateV1() {

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

func (c *UserConfig) migrateV2() {
	result := c.GetOrNilString("platform.external_access")
	if result == nil {
		return
	}

	c.SetIpv4Public(c.toBool(*result))
	c.Delete("platform.external_access")
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
	c.Delete("redirect.api_url")
	c.Upsert("redirect.domain", domain)
}

func (c *UserConfig) UpdateRedirectApiUrl(apiUrl string) {
	c.Upsert("redirect.api_url", apiUrl)
}

func (c *UserConfig) SetUserEmail(userEmail string) {
	c.Upsert("redirect.user_email", userEmail)
}

func (c *UserConfig) GetUserEmail() *string {
	return c.GetOrNilString("redirect.user_email")
}

func (c *UserConfig) SetUserUpdateToken(userUpdateToken string) {
	c.Upsert("redirect.user_update_token", userUpdateToken)
}

func (c *UserConfig) GetRedirectDomain() string {
	return c.Get("redirect.domain", DefaultRedirectDomain)
}

func (c *UserConfig) GetRedirectApiUrl() string {
	return c.Get("redirect.api_url", fmt.Sprintf("https://api.%s", c.GetRedirectDomain()))
}

func (c *UserConfig) IsIpv4Enabled() bool {
	result := c.Get("platform.ipv4_enabled", DbTrue)
	return c.toBool(result)
}

func (c *UserConfig) SetIpv4Enabled(enabled bool) {
	c.Upsert("platform.ipv4_enabled", c.fromBool(enabled))
}

func (c *UserConfig) IsIpv4Public() bool {
	result := c.Get("platform.ipv4_public", DbTrue)
	return c.toBool(result)
}

func (c *UserConfig) SetIpv4Public(enabled bool) {
	c.Upsert("platform.ipv4_public", c.fromBool(enabled))
}

func (c *UserConfig) IsIpv6Enabled() bool {
	result := c.Get("platform.ipv6_enabled", DbTrue)
	return c.toBool(result)
}

func (c *UserConfig) SetIpv6Enabled(enabled bool) {
	c.Upsert("platform.ipv6_enabled", c.fromBool(enabled))
}

func (c *UserConfig) IsRedirectEnabled() bool {
	result := c.Get("platform.redirect_enabled", DbFalse)
	return c.toBool(result)
}

func (c *UserConfig) GetDeprecatedExternalAccess() bool {
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

func (c *UserConfig) getDomain() *string {
	return c.GetOrNilString("platform.domain")
}

func (c *UserConfig) setDeprecatedUserDomain(domain string) {
	c.Upsert("platform.user_domain", domain)
}

func (c *UserConfig) getDeprecatedUserDomain() *string {
	return c.GetOrNilString("platform.user_domain")
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

func (c *UserConfig) GetOrDefaultInt(key string, defaultValue int) int {
	value := c.GetOrNilInt(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *UserConfig) GetOrDefaultInt64(key string, defaultValue int64) int64 {
	value := c.GetOrNilInt64(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *UserConfig) GetOrNilInt(key string) *int {
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

func (c *UserConfig) GetOrNilInt64(key string) *int64 {
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

func (c *UserConfig) GetOrDefaultString(key string, defaultValue string) string {
	value := c.GetOrNilString(key)
	if value == nil {
		return defaultValue
	}
	return *value
}

func (c *UserConfig) GetOrNilString(key string) *string {
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
	value := c.GetOrNilString(key)
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
	return c.GetOrNilString("dkim_key")
}

func (c *UserConfig) GetDomainUpdateToken() *string {
	return c.GetOrNilString("platform.domain_update_token")
}

func (c *UserConfig) SetPublicIp(publicIp *string) {
	if publicIp == nil {
		c.Delete("platform.public_ip")
	} else {
		c.Upsert("platform.public_ip", *publicIp)
	}
}

func (c *UserConfig) GetPublicIp() *string {
	return c.GetOrNilString("platform.public_ip")
}

func (c *UserConfig) SetPublicPort(port *int) {
	if port == nil {
		c.Delete("platform.manual_access_port")
	} else {
		c.Upsert("platform.manual_access_port", strconv.Itoa(*port))
	}
}

func (c *UserConfig) GetPublicPort() *int {
	return c.GetOrNilInt("platform.manual_access_port")
}

func (c *UserConfig) GetCustomDomain() *string {
	return c.GetOrNilString("platform.custom_domain")
}

func (c *UserConfig) GetBackupAuto() string {
	auto := c.GetOrNilString("platform.backup_auto")
	if auto == nil {
		return "no"
	}
	return *auto
}

func (c *UserConfig) SetBackupAuto(auto string) {
	switch auto {
	case
		"no",
		"backup",
		"restore":
		c.Upsert("platform.backup_auto", auto)
	}
}

func (c *UserConfig) GetBackupAutoDay() int {
	return c.GetOrDefaultInt("platform.backup_auto_day", 0)
}

func (c *UserConfig) SetBackupAutoDay(day int) {
	c.Upsert("platform.backup_auto_day", strconv.Itoa(day))
}

func (c *UserConfig) GetBackupAutoHour() int {
	return c.GetOrDefaultInt("platform.backup_auto_hour", 0)
}

func (c *UserConfig) SetBackupAutoHour(hour int) {
	c.Upsert("platform.backup_auto_hour", strconv.Itoa(hour))
}

func (c *UserConfig) GetBackupAppTime(app string, mode string) time.Time {
	value := c.GetOrNilInt64(fmt.Sprintf("platform.backup.%s.%s", app, mode))
	if value == nil {
		return time.Time{}
	}
	return time.Unix(*value, 0)
}

func (c *UserConfig) SetBackupAppTime(app string, mode string, time time.Time) {
	c.Upsert(fmt.Sprintf("platform.backup.%s.%s", app, mode), strconv.FormatInt(time.Unix(), 10))
}

func (c *UserConfig) SetCustomDomain(domain string) {
	c.Upsert("platform.custom_domain", domain)
}

func (c *UserConfig) GetDeviceDomain() string {
	result := "localhost"
	if c.IsRedirectEnabled() {
		domain := c.getDomain()
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
