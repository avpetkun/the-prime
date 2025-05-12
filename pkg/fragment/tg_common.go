package fragment

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

var rxApiHash = regexp.MustCompile(`\/api\?hash=([^"]+)`)

func getApiHash(body []byte) (apiHash string, err error) {
	matches := rxApiHash.FindSubmatch(body)
	if len(matches) < 2 || len(matches[1]) == 0 {
		return "", fmt.Errorf("api hash not found")
	}
	return string(matches[1]), nil
}

type SearchRecipientResponse struct {
	Error string `json:"error"`
	OK    bool   `json:"ok"`
	Found struct {
		// Myself    bool   `json:"myself"`
		Recipient string `json:"recipient"`
		// Photo     string `json:"photo"`
		// Name      string `json:"name"`
	} `json:"found"`
}

type InitBuyResponse struct {
	Error string `json:"error"`
	ReqID string `json:"req_id"`
	// Myself    bool   `json:"myself"`
	// ToBot     bool   `json:"to_bot"`
	Amount string `json:"amount"`
	// ItemTitle string `json:"item_title"`
	// Content   string `json:"content"`
	// Button    string `json:"button"`
}

type GetLinkResponse struct {
	Error         string `json:"error"`
	OK            bool   `json:"ok"`
	Transaction   TonTx  `json:"transaction"`
	ConfirmMethod string `json:"confirm_method"`
	ConfirmParams struct {
		ID string `json:"id"`
	} `json:"confirm_params"`
}

type TonTx struct {
	ValidUntil int64           `json:"validUntil"`
	From       string          `json:"from"`
	Messages   []*TonTxMessage `json:"messages"`
}

type TonTxMessage struct {
	Address *address.Address `json:"address"`
	Amount  *big.Int         `json:"amount"`
	Payload string           `json:"payload"`
}

func (tx *TonTx) GetMessage() (*wallet.Message, error) {
	if len(tx.Messages) != 1 {
		return nil, fmt.Errorf("tx messages count != 1")
	}
	msg := tx.Messages[0]
	boc64 := msg.Payload
	for len(boc64)%4 != 0 {
		boc64 += "="
	}
	boc, err := base64.URLEncoding.DecodeString(boc64)
	if err != nil {
		return nil, err
	}
	body, err := cell.FromBOC(boc)
	if err != nil {
		return nil, err
	}
	message := &wallet.Message{
		Mode: wallet.PayGasSeparately + wallet.IgnoreErrors,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      true,
			DstAddr:     msg.Address,
			Amount:      tlb.FromNanoTON(msg.Amount),
			Body:        body,
		},
	}
	return message, nil
}

type SentTx struct {
	Hash string `json:"hash"`
	LT   uint64 `json:"lt"`
}
