package engines

import (
	"fmt"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)
}

// PostgresDB returns a handler to the DB object
func RedisDB() *redis.Client {
	return client
}
