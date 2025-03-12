package handler

import (
	"math/big"
	"time"

	wallet_tracer "github.com/Emircaann/wallet-tracer/pkg/protos"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type Handler struct {
	ethclient *ethclient.Client
	wallet_tracer.UnimplementedTransactionServiceServer
}

func NewHandler(ethclient *ethclient.Client) *Handler {
	return &Handler{
		ethclient: ethclient,
	}
}

func (h *Handler) WatchTransactions(req *wallet_tracer.WatchRequest, stream wallet_tracer.TransactionService_WatchTransactionsServer) error {
	address := req.GetAddress()
	var lastBlock uint64

	for {

		blockNumber, err := h.ethclient.BlockNumber(stream.Context())
		if err != nil {
			zap.L().Error("Failed to fetch block number", zap.Error(err))
			time.Sleep(10 * time.Second)
			continue
		}

		if blockNumber > lastBlock {
			lastBlock = blockNumber
			block, err := h.ethclient.BlockByNumber(stream.Context(), big.NewInt(int64(blockNumber)))
			if err != nil {
				zap.L().Error("Failed to get block by number", zap.Error(err))
				continue
			}

			for _, tx := range block.Transactions() {
				if tx.To() != nil && tx.To().String() == address {
					zap.L().Info("Transaction detected", zap.String("hash", tx.Hash().String()), zap.String("amount", tx.Value().String()))
					if err := stream.Send(&wallet_tracer.TransactionResponse{
						TxHash: tx.Hash().String(),
						Amount: tx.Value().String(),
					}); err != nil {
						zap.L().Error("Failed to send transaction response", zap.Error(err))
						return err
					}
				}
			}
		}

		time.Sleep(7 * time.Second)
	}
}
