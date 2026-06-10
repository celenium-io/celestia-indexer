#!/bin/sh

#################################################### VCS config
git config --local core.attributesfile ./.gitattributes

########################################### Install third tools

# linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.11.3

# for generating swagger specification
go install github.com/swaggo/swag/cmd/swag@latest

# for generating enum types
curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.9.2/go-enum_$(uname -s)_$(uname -m)" -o ${GOPATH}/bin/go-enum
chmod +x ${GOPATH}/bin/go-enum
go install go.uber.org/mock/mockgen@v0.5.2

# for checking licenses used by project and deps
go install github.com/google/go-licenses@latest

# for checking known vulnerabilities in deps
go install golang.org/x/vuln/cmd/govulncheck@latest

# for api test, should have npm installed
npm install -g newman

# for setting up license header in each source code file
go install github.com/vvuwei/update-license@v0.0.1