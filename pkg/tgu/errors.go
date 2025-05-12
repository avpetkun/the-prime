package tgu

import (
	"strings"
	"time"
)

// telego: sendDocument: api: 429 "Too Many Requests: retry after 34", migrate to chat ID: 0, retry after: 34
func ParseErrTooManyRequests(err error) (ok bool, retry time.Duration) {
	if err == nil {
		return
	}
	text := err.Error()
	if !strings.Contains(text, "Too Many Requests") {
		return
	}
	ok = true
	if i := strings.Index(text, "retry after: "); i != -1 {
		retry, err = time.ParseDuration(text[i+13:] + "s")
		if err == nil {
			return
		}
	}
	retry = time.Minute
	return
}
