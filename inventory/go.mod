module github.com/paincake00/microservices-go/inventory

go 1.26.4

replace github.com/paincake00/microservices-go/shared => ../shared

require (
	github.com/brianvoe/gofakeit/v7 v7.15.0
	github.com/google/uuid v1.6.0
	github.com/paincake00/microservices-go/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.83.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
)
