package main

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/internal/worker"
	"github.com/avpetkun/the-prime/pkg/natsu"
	"github.com/avpetkun/the-prime/pkg/signalu"
	"github.com/avpetkun/the-prime/pkg/tgu"
	"github.com/rs/zerolog"
)

func main() {
	ctx, _ := signalu.WaitExitContext(context.TODO())

	db, err := dbx.Connect(context.TODO(), dbx.Config{
		User: os.Getenv("PG_USER"),
		Pass: os.Getenv("PG_PASS"),
		Addr: os.Getenv("PG_ADDR"),
		Name: os.Getenv("PG_NAME"),
	})
	check(err)
	defer db.Close()

	stream, err := natsu.Connect(zerolog.New(os.Stdout), os.Getenv("NATS_URL"))
	check(err)
	defer stream.Stop()

	products, err := db.GetAllProducts(ctx)
	check(err)

	productsMap := make(map[int64]*common.Product)
	for _, p := range products {
		productsMap[p.ID] = p
	}

	users := readUsersCsv("./users.csv")

	tickets := readTicketsCsv("./data-1745239657688.csv")

	for i, t := range tickets {
		log.Println(i, len(tickets))

		product := productsMap[t.ProductID]
		username, ok := users[t.UserID]
		if !ok || product == nil {
			panic(t)
		}

		err = stream.Publish(ctx, worker.KeyFragmentSend, worker.FragmentMessage{
			ProductClaimMessage: worker.ProductClaimMessage{
				Product: product,
				User:    tgu.User{ID: t.UserID, Username: username},
			},
			TicketID: t.TicketID,
		})
		check(err)
	}

	log.Println("done")
	<-ctx.Done()
}

type Ticket struct {
	TicketID  int64
	UserID    int64
	ProductID int64
}

func readTicketsCsv(filename string) (tickets []Ticket) {
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

		ticketID, err := strconv.ParseInt(line[0], 10, 64)
		check(err)
		userID, err := strconv.ParseInt(line[1], 10, 64)
		check(err)
		productID, err := strconv.ParseInt(line[2], 10, 64)
		check(err)

		tickets = append(tickets, Ticket{
			TicketID:  ticketID,
			UserID:    userID,
			ProductID: productID,
		})
	}
	return
}

func readUsersCsv(filename string) (usernames map[int64]string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	r := csv.NewReader(f)

	usernames = make(map[int64]string)
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

		usernames[userID] = line[1]
	}
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
