package db

import (
	"context"
	"sync"
	"time"

	"github.com/clin211/go-cicd-github-actions-docker/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	return redisClient
}

// InitializeRedis 初始化Redis连接
func InitializeRedis(logger *zap.Logger) error {
	var err error

	redisOnce.Do(func() {
		// 创建Redis客户端
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
			PoolSize: config.Redis.PoolSize,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = redisClient.Ping(ctx).Result()
		if err != nil {
			logger.Error("Failed to connect to Redis", zap.Error(err))
			return
		}

		logger.Info("Connected to Redis", zap.String("addr", config.Redis.Addr))
	})

	return err
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// RedisGet 从Redis获取字符串值
func RedisGet(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

// RedisSet 设置Redis字符串值和过期时间
func RedisSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

// RedisDelete 删除一个或多个键
func RedisDelete(ctx context.Context, keys ...string) error {
	return redisClient.Del(ctx, keys...).Err()
}

// RedisExists 检查键是否存在
func RedisExists(ctx context.Context, keys ...string) (bool, error) {
	n, err := redisClient.Exists(ctx, keys...).Result()
	return n > 0, err
}

// RedisExpire 设置键的过期时间
func RedisExpire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return redisClient.Expire(ctx, key, expiration).Result()
}

// RedisHSet 设置哈希表字段的字符串值
func RedisHSet(ctx context.Context, key, field string, value interface{}) error {
	return redisClient.HSet(ctx, key, field, value).Err()
}

// RedisHGet 获取哈希表中指定字段的值
func RedisHGet(ctx context.Context, key, field string) (string, error) {
	return redisClient.HGet(ctx, key, field).Result()
}

// RedisHGetAll 获取哈希表中所有的字段和值
func RedisHGetAll(ctx context.Context, key string) (map[string]string, error) {
	return redisClient.HGetAll(ctx, key).Result()
}
