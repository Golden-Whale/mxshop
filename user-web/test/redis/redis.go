package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	// 验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
	})
	value, err := rdb.Get(context.Background(), "13002672571").Result()
	if err == redis.Nil {
		fmt.Println("key 不存在")
		return
	}
	fmt.Println(value)
}
