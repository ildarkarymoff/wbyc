FROM golang:latest

COPY . /app

# Install deps
RUN go get github.com/gorilla/mux
RUN go get github.com/go-redis/redis

EXPOSE 8080