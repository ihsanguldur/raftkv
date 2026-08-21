FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/node ./cmd/node

FROM alpine:3.20
COPY --from=builder /out/node /usr/local/bin/node
ENTRYPOINT ["node"]
