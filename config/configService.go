package config

type ConfigService struct {
	sectionMap map[string]ConfigSection
}

func (cfg *ConfigService) GetConfig(uri string) ConfigSection {
	if ele, ok := cfg.sectionMap[uri]; ok {
		return ele
	}
	return EmptyConfigSection
}
