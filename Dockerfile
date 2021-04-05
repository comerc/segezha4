FROM golang:1.16-alpine AS golang
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /tmp/exe

FROM gcr.io/distroless/base:latest
WORKDIR /app
COPY --from=golang /tmp/exe /runner
# EXPOSE 8008
ENV HEADLESS_IP=173.16.0.42
ENTRYPOINT ["/runner"]

# RUN go get github.com/githubnemo/CompileDaemon
# ENTRYPOINT go run commands/webservice.go
# ENTRYPOINT CompileDaemon --build="go build commands/webservice.go" --command=./webservice