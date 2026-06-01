ARTIFACT_NAME := vole

build:
	@go build -o bin/${ARTIFACT_NAME} main.go

run:
	@go run main.go