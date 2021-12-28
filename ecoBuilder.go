package ecosystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ecoInstance *EcoSystem

func GetEco() *EcoSystem {
	return ecoInstance
}

type EcoSystemBuilder struct {
	eco *EcoSystem
}

func NewEcoSystemBuilder() *EcoSystemBuilder {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	configReader := &YmlConfigReader{
		fileDiscoveryPathList: []string{
			filepath.Join(execPath, "conf"),
		},
	}
	config := configReader.GetSection("application.yml")

	if serviceName, ok := config.GetValue("serviceName"); !ok {
		panic("服务没有配置serviceName")
	} else {
		ecoInstance := &EcoSystem{
			ServiceName: serviceName,
		}
		ecoBuilder := &EcoSystemBuilder{
			eco: ecoInstance,
		}
		return ecoBuilder
	}

}

func (builder *EcoSystemBuilder) Build() {
	ecoInstance = builder.eco
}

func (builder *EcoSystemBuilder) UseConfigService(f func(*ConfigServiceBuilder)) {
	exec, err := os.Executable()
	if err != nil {
		panic(errors.Unwrap(err))
	}
	configBuilder := &ConfigServiceBuilder{
		configService: &ConfigService{
			sectionMap: make(map[string]ConfigSection),
		},
		readers: []ConfigReader{
			&YmlConfigReader{
				fileDiscoveryPathList: []string{
					filepath.Join(exec, "conf"),
					fmt.Sprintf("/etc/zfy/conf.d/%s", builder.eco.ServiceName),
				},
			},
		},
	}

	configBuilder.
		AddConfig("application.yml").
		AddConfig("dsn.yml").
		AddConfig("redis.yml")

	f(configBuilder)
	builder.eco.configService = configBuilder.configService
}

func (builder *EcoSystemBuilder) UseMysql(keyspace string) {
	builder.eco.dsn = NewMysqlDbConnDsn(keyspace)
}
