# ---------------------------------------------------------------------
#  The first stage container, for building the application
# ---------------------------------------------------------------------
FROM golang:1.21.2-alpine as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux

RUN apk --no-cache add ca-certificates
RUN apk add --update git

RUN mkdir -p $GOPATH/src/github.com/dipdup-io/celestia-indexer/

COPY ./go.* $GOPATH/src/github.com/dipdup-io/celestia-indexer/
WORKDIR $GOPATH/src/github.com/dipdup-io/celestia-indexer
RUN go mod download

COPY cmd/api cmd/api
COPY internal internal
COPY pkg pkg

WORKDIR $GOPATH/src/github.com/dipdup-io/celestia-indexer/cmd/api/
RUN go build -a -installsuffix cgo -o /go/bin/api .

# ---------------------------------------------------------------------
#  The second stage container, for running the application
# ---------------------------------------------------------------------
FROM scratch

WORKDIR /app/api

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/api /go/bin/api
COPY ./build/dipdup.yml ./

ENTRYPOINT ["/go/bin/api", "-c", "dipdup.yml"]