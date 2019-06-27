FROM golang:latest

ADD . /app
WORKDIR /app

# Install deps
RUN go get github.com/gorilla/mux
RUN go get github.com/go-redis/redis


CMD ["go", "run", "main.go"]

EXPOSE 8080