#!/bin/sh

#################################################### VCS config
git config --local core.attributesfile ./.gitattributes

########################################### Install third tools

# linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1 # as long we use go@1.20

# for generating swagger specification
go install github.com/swaggo/swag/cmd/swag@latest

# for generating enum types
curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.5.7/go-enum_$(uname -s)_$(uname -m)" -o ${GOPATH}/bin/go-enum
chmod +x ${GOPATH}/bin/go-enum
go install go.uber.org/mock/mockgen@main

# for checking licenses used by project and deps
go install github.com/google/go-licenses@latest

# for api test, should have npm installed
npm install -g newman
