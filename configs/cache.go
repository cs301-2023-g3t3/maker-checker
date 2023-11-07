package configs

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.ClusterClient

func ConnectToRedis() {
	var addr string
	if os.Getenv("ENV") != "lambda" {
		addr = "redis:6379"
	} else {
		addr = os.Getenv("REDIS_HOST")
		log.Println(addr)
	}
	// RedisClient = redis.NewClient(&redis.Options{
	// 	TLSConfig: &tls.Config{
	// 		MinVersion: tls.VersionTLS12,
	// 	},
	// 	Addr:     addr,
	// 	Password: "",
	// 	DB:       0,
	// })
	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{"main-cluster-0001-001.main-cluster.erva1y.apse1.cache.amazonaws.com:6379", "main-cluster-0001-002.main-cluster.erva1y.apse1.cache.amazonaws.com:6379"},
		TLSConfig:      &tls.Config{},
		ReadOnly:       false,
		RouteRandomly:  false,
		RouteByLatency: false,
	})

	ctx := context.Background()
	err := RedisClient.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}
}
