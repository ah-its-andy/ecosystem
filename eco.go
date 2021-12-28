package ecosystem

type EcoSystem struct {
	ServiceName   string
	configService *ConfigService
	dsn           DbConnDsn
	redis         *RedisDsn
}

func (eco *EcoSystem) GetConfigService() *ConfigService {
	return eco.configService
}

func (eco *EcoSystem) GetDsn() DbConnDsn {
	return eco.dsn
}

func (eco *EcoSystem) GetRedisDsn() *RedisDsn {
	return eco.redis
}
