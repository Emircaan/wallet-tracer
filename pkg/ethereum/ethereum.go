package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

func İnit() *ethclient.Client {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/8741abbfceb34f96bcd929b2aba2985e")
	if err != nil {
		zap.L().Error("Failed to connect to the Ethereum client", zap.Error(err))

	}
	return client
}
