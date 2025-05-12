package tgu

import "github.com/mymmrac/telego"

type User struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	IsPremium       bool   `json:"is_premium"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
	PhotoUrl        string `json:"photo_url"`
}

func userTelego(u telego.User) User {
	return User{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Username:        u.Username,
		LanguageCode:    u.LanguageCode,
		IsPremium:       u.IsPremium,
		AllowsWriteToPm: true,
		PhotoUrl:        "",
	}
}
