package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

func İnit() *ethclient.Client {
	client, err := ethclient.Dial("ethclient")
	if err != nil {
		zap.L().Error("Failed to connect to the Ethereum client", zap.Error(err))

	}
	return client
}
