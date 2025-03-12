build_proto:
	protoc --proto_path=pkg/protos \
		--go_out=pkg/protos --go_opt=paths=source_relative \
		--go-grpc_out=pkg/protos --go-grpc_opt=paths=source_relative \
		pkg/protos/*.proto
