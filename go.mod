module github.com/openfinality/openfinality

go 1.26.0

replace github.com/openfinality/openfinality/chaincode/openfinality => ./chaincode/openfinality

require (
	github.com/hyperledger/fabric-gateway v1.13.0
	github.com/openfinality/openfinality/chaincode/openfinality v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.51.0
	google.golang.org/grpc v1.84.0
)

require (
	github.com/hyperledger/fabric-protos-go-apiv2 v0.3.7 // indirect
	github.com/miekg/pkcs11 v1.1.2 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260918162117-cecb64721679 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240123012728-ef4313101c80
