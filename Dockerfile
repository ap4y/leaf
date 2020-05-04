FROM golang:alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -tags netgo -ldflags '-extldflags "-static" -s -w' -o /leaf-server ./cmd/leaf-server

FROM scratch

COPY --from=builder /leaf-server /leaf-server

EXPOSE 8000

ENTRYPOINT ["/leaf-server"]
