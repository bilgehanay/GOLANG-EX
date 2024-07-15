package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func cError(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	ctx := context.Background()

	clnt := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := clnt.Ping(ctx).Result()
	cError(err)
	fmt.Println(pong)

	err = clnt.Set(ctx, "single", "line value implemantation", 0).Err()
	cError(err)

	value, err := clnt.Get(ctx, "single").Result()
	cError(err)

	fmt.Printf("single %s\n", value)

	err = clnt.Del(ctx, "single").Err()
	cError(err)

	value, err = clnt.Get(ctx, "single").Result()
	if err == redis.Nil {
		fmt.Println("single does not exist")
	} else if err != nil {
		fmt.Println("something went wrong")
	} else {
		fmt.Printf("single %s\n", value)
	}

	err = clnt.HSet(ctx, "hset", map[string]interface{}{
		"name":        "Deneme",
		"explanation": "Multiple Fields",
		"type":        "infinite",
	}).Err()
	cError(err)
	userInfo, err := clnt.HGetAll(ctx, "hset").Result()
	cError(err)
	fmt.Println("hset", userInfo)

}
