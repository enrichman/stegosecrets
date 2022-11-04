build:
	goreleaser build --single-target --snapshot --rm-dist

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
