# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.21-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git musl-dev gcc build-base

RUN mkdir -p $GOPATH/src/github.com/dipdup-io/celestia-indexer/

COPY ./go.* $GOPATH/src/github.com/dipdup-io/celestia-indexer/
WORKDIR $GOPATH/src/github.com/dipdup-io/celestia-indexer
RUN go mod download

COPY cmd/indexer cmd/indexer
COPY internal internal

WORKDIR $GOPATH/src/github.com/dipdup-io/celestia-indexer/cmd/indexer/
RUN go build -a -o /go/bin/indexer .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/celestia-indexer/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/indexer /go/bin/indexer
COPY ./build/dipdup.yml ./

ENTRYPOINT ["/go/bin/indexer", "-c", "dipdup.yml"]