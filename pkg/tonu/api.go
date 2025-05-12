package tonu

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type API struct{ ton.APIClientWrapped }

func ConnectGlobal(ctx context.Context, log zerolog.Logger) (*API, error) {
	log.Info().Msg("[ton] getting global config")
	cfg, err := liteclient.GetConfigFromUrl(ctx, "https://ton.org/global.config.json")
	if err != nil {
		cfg, err = liteclient.GetConfigFromUrl(ctx, "https://ton-blockchain.github.io/global.config.json")
	}
	if err != nil {
		log.Error().Err(err).Msg("[ton] failed to get global config")
		return nil, err
	}
	return Connect(ctx, log, cfg)
}

func Connect(ctx context.Context, log zerolog.Logger, cfg *liteclient.GlobalConfig) (*API, error) {
	client := liteclient.NewConnectionPool()

	log.Info().Msg("[ton] starting")
	err := client.AddConnectionsFromConfig(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Msg("[ton] add conn failed")
		return nil, err
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	return &API{api}, nil
}
