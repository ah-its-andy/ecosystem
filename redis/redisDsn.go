package redis

import (
	redis "gopkg.in/redis.v4"
)

type RedisDsn struct {
	keyspace string
	opts     *redis.Options
}

func NewRedisDsn(keyspace string, addr string, pwd string, db int) *RedisDsn {
	return &RedisDsn{
		keyspace: keyspace,
		opts: &redis.Options{
			Addr:     addr,
			Password: pwd,
			DB:       db,
		},
	}
}

func (dsn *RedisDsn) NewClient() (*redis.Client, error) {
	client := redis.NewClient(dsn.opts)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
