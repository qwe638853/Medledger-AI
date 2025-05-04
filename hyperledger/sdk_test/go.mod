module sdk_test

go 1.24.2

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-gateway v1.7.1
	github.com/hyperledger/fabric-protos-go-apiv2 v0.3.4
	github.com/pkg/errors v0.9.1
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
)

replace google.golang.org/protobuf => google.golang.org/protobuf v1.34.1

replace github.com/golang/protobuf => github.com/golang/protobuf v1.4.3
