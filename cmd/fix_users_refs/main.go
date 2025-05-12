package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/math"
)

func main() {
	c, err := cache.Connect(context.TODO(), cache.Config{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	db, err := dbx.Connect(context.TODO(), dbx.Config{
		User: os.Getenv("PG_USER"),
		Pass: os.Getenv("PG_PASS"),
		Addr: os.Getenv("PG_ADDR"),
		Name: os.Getenv("PG_NAME"),
	})
	check(err)
	defer db.Close()

	users := readUsersCsv("./users.csv")
	println(len(users))

	eta := math.NewETA()

	for i, r := range users {
		fmt.Println(i, len(users), r.UserID, eta.Update(float64(i)/float64(len(users))))

		refPoints, refCount, err := db.GetUserRefs(context.TODO(), r.UserID)
		check(err)
		err = c.SetUserRefCount(context.TODO(), r.UserID, refCount)
		check(err)
		err = c.SetUserRefPoints(context.TODO(), r.UserID, refPoints)
		check(err)
	}
}

type Record struct {
	UserID    int64
	RefPoints int64
	RefCount  int64
}

func readUsersCsv(filename string) (list []Record) {
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
		var r Record
		r.UserID, err = strconv.ParseInt(line[0], 10, 64)
		check(err)
		r.RefPoints, err = strconv.ParseInt(line[1], 10, 64)
		check(err)
		r.RefCount, err = strconv.ParseInt(line[2], 10, 64)
		check(err)

		if r.RefPoints != 0 || r.RefCount != 0 {
			list = append(list, r)
		}
	}
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
