package ecosystem

type ConfigSection interface {
	GetValue(key string) (string, bool)
	MustGetValue(key string) string
	Keys() []string
}

type ConfigReader interface {
	GetSection(uri string) ConfigSection
}

type ConfigService struct {
	sectionMap map[string]ConfigSection
}

func (cfg *ConfigService) GetConfig(uri string) ConfigSection {
	if ele, ok := cfg.sectionMap[uri]; ok {
		return ele
	}
	return EmptyConfigSection
}

type ConfigServiceBuilder struct {
	configService *ConfigService
	readers       []ConfigReader
}

func (builder *ConfigServiceBuilder) AddConfig(uri string) *ConfigServiceBuilder {
	results := &MultiConfigSection{
		sections: make([]ConfigSection, 0),
	}
	for _, reader := range builder.readers {
		tmp := reader.GetSection(uri)
		if tmp != nil {
			results.sections = append(results.sections, tmp)
		}
	}
	if len(results.sections) == 0 {
		builder.configService.sectionMap[uri] = EmptyConfigSection
	} else {
		builder.configService.sectionMap[uri] = results
	}

	return builder
}
