FROM golang:1.20-buster

RUN go version
ENV $GOPATH=/

WORKDIR go/src/app

COPY . .

RUN go mod download
RUN go build -o movieapi ./cmd/main.go

EXPOSE 8000

CMD ["./movieapi"]