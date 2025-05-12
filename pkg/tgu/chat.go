package tgu

import "github.com/mymmrac/telego"

type Chat struct {
	ChatID int64  `json:"chatID"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Link   string `json:"link"`
}

func chatTelego(c telego.Chat) Chat {
	link := c.Username
	if link != "" {
		link = "https://t.me/" + link
	}
	return Chat{
		ChatID: c.ID,
		Type:   c.Type,
		Title:  c.Title,
		Link:   link,
	}
}
