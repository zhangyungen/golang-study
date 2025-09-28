package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background() // Context 是一个可选的参数，用于设置超时等
var client *redis.Client

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务器地址
		Password: "",               // Redis 密码（没有则留空）
		DB:       0,                // 默认数据库编号，默认是0
		PoolSize: 20,               // 连接池大小
	})
	client.SetNX(ctx, "key", "value", 0)
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	// 设置键值对
	err = client.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		fmt.Println("Failed to set key:", err)
		return
	}

	person := Person{Name: "John Doe", Age: 30}

	setCacheInfo("structKey", person)
	obj, err := getCacheInfo[Person]("structKey")
	if err != nil {
		return
	}
	fmt.Println("getset structKey obj :", obj)

	infoByte, err := json.Marshal(person)
	person = Person{}
	// 设置键值对
	err = client.Set(ctx, "structKey", infoByte, 50*time.Second).Err()
	if err != nil {
		fmt.Println("Failed to set structKey:", err)
		return
	}

	// 获取键的值
	val, err := client.Get(ctx, "structKey").Result()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	err = json.Unmarshal([]byte(val), &person)
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	fmt.Println("structKey :", person)

	// 获取键的值
	val, err = client.Get(ctx, "key1").Result()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}

	fmt.Println("key1:", val)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	// 向列表尾部添加元素
	err = client.RPush(ctx, "list1", "item1", "item2", "item3").Err()
	if err != nil {
		fmt.Println("Failed to push to list:", err)
		return
	}

	defer cancel()
	// 获取列表范围内的元素
	list, err := client.LRange(ctx, "list1", 0, -1).Result()
	if err != nil {
		fmt.Println("Failed to get list items:", err)
		return
	}
	fmt.Println("list1 items:", list)

	//开启事务
	pipe := client.TxPipeline()

	//执行事务操作
	pipe.Incr(ctx, "counter")
	pipe.Expire(ctx, "counter", time.Hour)
	// 执行事务
	pipe.Exec(ctx)

	// 创建订阅者
	pubsub := client.Subscribe(ctx, "channel")

	// 接收消息
	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("Failed to receive message:", err)
		return
	}
	fmt.Println("Received message:", msg.Payload)

	// 发布消息
	err = client.Publish(ctx, "channel", "hello").Err()

	//	### 批量操作
	//
	//	通过管道可以实现批量操作，可以显著提高性能，特别是在需要频繁进行多个命令的情况下。

	pipe = client.Pipeline()

	// 批量设置
	pipe.Set(ctx, "key3", "value3", 0)
	pipe.Set(ctx, "key4", "value4", 0)

	// 执行管道操作
	pipe.Exec(ctx)

	// 关闭连接
	defer client.Close()
}

func setCacheInfo(cacheKey string, cacheInfo interface{}) {
	marshal, err := json.Marshal(cacheInfo)
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}
	client.Set(ctx, cacheKey, marshal, 100000*time.Second)
}

func getCacheInfo[T any](cacheKey string) (T, error) {
	var result T
	info, err := client.Get(ctx, cacheKey).Result()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return result, err
	}
	err = json.Unmarshal([]byte(info), &result)
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return result, err
	}
	return result, nil
}
