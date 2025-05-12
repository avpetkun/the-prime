package tonu

import (
	"encoding/base64"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

type Transaction struct {
	Timestamp time.Time `json:"ts"`
	Hash      string    `json:"hash"`
	SrcAddr   string    `json:"src"`
	DstAddr   string    `json:"dst"`
	Comment   string    `json:"comment"`
	TonAmount *big.Int  `json:"tonAmount"`
	IsDeposit bool      `json:"isDeposit"`
}

func NewTxWithComment(comment string) (payloadBase64 string, err error) {
	cell, err := wallet.CreateCommentCell(comment)
	if err == nil {
		payloadBase64 = base64.StdEncoding.EncodeToString(cell.ToBOC())
	}
	return
}
