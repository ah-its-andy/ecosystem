package config

type ConfigServiceBuilder struct {
	configService *ConfigService
	readers       []ConfigReader
}

func NewConfigServiceBuilder(readers []ConfigReader) *ConfigServiceBuilder {
	configBuilder := &ConfigServiceBuilder{
		configService: &ConfigService{
			sectionMap: make(map[string]ConfigSection),
		},
		readers: make([]ConfigReader, len(readers)),
	}
	copy(configBuilder.readers, readers)
	return configBuilder
}

func (builder *ConfigServiceBuilder) AddConfig(uri string) *ConfigServiceBuilder {
	subSections := make([]ConfigSection, 0)
	for _, reader := range builder.readers {
		tmp := reader.GetSection(uri)
		if tmp != nil {
			subSections = append(subSections, tmp)
		}
	}

	if len(subSections) == 0 {
		builder.configService.sectionMap[uri] = EmptyConfigSection
	} else if len(subSections) == 1 {
		builder.configService.sectionMap[uri] = subSections[0]
	} else {
		builder.configService.sectionMap[uri] = NewMultiConfigSection(subSections)
	}

	return builder
}

func (builder *ConfigServiceBuilder) Build() *ConfigService {
	return builder.configService
}
