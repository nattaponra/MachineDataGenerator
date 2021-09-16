FROM golang:1.16 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin .
RUN chmod +x ./bin
CMD ["./bin"]
