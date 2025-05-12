package main

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/avpetkun/the-prime/internal/worker"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/signalu"
	"github.com/rs/zerolog"
)

func main() {
	ctx, _ := signalu.WaitExitContext(context.TODO())

	stream, err := natsu.Connect(zerolog.New(os.Stdout), os.Getenv("NATS_URL"))
	check(err)
	defer stream.Stop()

	const taskID = 22
	users := readUsersCsv("./data-1744290550382.csv")

	for i, userID := range users {
		log.Println(i, len(users))
		err = stream.Publish(ctx, worker.KeyTaskDone, worker.TaskMessage{
			Time:   time.Now(),
			UserID: userID,
			TaskID: taskID,
			SubID:  0,
			Strict: false,
		})
		check(err)
	}

	log.Println("done")
	<-ctx.Done()
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
