FROM golang AS builder
WORKDIR /nlib-app-kv
COPY go.mod /nlib-app-kv/go.mod
COPY go.sum /nlib-app-kv/go.sum
RUN go mod download
COPY . /nlib-app-kv
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine
WORKDIR /nlib-app-kv
COPY --from=builder /nlib-app-kv/nlib-app-kv /nlib-app-kv/nlib-app-kv
ENTRYPOINT ["/nlib-app-kv/nlib-app-kv"]
