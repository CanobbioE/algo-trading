.PHONY:
GO111MODULE=on
PKG_NAME=.



default:

fmt:
	gofmt -s -w $$(go list -f "{{.Dir}}" ./...)

fieldalignment:
	fieldalignment -fix -test=false ./...

lint-diff:
	@golangci-lint run --new-from-rev=$$(git merge-base HEAD master) --timeout 6m0s ./...

lint:
	@golangci-lint run --fix --timeout 6m0s ./...

gci:
	@gci write -s standard -s default -s "prefix(github.com/CanobbioE/algo-trading)" -s blank -s dot ./pkg

install-tools:
	@echo "Installing tools..."
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	@go install go.uber.org/mock/mockgen@v0.4.0
	@echo "Installation completed!"