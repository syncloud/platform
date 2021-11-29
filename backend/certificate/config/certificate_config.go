package config

type GeneratorSystemConfig interface {
	SslCertificateFile() string
	SslKeyFile() string
}
