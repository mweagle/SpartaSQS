.DEFAULT_GOAL=provision
.PHONY: build test

format:
	go fmt .

vet:
	go vet .

build: format vet
	go build .

explore:
	go run main.go --level info explore

test: build
	go test ./test/...

provision:
	go run main.go --level info provision --s3Bucket $(S3_BUCKET)

describe:
	go run main.go --level info describe --out ./graph.html --s3Bucket $(S3_BUCKET) && open graph.html

delete:
	go run main.go --level info delete
