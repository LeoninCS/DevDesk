package main

import (
	"context"
	"log"
	"moment/internal/CodeShare"
	"moment/internal/Server"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	port := ":8081"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	// ========== Redis 配置 ==========
	redisAddr := os.Getenv("REDIS_ADDR")     // 如 "127.0.0.1:6379"
	redisPass := os.Getenv("REDIS_PASSWORD") // 默认空
	redisDB := 0                             // 默认数据库 0

	var cs *CodeShare.CodeShare

	// ========== 检查 Redis 是否可连接 ==========
	if redisAddr != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPass,
			DB:       redisDB,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("Redis 连接失败（%v），使用内存模式运行。", err)
			cs = CodeShare.New("", "", 0) // 退回内存模式
		} else {
			log.Printf("Redis 已连接: %s", redisAddr)
			cs = CodeShare.New(redisAddr, redisPass, redisDB)
		}
	} else {
		log.Printf("未设置 REDIS_ADDR，使用内存模式运行。")
		cs = CodeShare.New("", "", 0)
	}

	// ========== 启动服务 ==========
	handler := CodeShare.NewHandler(cs)
	srv, err := Server.NewServer(port, handler)
	if err != nil {
		log.Fatalf("服务器初始化失败: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("服务器启动错误: %v", err)
	} else {
		log.Printf("服务器已在 %s 启动成功", srv.Port)
	}
}
