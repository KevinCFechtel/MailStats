package logger

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
    client *http.Client
	url	string
}

func NewLogger(url string) *Logger {
    logger := &Logger{}
	logger.client = &http.Client{
        Timeout: 10 * time.Second,
    }
	logger.url = url

    return logger
}

func (logger *Logger) LogThis(message string, error bool) {
	if(error) {
		resp, err := logger.client.Post(logger.url + "/fail", "text/plain; charset=utf-8",  strings.NewReader(message))
		if err != nil {
			log.Println("Failed to send error log: " + err.Error())
		}
		if resp.StatusCode != 200 {
			log.Println("Failed to send error log, response code: " + resp.Status)
		}
	} else {
		resp, err := logger.client.Post(logger.url, "text/plain; charset=utf-8",  strings.NewReader(message))
		if err != nil {
			log.Println("Failed to send success log: " + err.Error())
		}
		if resp.StatusCode != 200 {
			log.Println("Failed to send success log, response code: " + resp.Status)
		}
	}
}