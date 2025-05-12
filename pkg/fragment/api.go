package fragment

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"strings"

	"github.com/rs/zerolog"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/avpetkun/the-prime/pkg/tonu"
)

type API struct {
	wal  *wallet.Wallet
	ton  *tonu.API
	auth string

	publicKey  string
	walAddress string
	walRawAddr string
}

func NewAPI(ton *tonu.API, log zerolog.Logger, fragmentCookie, walletW5Phrase string) (*API, error) {
	walletSeed := strings.Fields(walletW5Phrase)
	wal, err := wallet.FromSeed(ton, walletSeed, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	})
	if err != nil {
		return nil, err
	}

	publicKey := hex.EncodeToString(wal.PrivateKey().Public().(ed25519.PublicKey))
	walAddress := wal.WalletAddress().String()
	walRawAddr := wal.WalletAddress().StringRaw()

	log.Info().
		Stringer("wallet", wal.WalletAddress()).
		Str("raw_addr", walRawAddr).
		Str("pub_key", publicKey).
		Msg("[fragment] api created")

	api := &API{
		ton:  ton,
		wal:  wal,
		auth: fragmentCookie,

		publicKey:  publicKey,
		walAddress: walAddress,
		walRawAddr: walRawAddr,
	}
	return api, nil
}

func (api *API) WalletAddress() string {
	return api.walAddress
}

func (api *API) GetBalance(ctx context.Context) (float64, error) {
	master, err := api.ton.CurrentMasterchainInfo(ctx)
	if err != nil {
		return 0, err
	}
	balance, err := api.wal.GetBalance(ctx, master)
	if err != nil {
		return 0, err
	}

	fbalance, _ := balance.Nano().Float64()
	fbalance /= 1e9

	return fbalance, nil
}

func (api *API) sendTx(ctx context.Context, log zerolog.Logger, msg *wallet.Message) (*SentTx, error) {
	log.Info().Msg("[fragment] send tx")
	tx, _, err := api.wal.SendWaitTransaction(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msg("[fragment] failed to send tx")
		return nil, err
	}
	txHash := hex.EncodeToString(tx.Hash)
	log.Info().
		Str("tx_hash", txHash).
		Uint64("tx_lt", tx.LT).
		Msg("[fragment] tx successfully sent")

	return &SentTx{Hash: txHash, LT: tx.LT}, nil
}
