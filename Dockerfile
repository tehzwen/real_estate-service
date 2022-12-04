FROM golang:1.19.3-alpine3.17 AS builder
WORKDIR /app/build
COPY go.mod .
RUN go mod download
COPY . .
WORKDIR /app/build/cmd/backend
RUN go build -o app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /home/app
COPY --from=builder /app/build/cmd/backend/app ./
CMD [ "./app" ]
