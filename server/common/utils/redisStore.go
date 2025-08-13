/*
 * @Author: JimZhang
 * @Date: 2025-07-24 21:54:47
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-24 22:10:06
 * @FilePath: /server/common/utils/redisStore.go
 * @Description: Redis 存取验证码
 *
 */

package utils

import (
	"context"
	"server/common/constant"
	"server/pkg/redis"
	"strings"
	"time"
)

var ctx = context.Background()

type RedisStore struct {
}

func (r RedisStore) Set(id string, value string) error {
	key := constant.LOGIN_CODE + id
	return redis.RedisDb.Set(ctx, key, value, time.Minute*5).Err()
}

func (r RedisStore) Get(id string, clear bool) string {
	key := constant.LOGIN_CODE + id
	value, err := redis.RedisDb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return value
}

// 验证码校验
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)

	// 不区分大小写比较
	return strings.ToLower(v) == strings.ToLower(answer)
}
