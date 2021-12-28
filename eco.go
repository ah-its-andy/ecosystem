package ecosystem

import (
	"fmt"
	"path/filepath"

	"github.com/ah-its-andy/ecosystem/config"
	"github.com/ah-its-andy/ecosystem/db"
	ecohttp "github.com/ah-its-andy/ecosystem/http"
	"github.com/ah-its-andy/ecosystem/logging"
	"github.com/ah-its-andy/ecosystem/redis"
)

type EcoSystem struct {
	ExecutePath string
	ServiceName string
	ServiceCfg  config.ConfigSection

	configService *config.ConfigService
	dsn           db.DbConnDsn
	redis         *redis.RedisDsn
	smartLogger   logging.Logger
	app           ecohttp.Application
}

func (eco *EcoSystem) GetConfigService() *config.ConfigService {
	return eco.configService
}

func (eco *EcoSystem) GetDsn() db.DbConnDsn {
	return eco.dsn
}

func (eco *EcoSystem) GetRedisDsn() *redis.RedisDsn {
	return eco.redis
}

func (eco *EcoSystem) GetLogger(name string) logging.Logger {
	return eco.smartLogger
}

func (eco *EcoSystem) GetApplication(name string) ecohttp.Application {
	return eco.app
}

func (eco *EcoSystem) GetDefaultConfigDiscoveryPathList() []string {
	return []string{
		filepath.Join(eco.ExecutePath, "conf"),
		fmt.Sprintf("/etc/zfy/conf.d/%s", eco.ServiceName),
	}
}
