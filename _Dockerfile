FROM golang:alpine

WORKDIR /app

COPY ./ /app

RUN go mod download

# RUN go get github.com/githubnemo/CompileDaemon

# ENTRYPOINT go run commands/webservice.go
ENTRYPOINT go run main.go

# ENTRYPOINT CompileDaemon --build="go build commands/webservice.go" --command=./webservice