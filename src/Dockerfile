FROM golang:1.25.3 AS builder
WORKDIR /build
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o app ./cmd/server

FROM scratch
COPY --from=builder /build/app .
CMD ["/app"]