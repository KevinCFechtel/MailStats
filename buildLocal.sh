rm deploy/mailstats
GOOS=linux GOARCH=amd64 go build -o deploy/mailstats cmd/mailstats/main.go