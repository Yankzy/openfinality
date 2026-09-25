module github.com/openfinality/openfinality

go 1.26.0

replace github.com/openfinality/openfinality/chaincode/afrorail => ./chaincode/afrorail

require (
	github.com/openfinality/openfinality/chaincode/afrorail v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.50.0
)

require (
	github.com/hyperledger/fabric-gateway v1.12.1 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
)
