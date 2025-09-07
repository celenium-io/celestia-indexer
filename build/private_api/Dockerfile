# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.24.4-alpine AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates

RUN mkdir -p $GOPATH/src/github.com/celenium-io/celestia-indexer/

COPY ./go.* $GOPATH/src/github.com/celenium-io/celestia-indexer/
WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer
RUN go mod download

COPY cmd/private_api cmd/private_api
COPY cmd/api cmd/api
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer/cmd/private_api/
RUN go build -a -installsuffix cgo -o /go/bin/private_api .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/private_api

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/private_api /go/bin/private_api
COPY ./configs/dipdup.yml ./
COPY database database

ENTRYPOINT ["/go/bin/private_api", "-c", "dipdup.yml"]