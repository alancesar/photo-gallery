FROM golang:1.16-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api
ENTRYPOINT ["/app/api"]
