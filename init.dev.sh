#!/bin/sh

# VCS config
git config --local core.attributesfile ./.gitattributes

# Install third tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1 # as long we use go@1.20
go install github.com/swaggo/swag/cmd/swag@latest
curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.5.7/go-enum_$(uname -s)_$(uname -m)" -o ${GOPATH}/bin/go-enum
chmod +x ${GOPATH}/bin/go-enum
go install go.uber.org/mock/mockgen@main
go install github.com/google/go-licenses@latest
