package tgu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

type Auth struct {
	AuthDate     int64  `json:"auth_date"`
	ChatInstance int64  `json:"chat_instance"`
	ChatType     string `json:"chat_type"`
	User         User   `json:"user"`
}

var tgDataSecretConst = []byte("WebAppData")
var tgDataEndline = []byte("\n")

func ParseAuthBase64(botToken, webappInitDataBase64 string) (*Auth, error) {
	decoded, err := base64.StdEncoding.DecodeString(webappInitDataBase64)
	if err != nil {
		return nil, err
	}
	return ParseAuth(botToken, string(decoded))
}

func ParseAuth(botToken, webappInitData string) (*Auth, error) {
	return parseAuth([]byte(botToken), webappInitData)
}

func parseAuth(botToken []byte, webappInitData string) (*Auth, error) {
	params, err := url.ParseQuery(webappInitData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tg initData: %w", err)
	}

	expectedHash := params.Get("hash")
	if len(expectedHash) == 0 {
		return nil, errors.New("missing hash in tg initData")
	}

	td := new(Auth)

	dataStrings := make([]string, 0, len(params))
	for key, values := range params {
		switch key {
		case "hash":
			continue
		case "auth_date":
			td.AuthDate, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return td, fmt.Errorf("failed to parse auth_data: %w", err)
			}
		case "chat_instance":
			td.ChatInstance, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return td, fmt.Errorf("failed to parse chat_instance: %w", err)
			}
		case "chat_type":
			td.ChatType = values[0]
		case "user":
			err = json.Unmarshal([]byte(values[0]), &td.User)
			if err != nil {
				return td, fmt.Errorf("failed to parse user data: %w", err)
			}
		}
		if len(values) > 0 {
			dataStrings = append(dataStrings, key+"="+values[0])
		}
	}
	sort.Strings(dataStrings)

	h := hmac.New(sha256.New, tgDataSecretConst)
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	h = hmac.New(sha256.New, secretKey)
	for i, s := range dataStrings {
		if i != 0 {
			h.Write(tgDataEndline)
		}
		h.Write([]byte(s))
	}
	actualHash := h.Sum(nil)

	if hex.EncodeToString(actualHash) != expectedHash {
		return td, errors.New("invalid hash of tg initData")
	}
	return td, nil
}
