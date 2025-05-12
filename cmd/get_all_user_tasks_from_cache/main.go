package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/avpetkun/the-prime/pkg/math"
	"golang.org/x/sync/errgroup"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/common"
)

func main() {
	c, err := cache.Connect(context.TODO(), cache.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	check(err)
	defer c.Close()

	users := readUsersCsv("./all_users.csv")

	flows := map[int64][]common.TaskFlow{}
	flowsLen := 0
	var mx sync.Mutex

	usersChan := make(chan int64, 100)
	eg := errgroup.Group{}

	for range 20 {
		eg.Go(func() error {
			for user := range usersChan {
				var userTasks map[common.TaskKey]common.TaskFlow
				for {
					userTasks, err = c.GetUserTasks(context.TODO(), user)
					if err == nil {
						break
					}
				}
				var list []common.TaskFlow
				for _, t := range userTasks {
					list = append(list, t)
				}
				if len(list) > 0 {
					mx.Lock()
					flows[user] = list
					flowsLen++
					mx.Unlock()
				}
			}
			return nil
		})
	}

	eta := math.NewETA()
	size := float64(len(users))
	for i, user := range users {
		if i%100 == 0 {
			fmt.Printf("%d / %.0F (%.2F%%) found %d eta %s\n", i, size, float64(i)/size*100, flowsLen, eta.Update(float64(i)/size))
		}
		usersChan <- user
	}
	close(usersChan)

	eg.Wait()

	data, err := json.MarshalIndent(flows, "", "\t")
	check(err)
	os.WriteFile("users_tasks.json", data, os.ModePerm)
}

func readUsersCsv(filename string) (users []int64) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	r := csv.NewReader(f)

	for {
		line, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		userID, err := strconv.ParseInt(line[0], 10, 64)
		check(err)

		users = append(users, userID)
	}
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
