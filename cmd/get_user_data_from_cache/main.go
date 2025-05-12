package main

import (
	"context"
	"fmt"
	"os"

	"github.com/avpetkun/the-prime/internal/cache"
)

func main() {
	ctx := context.TODO()
	c, err := cache.Connect(ctx, cache.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	check(err)
	defer c.Close()

	const userID = 123456789

	exist, err := c.CheckUser(ctx, userID)
	check(err)
	inited, err := c.CheckUserInit(ctx, userID)
	check(err)
	fmt.Println("exist", exist, "inited", inited)

	points, err := c.GetUserPoints(ctx, userID)
	fmt.Println("points balance:", points)

	flows, err := c.GetUserTasks(ctx, userID)
	check(err)
	for _, f := range flows {
		fmt.Printf("%+v\n", f)
	}

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
