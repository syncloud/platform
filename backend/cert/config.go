package cert

type GeneratorSystemConfig interface {
	SslCertificateFile() string
	SslKeyFile() string
}
