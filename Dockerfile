FROM golang:1.18-buster

RUN go version
ENV GOPATH=/

#ADD migrate.linux-amd64.tar.gz /util/

COPY ./ ./

RUN tar xzf ./util/migrate.linux-amd64.tar.gz -C ./util

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o toimi-rest ./cmd/main.go

CMD ["./toimi-rest"]