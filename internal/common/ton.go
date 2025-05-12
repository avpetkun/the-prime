package common

import (
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"ts"`
	TonAmount int64     `json:"ton_amount"`
	IsDeposit bool      `json:"is_deposit"`
	TxHash    string    `json:"tx_hash"`
	TxComment string    `json:"tx_comment"`
	InvoiceID string    `json:"invoice_id"`
	SrcAddr   string    `json:"src_addr"`
	DstAddr   string    `json:"dst_addr"`
}

type TonCTx struct {
	ValidUntil int64           `json:"validUntil"`
	Messages   []TonCTxMessage `json:"messages"`
}

type TonCTxMessage struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
	Payload string `json:"payload"`
}

func NewTonTxComment(commentPrefix string, userID, taskID int64) string {
	return commentPrefix + " " + strconv.FormatInt(userID, 10) + " " + strconv.FormatInt(taskID, 10)
}

func ParseTonTxComment(comment string) (userID, taskID int64, success bool) {
	if parts := strings.Fields(comment); len(parts) >= 2 {
		userID, _ = strconv.ParseInt(parts[len(parts)-2], 10, 64)
		taskID, _ = strconv.ParseInt(parts[len(parts)-1], 10, 64)
		success = userID > 0 && taskID > 0
	}
	return
}
