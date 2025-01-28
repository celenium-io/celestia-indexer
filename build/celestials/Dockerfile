# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates

RUN mkdir -p $GOPATH/src/github.com/celenium-io/celestia-indexer/

COPY ./go.* $GOPATH/src/github.com/celenium-io/celestia-indexer/
WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer
RUN go mod download

COPY cmd/celestials cmd/celestials
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/celenium-io/celestia-indexer/cmd/celestials/
RUN go build -a -installsuffix cgo -o /go/bin/celestials .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/celestials

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/celestials /go/bin/celestials
COPY ./configs/dipdup.yml ./
COPY database database

ENTRYPOINT ["/go/bin/celestials", "-c", "dipdup.yml"]