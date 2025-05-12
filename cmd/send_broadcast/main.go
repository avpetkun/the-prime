package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/avpetkun/the-prime/pkg/math"
	"github.com/mymmrac/telego"
	"golang.org/x/time/rate"
)

func main() {
	users := readUsersCsv("./all_users.csv")
	// users = map[int64]string{
	// 	324324234: "ru",
	// }
	subUsers(users, readUsersTxt("./errs.txt"))

	sent, err := os.Create("./sent.txt")
	check(err)
	defer sent.Close()

	errs, err := os.Create("./errs.txt")
	check(err)
	defer errs.Close()

	msgLimit := rate.NewLimiter(rate.Every(time.Millisecond*50), 1)

	eta := math.NewETA()

	bot, err := telego.NewBot(os.Getenv("BOT_TOKEN"))
	check(err)

	ctx := context.TODO()

	i := 0
	for userID, lang := range users {
		msgLimit.Wait(ctx)
		i++
		perc := float64(i) / float64(len(users))
		fmt.Printf("%d / %d (%.2F%%) %s\n", i, len(users), perc*100, eta.Update(perc))
		sendMessage(userID, lang, sent, errs, bot)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sendMessage(userID int64, lang string, sent, errs io.Writer, bot *telego.Bot) {
	msg := "We decided to boost your rewards to 50% for 24 hours\n\n" +
		"Thanks for your activity"
	if lang == "ru" {
		msg = "Мы решили увеличить ваше вознаграждение на 50% на 24 часа\n\n" +
			"Спасибо вам за вашу активность"
	}

	fmt.Fprintf(sent, "%d\n", userID)

	_, err := bot.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{ID: userID},
		Text:   msg,
	})
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(errs, "%d\n", userID)
	}
}

// c = a - b
func subUsers(a, b map[int64]string) {
	for id := range b {
		delete(a, id)
	}
}

func readUsersCsv(filename string) map[int64]string {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	r := csv.NewReader(f)

	users := make(map[int64]string)

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
		lang := line[1]

		users[userID] = lang
	}
	return users
}

func readUsersTxt(filename string) map[int64]string {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	users := make(map[int64]string)

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		userID, err := strconv.ParseInt(scan.Text(), 10, 64)
		check(err)
		users[userID] = ""
	}
	return users
}
