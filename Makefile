ENV_FILE=.env
MIGRATIONS_FOLDER=./data/migrations
PROTO_FOLDER=protobuf/proto
BUF_VERSION?=1.58.0

build-oidc:
	CGO_ENABLED=0 GOOS=linux go build -o ciam ./cmd/oidc

env-example:
	awk -F'=' 'BEGIN {OFS="="} \
    	/^[[:space:]]*#/ {print; next} \
    	/^[[:space:]]*$$/ {print ""; next} \
    	NF>=1 {gsub(/^[[:space:]]+|[[:space:]]+$$/, "", $$1); print $$1"="}' .env > .env.example
	echo ".env.example generated successfully."

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	ls "$(shell go env GOPATH)/bin/" | grep goose

migrate-sql:
	goose -dir=$(MIGRATIONS_FOLDER)/postgres create $(NAME) sql

migrate-up:
	goose -env $(ENV_FILE) up

migrate-down:
	goose -env $(ENV_FILE) down

install-sqlboiler:
	go install github.com/aarondl/sqlboiler/v4@latest
	go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest

sqlboiler-gen:
	echo "Generating Go models from the database schema using sqlboiler"
	sqlboiler psql -o internal/repository/models -p models --no-tests --wipe

mockery-gen:
	echo "Generating mock implementations using mockery"
	docker run --rm -v $(PWD):/src -w /src vektra/mockery:v3.7.0

install-buf:
	go install github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}

buf-dev:
	buf dep update
	buf export buf.build/bufbuild/protovalidate --output=.

buf-gen:
	buf dep update
	buf generate
