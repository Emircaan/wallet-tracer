package handler

import (
	"math/big"
	"time"

	wallet_tracer "github.com/Emircaann/wallet-tracer/pkg/protos"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Handler struct {
	ethclient *ethclient.Client
	wallet_tracer.UnimplementedTransactionServiceServer
	tracer trace.Tracer
}

func NewHandler(ethclient *ethclient.Client, tracer trace.Tracer) *Handler {
	return &Handler{
		ethclient: ethclient,
		tracer:    tracer,
	}
}

func (h *Handler) WatchTransactions(req *wallet_tracer.WatchRequest, stream wallet_tracer.TransactionService_WatchTransactionsServer) error {
	ctx, transactionSpan := h.tracer.Start(stream.Context(), "WatchTransactions",
		trace.WithAttributes(
			attribute.String("service", "wallet-tracer"),
			attribute.String("method", "WatchTransactions"),
			attribute.String("address", req.GetAddress()),
		))

	address := req.GetAddress()
	var lastBlock uint64

	for {
		blockSpanCtx, blockSpan := h.tracer.Start(ctx, "GetBlockNumber", trace.WithAttributes(
			attribute.String("service", "wallet-tracer"),
			attribute.String("method", "GetBlockNumber"),
		))

		blockNumber, err := h.ethclient.BlockNumber(blockSpanCtx)
		blockSpan.End()
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
		transactionSpan.End()

		time.Sleep(7 * time.Second)
	}
}
