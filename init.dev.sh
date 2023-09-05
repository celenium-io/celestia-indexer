#!/bin/sh

# VCS config
git config --local core.attributesfile ./.gitattributes

# Install third tools
go install github.com/swaggo/swag/cmd/swag@latest
curl -fsSL "https://github.com/abice/go-enum/releases/download/v0.5.7/go-enum_$(uname -s)_$(uname -m)" -o ${GOPATH}/bin/go-enum
chmod +x ${GOPATH}/bin/go-enum
go install go.uber.org/mock/mockgen@main
