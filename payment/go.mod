module github.com/paincake00/microservices-go/payment

go 1.26.4

replace github.com/paincake00/microservices-go/shared => ../shared

replace github.com/paincake00/microservices-go/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/paincake00/microservices-go/platform v0.0.0-00010101000000-000000000000
	github.com/paincake00/microservices-go/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.12.1
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.83.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260908043556-f8649ddbbfe6 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
