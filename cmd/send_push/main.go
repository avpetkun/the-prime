package main

import (
	"os"

	"github.com/mymmrac/telego"
)

func main() {
	bot, err := telego.NewBot(os.Getenv("BOT_TOKEN"))
	check(err)

	ids := []int64{
		// 23432432423424, my
	}

	for _, id := range ids {
		_, err = bot.SendMessage(&telego.SendMessageParams{
			ChatID: telego.ChatID{ID: id},
			Text:   MsgStars250,
		})
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const MsgPremium = "Congratulations! You have received Telegram Premium for 3 months!\nStay with us, complete tasks and earn rewards!"
const MsgStars250 = "Congratulations!\nYou have received 250 Stars!\nStay with us, complete tasks and earn rewards!"
