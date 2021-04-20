FROM golang:1.15-alpine3.12 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
RUN go build -o main .

# use alpine image instead of scratch, for support a clear function

FROM alpine:3.12

COPY --from=builder /build/main /

ENTRYPOINT ["/main"]
