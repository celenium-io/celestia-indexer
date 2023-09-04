#!/bin/sh

# VCS config
git config --local core.attributesfile ./.gitattributes

# Install third tools
go install github.com/swaggo/swag/cmd/swag@latest
go install go.uber.org/mock/mockgen@latest
