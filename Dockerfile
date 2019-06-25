FROM golang:latest

COPY . /app

# Install deps
RUN go get -t github.com/gorilla/mux
RUN go get -t github.com/go-redis/redis
RUN go get -t github.com/joho/godotenv

EXPOSE 8080