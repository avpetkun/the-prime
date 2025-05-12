package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/avpetkun/the-prime/internal/cache"
)

func main() {
	c, err := cache.Connect(context.TODO(), cache.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Println(c.GetTaskClicks(context.TODO(), 14))
}
