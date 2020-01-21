all:
	env GOBIN="$(PWD)/bin"  go install -i ./cmd/client/
	env GOBIN="$(PWD)/bin"  go install -i ./cmd/server/
