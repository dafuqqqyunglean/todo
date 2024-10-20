FROM golang:1.23.1-alpine

RUN go version

ENV GOPATH=/

COPY . .

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o app ./cmd

EXPOSE 8000

CMD ["./app"]