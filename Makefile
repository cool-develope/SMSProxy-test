.PHONY: test
test:
	go test ./... -v -count=1 -race