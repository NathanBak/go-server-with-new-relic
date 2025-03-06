FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app

FROM alpine:latest  
WORKDIR /
COPY --from=builder /out/app /app
ENTRYPOINT ["/app"]
EXPOSE 8081