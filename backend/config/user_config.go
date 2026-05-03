package config

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserConfig struct {
	db     *Db
	logger *zap.Logger
}

func NewUserConfig(
	db *Db,
	logger *zap.Logger,
) *UserConfig {
	return &UserConfig{
		db:     db,
		logger: logger,
	}
}

func (c *UserConfig) SetWebSecretKey(key string) {
	c.db.Upsert("platform.web_secret_key", key)
}

func (c *UserConfig) GetWebSecretKey() string {
	return c.db.GetOrDefaultString("platform.web_secret_key", "default")
}

func (c *UserConfig) IsIpv4Enabled() bool {
	return c.db.GetBool("platform.ipv4_enabled", true)
}

func (c *UserConfig) SetIpv4Enabled(enabled bool) {
	c.db.UpsertBool("platform.ipv4_enabled", enabled)
}

func (c *UserConfig) IsIpv4Public() bool {
	return c.db.GetBool("platform.ipv4_public", true)
}

func (c *UserConfig) SetIpv4Public(enabled bool) {
	c.db.UpsertBool("platform.ipv4_public", enabled)
}

func (c *UserConfig) IsIpv6Enabled() bool {
	return c.db.GetBool("platform.ipv6_enabled", true)
}

func (c *UserConfig) SetIpv6Enabled(enabled bool) {
	c.db.UpsertBool("platform.ipv6_enabled", enabled)
}

func (c *UserConfig) IsRedirectEnabled() bool {
	return c.db.GetBool("platform.redirect_enabled", false)
}

func (c *UserConfig) SetRedirectEnabled(enabled bool) {
	c.db.UpsertBool("platform.redirect_enabled", enabled)
}

func (c *UserConfig) SetActivated() {
	c.db.UpsertBool("platform.activated", true)
}

func (c *UserConfig) SetDeactivated() {
	c.db.UpsertBool("platform.activated", false)
}

func (c *UserConfig) IsActivated() bool {
	return c.db.GetBool("platform.activated", false)
}

func (c *UserConfig) IsCertbotStaging() bool {
	return c.db.GetBool("certbot.staging", false)
}

func (c *UserConfig) SetDomain(domain string) {
	c.db.Upsert("platform.domain", domain)
}

func (c *UserConfig) getDomain() *string {
	return c.db.GetOrNilString("platform.domain")
}

func (c *UserConfig) UpdateDomainToken(token string) {
	c.db.Upsert("platform.domain_update_token", token)
}

func (c *UserConfig) GetDkimKey() *string {
	return c.db.GetOrNilString("dkim_key")
}

func (c *UserConfig) SetDkimKey(key *string) {
	if key == nil {
		c.db.Delete("dkim_key")
	} else {
		c.db.Upsert("dkim_key", *key)
	}
}
func (c *UserConfig) GetDomainUpdateToken() *string {
	return c.db.GetOrNilString("platform.domain_update_token")
}

func (c *UserConfig) SetPublicIp(publicIp *string) {
	if publicIp == nil {
		c.db.Delete("platform.public_ip")
	} else {
		c.db.Upsert("platform.public_ip", *publicIp)
	}
}

func (c *UserConfig) GetPublicIp() *string {
	return c.db.GetOrNilString("platform.public_ip")
}

func (c *UserConfig) SetPublicPort(port *int) {
	if port == nil {
		c.db.Delete("platform.manual_access_port")
	} else {
		c.db.Upsert("platform.manual_access_port", strconv.Itoa(*port))
	}
}

func (c *UserConfig) GetPublicPort() *int {
	return c.db.GetOrNilInt("platform.manual_access_port")
}

func (c *UserConfig) GetCustomDomain() *string {
	return c.db.GetOrNilString("platform.custom_domain")
}

func (c *UserConfig) GetBackupAuto() string {
	auto := c.db.GetOrNilString("platform.backup_auto")
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
		c.db.Upsert("platform.backup_auto", auto)
	}
}

func (c *UserConfig) GetBackupAutoDay() int {
	return c.db.GetOrDefaultInt("platform.backup_auto_day", 0)
}

func (c *UserConfig) SetBackupAutoDay(day int) {
	c.db.Upsert("platform.backup_auto_day", strconv.Itoa(day))
}

func (c *UserConfig) GetBackupAutoHour() int {
	return c.db.GetOrDefaultInt("platform.backup_auto_hour", 0)
}

func (c *UserConfig) SetBackupAutoHour(hour int) {
	c.db.Upsert("platform.backup_auto_hour", strconv.Itoa(hour))
}

func (c *UserConfig) GetBackupAppTime(app string, mode string) time.Time {
	value := c.db.GetOrNilInt64(fmt.Sprintf("platform.backup.%s.%s", app, mode))
	if value == nil {
		return time.Time{}
	}
	return time.Unix(*value, 0)
}

func (c *UserConfig) SetBackupAppTime(app string, mode string, time time.Time) {
	c.db.Upsert(fmt.Sprintf("platform.backup.%s.%s", app, mode), strconv.FormatInt(time.Unix(), 10))
}

func (c *UserConfig) SetCustomDomain(domain string) {
	c.db.Upsert("platform.custom_domain", domain)
}

func (c *UserConfig) GetDeviceDomain() string {
	result := "www.localhost"
	if c.IsRedirectEnabled() {
		domain := c.getDomain()
		if domain != nil {
			result = *domain
		}
	} else {
		customDomain := c.GetCustomDomain()
		if customDomain != nil {
			result = *customDomain
		}
	}
	return result
}

func (c *UserConfig) DeviceUrl() string {
	port := c.GetPublicPort()
	domain := c.GetDeviceDomain()
	return ConstructUrl(port, domain)
}

func ConstructUrl(port *int, domain string) string {
	externalPort := ""
	if port != nil && *port != 80 && *port != 443 {
		externalPort = fmt.Sprintf(":%d", *port)
	}
	return fmt.Sprintf("https://%s%s", domain, externalPort)
}

func (c *UserConfig) AppDomain(app string) string {
	return fmt.Sprintf("%s.%s", app, c.GetDeviceDomain())
}

func (c *UserConfig) IsTwoFactorEnabled() bool {
	return c.db.GetBool("platform.two_factor_enabled", false)
}

func (c *UserConfig) SetTwoFactorEnabled(enabled bool) {
	c.db.UpsertBool("platform.two_factor_enabled", enabled)
}

func (c *UserConfig) GetTimezone() string {
	return c.db.Get("platform.timezone", "UTC")
}

func (c *UserConfig) SetTimezone(timezone string) {
	c.db.Upsert("platform.timezone", timezone)
}

func (c *UserConfig) Url(app string) string {
	port := c.GetPublicPort()
	domain := c.GetDeviceDomain()
	return ConstructUrl(port, fmt.Sprintf("%s.%s", app, domain))
}
