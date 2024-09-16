FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x ./wait_database.sh

RUN go mod download
RUN go build -o app ./cmd/main.go

CMD ["./app"]