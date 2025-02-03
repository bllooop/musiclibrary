FROM golang:1.22

RUN go version
ENV $GOPATH=/

WORKDIR go/src/app

COPY . .

RUN go mod download
RUN go build -o musiclibrary ./cmd/main.go

EXPOSE 8000

CMD ["./musiclibrary"]