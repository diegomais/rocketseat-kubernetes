# Use Go base image
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o server .

# Final stage
FROM alpine:3.21
WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/*.sql .
COPY --from=builder /usr/src/app/server .

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]
