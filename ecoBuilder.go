package ecosystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ah-its-andy/ecosystem/config"
	"github.com/ah-its-andy/ecosystem/db"
	ecohttp "github.com/ah-its-andy/ecosystem/http"
	"github.com/ah-its-andy/ecosystem/logging"
	"github.com/ah-its-andy/ecosystem/redis"
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
	configReader := config.NewYmlConfigReader([]string{
		filepath.Join(execPath, "conf"),
	})
	config := configReader.GetSection("application.yml")

	if serviceName, ok := config.GetString("serviceName"); !ok {
		panic("服务没有配置serviceName")
	} else {
		ecoInstance := &EcoSystem{
			ServiceName: serviceName,
			ServiceCfg:  config,
			ExecutePath: execPath,
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

func (builder *EcoSystemBuilder) UseExecutePath(f func() string) {
	builder.eco.ExecutePath = f()
}

func (builder *EcoSystemBuilder) UseConfigService(f func(*config.ConfigServiceBuilder)) {
	configBuilder := config.NewConfigServiceBuilder([]config.ConfigReader{
		config.NewYmlConfigReader(builder.eco.GetDefaultConfigDiscoveryPathList()),
	})

	configBuilder.
		AddConfig("application.yml").
		AddConfig("dsn.yml").
		AddConfig("redis.yml")

	f(configBuilder)

	builder.eco.configService = configBuilder.Build()
}

func (builder *EcoSystemBuilder) UseMysql(keyspace string) {
	databaseConfig := builder.eco.configService.GetConfig("dsn.yml")
	dbConnDsn := databaseConfig.MustGetString(keyspace)
	builder.eco.dsn = db.NewMysqlDbConnDsn(dbConnDsn, keyspace)
}

func (builder *EcoSystemBuilder) UseRedis(keyspace string) {
	redisCfg := builder.eco.configService.GetConfig("redis.yml")
	dbNum := redisCfg.MustGetString(fmt.Sprintf("%s.db", keyspace))
	dbNumInt, err := strconv.Atoi(dbNum)
	if err != nil {
		panic(errors.Unwrap(err))
	}
	addr := redisCfg.MustGetString(fmt.Sprintf("%s.addr", keyspace))
	pwd := redisCfg.MustGetString(fmt.Sprintf("%s.password", keyspace))
	builder.eco.redis = redis.NewRedisDsn(keyspace, addr, pwd, dbNumInt)
}

func (builder *EcoSystemBuilder) UseLogging(path string) {
	smartLogger := logging.NewSmartLogger(builder.eco.ServiceName, path)
	builder.eco.smartLogger = smartLogger
}

func (builder *EcoSystemBuilder) UseWebServer(f func(ecohttp.Application)) {
	builder.eco.app = ecohttp.NewWebApplication(builder.eco.smartLogger, builder.eco.ServiceCfg)
	f(builder.eco.app)
}
