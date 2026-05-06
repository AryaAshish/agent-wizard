.PHONY: test fmt verify

test:
	go test ./...

fmt:
	gofmt -w .

verify:
	bash scripts/verify_docs.sh
