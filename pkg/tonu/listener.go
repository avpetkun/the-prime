package tonu

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/avpetkun/the-prime/pkg/timeu"
)

type Listener struct {
	API *API
	Log zerolog.Logger

	WatchAddress string

	ProcessedTransaction func(ctx context.Context, txHash string) (exist bool, err error)
	ProcessTransaction   func(ctx context.Context, tx *Transaction) error
}

func (lis *Listener) Start(ctx context.Context) error {
	if lis.API == nil {
		return fmt.Errorf("[ton] connection is no set")
	}
	if lis.WatchAddress == "" {
		return fmt.Errorf("[ton] watch address is no set")
	}
	if lis.ProcessedTransaction == nil {
		return fmt.Errorf("[ton] func ProcessedTransaction is no set")
	}
	if lis.ProcessTransaction == nil {
		return fmt.Errorf("[ton] func ProcessTransaction is no set")
	}

	watchAddr, err := ParseAnyAddress(lis.WatchAddress)
	if err != nil {
		return fmt.Errorf("[ton] failed to parse watch address: %w", err)
	}

	go lis.runFetchLoop(ctx, watchAddr)
	return nil
}

func (lis *Listener) runFetchLoop(ctx context.Context, watchAddr *address.Address) {
	lis.Log.Info().Msg("[ton] started")

mainLoop:
	for {
		if timeu.SleepContext(ctx, time.Second*2) {
			return
		}

		master, err := lis.API.CurrentMasterchainInfo(ctx)
		if err != nil {
			lis.Log.Warn().Err(err).Msg("[ton] failed to get masterchain info")
			continue
		}

		acc, err := lis.API.GetAccount(ctx, master, watchAddr)
		if err != nil {
			lis.Log.Warn().Err(err).Msg("[ton] failed to get wallet account")
			continue
		}

		lt := acc.LastTxLT
		hash := acc.LastTxHash

		var transactions []*Transaction

	listLoop:
		for {
			var txs []*tlb.Transaction
			for try := 0; ; try++ {
				txs, err = lis.API.ListTransactions(ctx, watchAddr, 1, lt, hash)
				if err == nil {
					break
				}
				if errors.Is(err, ton.ErrNoTransactionsWereFound) {
					break listLoop
				}
				lis.Log.Warn().Err(err).Int("try", try).Msg("[ton] failed to call ListTransactions")
				if timeu.SleepContext(ctx, time.Second) {
					return
				}
			}
			for _, tx := range txs {
				txx := unpackTransferTransaction(tx)
				if txx != nil {
					exist, err := lis.ProcessedTransaction(ctx, txx.Hash)
					if err != nil {
						continue mainLoop
					}
					if exist {
						break listLoop
					}
					transactions = append(transactions, txx)
				}
				lt = tx.PrevTxLT
				hash = tx.PrevTxHash
			}
		}

		for i := len(transactions) - 1; i >= 0; i-- {
			tx := transactions[i]

			log := lis.Log.With().
				Str("src", tx.SrcAddr).
				Str("dst", tx.DstAddr).
				Str("tx_hash", tx.Hash).
				Str("tx_comment", tx.Comment).
				Bool("is_deposit", tx.IsDeposit).
				Stringer("tx_amount", tx.TonAmount).
				Logger()
			if err = lis.ProcessTransaction(ctx, tx); err != nil {
				log.Error().Err(err).Msg("[ton] failed to process transaction")
				continue mainLoop
			}
			log.Info().Msg("[ton] processed transaction")
		}
	}
}

func unpackTransferTransaction(tx *tlb.Transaction) *Transaction {
	msg, received, unpacked := unpackTransferMessage(tx)
	if !unpacked || msg.Bounced {
		return nil
	}
	amount := msg.Amount.Nano()
	if amount.Sign() != 1 {
		return nil
	}
	return &Transaction{
		Timestamp: time.Unix(int64(msg.CreatedAt), 0),
		Hash:      hex.EncodeToString(tx.Hash),
		SrcAddr:   msg.SrcAddr.String(),
		DstAddr:   msg.DstAddr.String(),
		Comment:   msg.Comment(),
		TonAmount: amount,
		IsDeposit: received,
	}
}

func unpackTransferMessage(tx *tlb.Transaction) (msg *tlb.InternalMessage, received, unpacked bool) {
	if tx.IO.Out != nil {
		msgs, err := tx.IO.Out.ToSlice()
		if err == nil && len(msgs) == 1 && msgs[0].MsgType == tlb.MsgTypeInternal {
			return msgs[0].AsInternal(), false, true
		}
	} else if tx.IO.In.MsgType == tlb.MsgTypeInternal {
		return tx.IO.In.AsInternal(), true, true
	}
	return
}
