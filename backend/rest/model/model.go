package model

type Access struct {
	RelayEnabled bool    `json:"relay_enabled"`
	Ipv4         *string `json:"ipv4,omitempty"`
	Ipv4Enabled  bool    `json:"ipv4_enabled"`
	Ipv4Public   bool    `json:"ipv4_public"`
	AccessPort   *int    `json:"access_port,omitempty"`
	Ipv6Enabled  bool    `json:"ipv6_enabled"`
}

func (a Access) Ipv4Manual() bool {
	return !a.RelayEnabled && a.Ipv4Enabled
}

func (a Access) Ipv4Active() bool {
	return a.RelayEnabled || a.Ipv4Enabled
}

func (a Access) Ipv4PublicDirect() bool {
	return a.Ipv4Public && !a.RelayEnabled
}

type MailRelay struct {
	Enabled bool `json:"enabled"`
}

type MailRelayCredentials struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type RedirectInfoResponse struct {
	Domain string `json:"domain"`
}

type BackupCreateRequest struct {
	App string `json:"app"`
}

type BackupRestoreRequest struct {
	File string `json:"file"`
}

type BackupRemoveRequest struct {
	File string `json:"file"`
}

type StorageActivatePartitionRequest struct {
	Device string `json:"device"`
	Format bool   `json:"format"`
}

type StorageActivateDisksRequest struct {
	Devices []string `json:"devices"`
	Format  bool     `json:"format"`
}

type EventTriggerRequest struct {
	Event string `json:"event"`
}

type UserInfo struct {
	Admin    bool   `json:"admin"`
	Username string `json:"username"`
}

type Response struct {
	Success            bool                 `json:"success"`
	Message            string               `json:"message,omitempty"`
	Code               string               `json:"code,omitempty"`
	Data               *interface{}         `json:"data,omitempty"`
	ParametersMessages *[]ParameterMessages `json:"parameters_messages,omitempty"`
}
