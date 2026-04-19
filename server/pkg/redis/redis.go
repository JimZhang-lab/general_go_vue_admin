/*
 * @Author: JimZhang
 * @Date: 2025-07-24 11:39:02
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-24 11:39:03
 * @FilePath: /server/pkg/redis/redis.go
 * @Description:
 *
 */
package redis

import (
	"context"
	"server/common/config"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	RedisDb *redis.Client
)

// SetupRedisDb Initialize the Redis instance
func SetupRedisDb() error {
	var ctx = context.Background()
	redisAddress := config.Config.Redis.Host + ":" + strconv.Itoa(config.Config.Redis.Port)
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: config.Config.Redis.Password,
		DB:       9,
	})
	_, err := RedisDb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}
