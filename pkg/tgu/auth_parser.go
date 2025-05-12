package tgu

import "encoding/base64"

type AuthParser struct {
	botToken []byte
}

func NewAuthParser(botToken string) *AuthParser {
	return &AuthParser{[]byte(botToken)}
}

func (p *AuthParser) ParseBase64(webappInitDataBase64 string) (*Auth, error) {
	decoded, err := base64.StdEncoding.DecodeString(webappInitDataBase64)
	if err != nil {
		return nil, err
	}
	return p.Parse(string(decoded))
}

func (p *AuthParser) Parse(webappInitData string) (*Auth, error) {
	return parseAuth(p.botToken, webappInitData)
}
