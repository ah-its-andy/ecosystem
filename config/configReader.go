package config

type ConfigReader interface {
	GetSection(uri string) ConfigSection
}
