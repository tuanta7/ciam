ENV_FILE=.env
MIGRATIONS_FOLDER=./migrations
PROTO_FOLDER=protobuf/proto
BUF_VERSION?=1.58.0

env-example:
	awk -F'=' 'BEGIN {OFS="="} \
    	/^[[:space:]]*#/ {print; next} \
    	/^[[:space:]]*$$/ {print ""; next} \
    	NF>=1 {gsub(/^[[:space:]]+|[[:space:]]+$$/, "", $$1); print $$1"="}' .env > .env.example
	echo ".env.example generated successfully."

setup:
	docker compose -f ./build/docker-compose.local.yml up -d 

run-local-op:
	go run ./cmd/op

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	ls "$(shell go env GOPATH)/bin/" | grep goose

migrate-sql:
	goose -dir=$(MIGRATIONS_FOLDER) create $(NAME) sql

migrate-up:
	goose -env $(ENV_FILE) up

migrate-down:
	goose -env $(ENV_FILE) down

install-sqlboiler:
	go install github.com/aarondl/sqlboiler/v4@latest
	go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest

sqlboiler-gen:
	echo "Generating Go ORM models from the database schema using sqlboiler"
	sqlboiler psql -o internal/repository/orm -p orm --no-tests --wipe

mockery-gen:
	echo "Generating mock implementations using mockery"
	docker run --rm -v $(PWD):/src -w /src vektra/mockery:v3.7.0
