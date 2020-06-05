FROM golang:alpine AS builder
WORKDIR /jwplayer
COPY . .
RUN GOOS=linux go build -o server .

FROM alpine AS server
WORKDIR /
COPY --from=builder /jwplayer/server /server
ENTRYPOINT ["/server"]