# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.23-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates

RUN mkdir -p $GOPATH/src/github.com/celenium-io/celestia-indexer/

COPY ./go.* $GOPATH/src/github.com/celenium-io/celestia-indexer/
WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer
RUN go mod download

COPY cmd/indexer cmd/indexer
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer/cmd/indexer/
RUN go build -a -o /go/bin/indexer .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/celestia-indexer/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/indexer /go/bin/indexer
COPY configs/dipdup.yml ./
COPY database database

ENTRYPOINT ["/go/bin/indexer", "-c", "dipdup.yml"]