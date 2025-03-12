package main

import (
	"net"

	"github.com/Emircaann/wallet-tracer/pkg/handler"
	wallet_tracer "github.com/Emircaann/wallet-tracer/pkg/protos"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		zap.L().Fatal("Failed to starting server", zap.Error(err), zap.String("port", "8080"))
	}
	client, err := ethclient.Dial("ethclient")
	if err != nil {
		zap.L().Error("Failed to connect to the Ethereum client", zap.Error(err), zap.String("url", "ethclient"))
	}
	grpcServer := grpc.NewServer()
	handler := handler.NewHandler(client)
	wallet_tracer.RegisterTransactionServiceServer(grpcServer, handler)
	zap.L().Info("Starting server", zap.String("port", "8080"))
	if err := grpcServer.Serve(lis); err != nil {
		zap.L().Fatal("Failed to serve the server", zap.Error(err))
	}

}
