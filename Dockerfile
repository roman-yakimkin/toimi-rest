FROM golang:1.18-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz -O - | tar -xz

RUN go mod download
RUN go build -o toimi-rest ./cmd/main.go

CMD ["./toimi-rest"]