package ecosystem

import (
	"fmt"
	"strconv"

	redis "gopkg.in/redis.v4"
)

type RedisDsn struct {
	keyspace string
}

func NewRedisDsn(keyspace string) *RedisDsn {
	return &RedisDsn{
		keyspace: keyspace,
	}
}

func (dsn *RedisDsn) NewClient() (*redis.Client, error) {
	configService := GetEco().GetConfigService()
	redisDsn := configService.GetConfig(RedisConfigUri)
	dbNum := redisDsn.MustGetValue(fmt.Sprintf("%s.db", dsn.keyspace))
	dbNumInt, err := strconv.Atoi(dbNum)
	if err != nil {
		return nil, err
	}
	addr := redisDsn.MustGetValue(fmt.Sprintf("%s.addr", dsn.keyspace))

	pwd := redisDsn.MustGetValue(fmt.Sprintf("%s.password", dsn.keyspace))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       dbNumInt,
	})
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
