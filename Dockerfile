FROM golang:latest

ADD . /go/app

WORKDIR  /go/app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o main ./main.go 

ENTRYPOINT ./main