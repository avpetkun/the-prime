package common

import (
	"strconv"
	"strings"
)

func NewStarsInvoicePayload(userID, taskID int64) string {
	return strconv.FormatInt(userID, 10) + " " + strconv.FormatInt(taskID, 10)
}

func ParseStarsInvoicePayload(payload string) (userID, taskID int64, success bool) {
	if i := strings.IndexByte(payload, ' '); i != -1 {
		userID, _ = strconv.ParseInt(payload[:i], 10, 64)
		taskID, _ = strconv.ParseInt(payload[i+1:], 10, 64)
		success = userID > 0 && taskID > 0
	}
	return
}
