.PHONY: plan apply destroy build

plan:
	cd sandbox && terraform plan -no-color | go run ../main.go

apply:
	cd sandbox && terraform apply -no-color | go run ../main.go

destroy:
	cd sandbox && terraform destroy -no-color | go run ../main.go

install:
	go install github.com/orangekame3/tfsummary@latest

localstack:
	cd sandbox && docker compose -f compose.yml up -d localstack

test:
	cd cmd && go test

coverage:
	cd cmd && go test -covermode=count -coverprofile=c.out && go tool cover -html=c.out -o coverage.html

build:
	cd sandbox && GOOS=linux GOARCH=amd64 go build -o hello && zip lambda.zip hello
